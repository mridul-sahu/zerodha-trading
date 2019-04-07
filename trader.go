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

func NewPaperTrader(instruments kt.Instruments, broker Broker, feed *Feed) *PaperTrader {
	p := PaperTrader{}
	p.controllers = make(map[uint32]*Controller)
	p.broker = broker
	for i := range instruments {
		id := uint32(instruments[i].InstrumentToken)
		p.controllers[id] = NewController(&instruments[i], feed.GetBars(id), broker)
	}
	p.feed = feed
	return &p
}

func (t *PaperTrader) StartTrading() {
	go func() {
		for {
			select {
			case b := <-t.feed.OnBar:
				t.OnBar(b)
			case <-t.stop:
				return
			}
		}
	}()
}

func (t *PaperTrader) OnBar(b *Bar) {
	t.controllers[b.Instrument].OnBar(b)
}

func (t *PaperTrader) End() {
	t.stop <- true
	for k := range t.controllers {
		t.controllers[k].End()
	}
	log.Print("Stopped")
}
