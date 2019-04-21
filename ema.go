package trader

type EMA struct {
	series      *[]float64
	Window      int
	Data        *[]float64
	processFrom int
}

func NewEMA(window int, series *[]float64) *EMA {
	return &EMA{
		series:      series,
		Window:      window,
		Data:        &[]float64{},
		processFrom: window - 1,
	}
}

func (em *EMA) Update() {
	till := len(*em.series)
	if till <= em.processFrom {
		return
	}
	for i := em.processFrom; i < till; i++ {
		if len(*em.Data) > 0 {
			lastVal := (*em.Data)[len(*em.Data)-1]
			newVal := lastVal + ((*em.series)[i]-lastVal)/float64(em.Window)
			*em.Data = append(*em.Data, newVal)
		} else {
			var sum float64
			for j := 0; j < em.Window; j++ {
				sum += (*em.series)[j]
			}
			*em.Data = append(*em.Data, sum/float64(em.Window))
		}
	}
	em.processFrom = till
}
