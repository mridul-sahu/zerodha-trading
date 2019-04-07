package trader

import "time"

type Bar struct {
	Instrument uint32
	Datetime   time.Time
	Open       float64
	High       float64
	Low        float64
	Close      float64
}

type Bars struct {
	Instrument uint32
	dates      []time.Time
	opens      []float64
	highs      []float64
	lows       []float64
	closes     []float64
}

func (b *Bars) GetDates() []time.Time {
	return b.dates
}

func (b *Bars) GetOpenSeries() []float64 {
	return b.opens
}

func (b *Bars) GetHighSeries() []float64 {
	return b.highs
}

func (b *Bars) GetLowSeries() []float64 {
	return b.lows
}

func (b *Bars) GetCloseSeries() []float64 {
	return b.closes
}

func (b *Bars) Len() int {
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
	n := b.Len()

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
