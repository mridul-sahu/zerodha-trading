package trader

import (
	"log"

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
	if c.position != nil && c.position.StoplossHit(b) {
		if c.position.Type() == BOUGHT {
			// EXIT BOUGHT
			log.Print("Exit Bought")
			c.position = nil
		} else if c.position.Type() == BORROWED {
			// EXIT BORROWED
			log.Print("Exit Borrowed")
			c.position = nil
		}
	}

	signal := c.strategy.OnBar(b)
	switch signal {
	case HOLD:
		log.Print("Hold")
	case BUY:
		if (c.position != nil) && (c.position.Type() == BORROWED) {
			// EXIT BORROWED
			log.Print("Exit Borrowed")
			c.position = nil
		}
		if c.position == nil {
			shares := int(c.getFunds() / b.Close)
			// BUY SHARES
			log.Print("Buy Shares")
			c.position = NewPosition(BOUGHT, shares, 0.95*b.Close)
		} else {
			//TODO: (Maybe Buy More)
		}
	case SELL:
		if (c.position == nil) || (c.position.Type() == BORROWED) {
			return
		}
		// EXIT BOUGHT
		log.Print("Exit Bought")
		c.position = nil
	case SHORT:
		if (c.position != nil) && (c.position.Type() == BOUGHT) {
			// EXIT BOUGHT
			log.Print("Exit Bought")
			c.position = nil
		}
		if c.position == nil {
			shares := int(c.getFunds() / b.Close)
			// Bowrrow SHARES
			log.Print("Borrow Shares")
			c.position = NewPosition(BORROWED, shares, 1.05*b.Close)
		} else {
			//TODO: (Maybe Borrow More)
		}
	case COVER:
		if (c.position == nil) || (c.position.Type() == BOUGHT) {
			return
		}
		// EXIT BORROWED
		log.Print("Exit Borrowed")
		c.position = nil
	}

	// MAYBE UPDATE STOPLOSS
}
