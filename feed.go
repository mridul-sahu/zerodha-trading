package trader

import (
	"log"

	"github.com/zerodhatech/gokiteconnect/ticker"
)

type Feed struct {
	data  map[uint32]*Bars
	OnBar map[uint32]chan *Bar
}

func NewFeed(instruments []uint32) *Feed {
	f := Feed{}
	f.data = make(map[uint32]*Bars)
	f.OnBar = make(map[uint32]chan *Bar, 10)

	for _, inst := range instruments {
		f.data[inst] = NewBars(inst)
		f.OnBar[inst] = make(chan *Bar)
	}

	return &f
}

func (f *Feed) GetBars(instrument uint32) *Bars {
	return f.data[instrument]
}

func (f *Feed) OnTick(tick kiteticker.Tick) {
	ohlc := tick.OHLC
	bar := Bar{
		Open:     ohlc.Open,
		Close:    tick.LastPrice,
		High:     ohlc.High,
		Low:      ohlc.Low,
		Datetime: tick.Timestamp.Time,
	}
	f.data[tick.InstrumentToken].AddBar(&bar)
	select {
	case f.OnBar[tick.InstrumentToken] <- &bar:
		log.Printf("Received: %v", tick.InstrumentToken)
	default:
		log.Printf("Tick Dropped: %v", tick.InstrumentToken)
	}
}
