package technical

import (
	"math"

	trader "github.com/mridul-sahu/zerodha-trading"
)

type TR struct {
	bars        *trader.Bars
	Data        []float64
	processFrom int
}

func NewTR(bars *Bars) *TR {
	return &TR{bars: bars, processFrom: 1}
}

func (tr *TR) Update() {
	close := tr.bars.GetCloseSeries()
	high := tr.bars.GetHighSeries()
	low := tr.bars.GetLowSeries()

	till := len(high)
	if till <= tr.processFrom {
		return
	}
	for i := tr.processFrom; i < till; i++ {
		tr.Data = append(tr.Data, math.Max(high[i]-low[i],
			math.Max(math.Abs(high[i]-close[i-1]), math.Abs(low[i]-close[i-1]))))
	}
	tr.processFrom = till
}
