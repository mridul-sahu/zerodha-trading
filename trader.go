package trader

import (
	"log"

	kt "github.com/zerodhatech/gokiteconnect"
)

type PaperTrader struct {
	feed        *Feed
	broker      Broker
	controllers map[uint32]*Controller
	stop        chan bool
}

func NewPaperTrader(instruments kt.Instruments, broker Broker, feed *Feed, sb StrategyBuilder) *PaperTrader {
	p := PaperTrader{}
	p.controllers = make(map[uint32]*Controller)
	p.broker = broker
	for i := range instruments {
		id := uint32(instruments[i].InstrumentToken)
		p.controllers[id] = NewController(&instruments[i], feed.GetBars(id), broker, sb)
	}
	p.feed = feed
	p.stop = make(chan bool)
	return &p
}

func (t *PaperTrader) StartTrading() {
	for k := range t.controllers {
		go func(id uint32, c *Controller, b <-chan *Bar) {
			for {
				select {
				case bar := <-b:
					c.OnBar(bar)
				case <-t.stop:
					log.Println("Stopped")
					return
				}
			}
		}(k, t.controllers[k], t.feed.OnBar[k])
	}
}

func (t *PaperTrader) End() {
	log.Print("Stopping")
	t.stop <- true
	for k := range t.controllers {
		t.controllers[k].End()
	}
}
