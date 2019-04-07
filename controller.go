package trader

import (
	"log"
	"math"

	kt "github.com/zerodhatech/gokiteconnect"
)

type Controller struct {
	instrument *kt.Instrument
	position   *Position
	strategy   *Strategy
	broker     Broker
}

func NewController(inst *kt.Instrument, bars *Bars, broker Broker) *Controller {
	return &Controller{
		instrument: inst,
		position:   nil,
		strategy:   NewStrategy(bars),
		broker:     broker,
	}
}

//TODO: (Make Fund allocation better)
func (c *Controller) getFunds() float64 {
	return c.broker.GetAvailableFunds() / 10
}

func (c *Controller) OnBar(b *Bar) {

	// Check Stoploss
	if c.position != nil {
		if c.position.StoplossHit(b) {
			if c.position.Type() == BOUGHT {
				c.broker.Sell(c.instrument, b.Close, c.position.Shares())
				log.Print("Exit Bought")
				c.position = nil
			} else if c.position.Type() == BORROWED {
				c.broker.Buy(c.instrument, b.Close, c.position.Shares())
				log.Print("Exit Borrowed")
				c.position = nil
			}
		}
		c.position.StoplossPrice = math.Max(c.position.StoplossPrice, b.Close*0.95)
	}

	signal := c.strategy.OnBar(b)
	switch signal {
	case HOLD:
		log.Print("Hold")
	case BUY:
		if (c.position != nil) && (c.position.Type() == BORROWED) {
			c.broker.Buy(c.instrument, b.Close, c.position.Shares())
			log.Print("Exit Borrowed")
			c.position = nil
		}
		if c.position == nil {
			shares := int(c.getFunds() / b.Close)
			if shares > 0 {
				c.broker.Buy(c.instrument, b.Close, shares)
				log.Print("Buy Shares")
				c.position = NewPosition(BOUGHT, shares, 0.95*b.Close)
			}
		} else {
			shares := int(c.getFunds()/b.Close) - c.position.Shares()
			if shares > 0 {
				c.broker.Buy(c.instrument, b.Close, shares)
				log.Print("Bought More")
				c.position.AddShares(shares)
			}
		}
	case SELL:
		if (c.position == nil) || (c.position.Type() == BORROWED) {
			return
		}
		c.broker.Sell(c.instrument, b.Close, c.position.Shares())
		log.Print("Exit Bought")
		c.position = nil
	case SHORT:
		if (c.position != nil) && (c.position.Type() == BOUGHT) {
			c.broker.Sell(c.instrument, b.Close, c.position.Shares())
			log.Print("Exit Bought")
			c.position = nil
		}
		if c.position == nil {
			shares := int(c.getFunds() / b.Close)
			if shares > 0 {
				c.broker.Sell(c.instrument, b.Close, shares)
				log.Print("Borrow Shares")
				c.position = NewPosition(BORROWED, shares, 1.05*b.Close)
			}
		} else {
			shares := int(c.getFunds()/b.Close) - c.position.Shares()
			if shares > 0 {
				c.broker.Sell(c.instrument, b.Close, shares)
				log.Print("Borrowed More")
				c.position.AddShares(shares)
			}
		}
	case COVER:
		if (c.position == nil) || (c.position.Type() == BOUGHT) {
			return
		}
		c.broker.Buy(c.instrument, b.Close, c.position.Shares())
		log.Print("Exit Borrowed")
		c.position = nil
	}
}

func (c *Controller) End() {
	if c.position != nil {
		log.Print("Exiting")
		if c.position.Type() == BOUGHT {
			c.broker.PlaceOrder(kt.VarietyRegular, kt.OrderParams{
				Exchange:        c.instrument.Exchange,
				Tradingsymbol:   c.instrument.Tradingsymbol,
				TransactionType: kt.TransactionTypeSell,
				OrderType:       kt.OrderTypeMarket,
				Product:         kt.ProductMIS,
				Validity:        kt.ValidityIOC,
				Quantity:        c.position.Shares(),
			})
			log.Print("Exit Bought")
			c.position = nil
		} else if c.position.Type() == BORROWED {
			c.broker.PlaceOrder(kt.VarietyRegular, kt.OrderParams{
				Exchange:        c.instrument.Exchange,
				Tradingsymbol:   c.instrument.Tradingsymbol,
				TransactionType: kt.TransactionTypeBuy,
				OrderType:       kt.OrderTypeMarket,
				Product:         kt.ProductMIS,
				Validity:        kt.ValidityIOC,
				Quantity:        c.position.Shares(),
			})
			log.Print("Exit Borrowed")
			c.position = nil
		}
	}
}
