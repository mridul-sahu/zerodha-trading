package trader

type SMA struct {
	series      *[]float64
	Window      int
	Data        *[]float64
	processFrom int
}

func NewSMA(window int, series *[]float64) *SMA {
	return &SMA{
		series:      series,
		Window:      window,
		Data:        &[]float64{},
		processFrom: window - 1,
	}
}

func (sm *SMA) Update() {
	till := len(*sm.series)
	if till <= sm.processFrom {
		return
	}
	for i := sm.processFrom; i < till; i++ {
		if len(*sm.Data) > 0 {
			lastVal := (*sm.Data)[len(*sm.Data)-1]
			delta := ((*sm.series)[i] - (*sm.series)[i-sm.Window]) / float64(sm.Window)
			*sm.Data = append(*sm.Data, lastVal+delta)
		} else {
			var sum float64
			for j := 0; j < sm.Window; j++ {
				sum += (*sm.series)[j]
			}
			*sm.Data = append(*sm.Data, sum/float64(sm.Window))
		}
	}
	sm.processFrom = till
}
