package technical

type ATR struct {
	tr          *TR
	processFrom int
	Window      int
	Data        []float64
}

func NewATR(bars *Bars, window int) *ATR {
	return &ATR{
		tr:          NewTR(bars),
		Window:      window,
		processFrom: window - 1,
	}
}

func (atr *ATR) Update() {
	atr.tr.Update()
	trData := atr.tr.Data
	till := len(trData)
	if till <= atr.processFrom {
		return
	}
	for i := atr.processFrom; i < till; i++ {
		if len(atr.Data) > 0 {
			currentTr := trData[i]
			lastVal := atr.Data[len(atr.Data)-1]
			newVal := lastVal + (currentTr-lastVal)/float64(atr.Window)
			atr.Data = append(atr.Data, newVal)
		} else {
			var sum float64
			for j := 0; j < atr.Window; j++ {
				sum += trData[j]
			}
			atr.Data = append(atr.Data, sum/float64(atr.Window))
		}
	}
	atr.processFrom = till
}
