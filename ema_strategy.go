package trader

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

type EMAStrategy struct {
	bars *Bars
	ema  *EMA

	startFrom time.Time
}

func NewEMAStrategy(bars *Bars) Strategy {
	return &EMAStrategy{
		bars: bars,
		ema:  NewEMA(100, bars.GetCloseSeries()),
	}
}

func (s *EMAStrategy) OnBar(b *Bar) Signal {
	s.ema.Update()
	if s.startFrom.IsZero() && len(*s.ema.Data) > 0 {
		log.Println("Here: ", b.Datetime)
		s.startFrom = b.Datetime
	}
	return HOLD
}

func (s *EMAStrategy) End() {
	if st, err := json.Marshal(s.ema.Data); err == nil {
		if err := ioutil.WriteFile(strconv.Itoa(int(s.bars.Instrument))+"_ema.json", st, os.ModePerm); err != nil {
			log.Println("Error Writting SMA")
		}
	}
	log.Printf("EMA Values start from: %v\n", s.startFrom)
}
