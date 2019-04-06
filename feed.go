package trader

import (
	"log"
	"time"

	"github.com/zerodhatech/gokiteconnect/ticker"
)

type Bar struct {
	Datetime time.Time
	Open     float64
	High     float64
	Low      float64
	Close    float64
}

type Bars struct {
	Instrument uint32
	dates      []time.Time
	opens      []float64
	highs      []float64
	lows       []float64
	closes     []float64
}

func (b *Bars) getDates() []time.Time {
	return b.dates
}

func (b *Bars) getOpenSeries() []float64 {
	return b.opens
}

func (b *Bars) getHighSeries() []float64 {
	return b.highs
}

func (b *Bars) getLowSeries() []float64 {
	return b.lows
}

func (b *Bars) getCloseSeries() []float64 {
	return b.closes
}

func (b *Bars) size() int {
	return len(b.dates)
}

func (b *Bars) AddBar(bar *Bar) {
	b.dates = append(b.dates, bar.Datetime)
	b.lows = append(b.lows, bar.Low)
	b.highs = append(b.highs, bar.High)
	b.lows = append(b.lows, bar.Low)
	b.closes = append(b.closes, bar.Close)
}

func (b *Bars) AddBars(bars []*Bar) {
	n := b.size()

	dates := make([]time.Time, n, len(bars)+n)
	opens := make([]float64, n, len(bars)+n)
	closes := make([]float64, n, len(bars)+n)
	highs := make([]float64, n, len(bars)+n)
	lows := make([]float64, n, len(bars)+n)

	copy(dates, b.dates)
	copy(opens, b.opens)
	copy(closes, b.closes)
	copy(highs, b.highs)
	copy(lows, b.lows)

	for i, v := range bars {
		dates[n+i] = v.Datetime
		opens[n+1] = v.Open
		closes[n+1] = v.Close
		highs[n+1] = v.High
		lows[n+1] = v.Low
	}

	b.dates = dates
	b.opens = opens
	b.closes = closes
	b.highs = highs
	b.lows = lows
}

type Feed struct {
	data  map[uint32]*Bars
	OnBar chan *Bar
}

func NewFeed(instruments []uint32) *Feed {
	f := Feed{}
	f.data = make(map[uint32]*Bars)
	f.OnBar = make(chan *Bar)

	for _, inst := range instruments {
		f.data[inst] = &Bars{Instrument: inst}
	}

	return &f
}

func (f *Feed) onTick(tick kiteticker.Tick) {
	ohlc := tick.OHLC
	bar := Bar{
		Open:     ohlc.Open,
		Close:    ohlc.Close,
		High:     ohlc.High,
		Low:      ohlc.Low,
		Datetime: tick.Timestamp.Time,
	}
	f.data[tick.InstrumentToken].AddBar(&bar)

	select {
	case f.OnBar <- &bar:
	default:
		log.Printf("Tick Dropped: %v", tick.InstrumentToken)
	}
}
