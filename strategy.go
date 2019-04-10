package trader

import (
	"log"
)

type Signal int

const (
	BUY Signal = iota
	SELL
	SHORT
	COVER
	HOLD
)

type Strategy struct {
	bars *Bars
}

func NewStrategy(bars *Bars) *Strategy {
	return &Strategy{
		bars: bars,
	}
}

func (s *Strategy) OnBar(b *Bar) Signal {
	log.Println(s.bars.Instrument, s.bars.Len())
	return HOLD
}
