package trader

import (
	"math"
)

type TR struct {
	high        *[]float64
	low         *[]float64
	close       *[]float64
	Data        *[]float64
	processFrom int
}

func NewTR(bars *Bars) *TR {
	return &TR{
		high:        bars.GetHighSeries(),
		close:       bars.GetCloseSeries(),
		low:         bars.GetLowSeries(),
		processFrom: 1,
		Data:        &[]float64{},
	}
}

func (tr *TR) Update() {
	till := len(*tr.high)
	if till <= tr.processFrom {
		return
	}
	for i := tr.processFrom; i < till; i++ {
		*tr.Data = append(*tr.Data, math.Max((*tr.high)[i]-(*tr.low)[i],
			math.Max(math.Abs((*tr.high)[i]-(*tr.close)[i-1]), math.Abs((*tr.low)[i]-(*tr.close)[i-1]))))
	}
	tr.processFrom = till
}
