package trader

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mridul-sahu/zerodha-trading/technical"
)

type SuperTrendStrategy struct {
	bars *Bars
	st   *technical.SuperTrend

	startFrom time.Time
}

func NewSuperTrendStrategy(bars *Bars) Strategy {
	return &SuperTrendStrategy{
		bars: bars,
		st:   NewSuperTrend(bars, 10, 3),
	}
}

func (s *SuperTrendStrategy) OnBar(b *Bar) Signal {
	s.st.Update()
	if s.startFrom.IsZero() && len(s.st.Data) > 0 {
		log.Println("Here: ", b.Datetime)
		s.startFrom = b.Datetime
	}
	return HOLD
}

func (s *SuperTrendStrategy) End() {
	if st, err := json.Marshal(s.st.Data); err == nil {
		if err := ioutil.WriteFile(strconv.Itoa(int(s.bars.Instrument))+"_supertrend.json", st, os.ModePerm); err != nil {
			log.Println("Error Writting SuperTrend")
		}
	}
	log.Printf("SuperTrend Values start from: %v\n", s.startFrom)
}
