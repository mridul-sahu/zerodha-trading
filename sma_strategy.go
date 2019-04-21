package trader

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

type SMAStrategy struct {
	bars *Bars
	sma  *SMA

	startFrom time.Time
}

func NewSMAStrategy(bars *Bars) Strategy {
	return &SMAStrategy{
		bars: bars,
		sma:  NewSMA(10, bars.GetCloseSeries()),
	}
}

func (s *SMAStrategy) OnBar(b *Bar) Signal {
	s.sma.Update()
	if s.startFrom.IsZero() && len(*s.sma.Data) > 0 {
		log.Println("Here: ", b.Datetime)
		s.startFrom = b.Datetime
	}
	return HOLD
}

func (s *SMAStrategy) End() {
	if st, err := json.Marshal(s.sma.Data); err == nil {
		if err := ioutil.WriteFile(strconv.Itoa(int(s.bars.Instrument))+"_sma.json", st, os.ModePerm); err != nil {
			log.Println("Error Writting SMA")
		}
	}
	log.Printf("SMA Values start from: %v\n", s.startFrom)
}
