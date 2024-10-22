package trader

type SuperTrend struct {
	high       *[]float64
	low        *[]float64
	close      *[]float64
	atr        *ATR
	multiplier float64

	fuBand []float64
	flBand []float64
	Data   *[]float64

	processFrom int
}

func NewSuperTrend(bars *Bars, window int, multiplier float64) *SuperTrend {
	return &SuperTrend{
		high:       bars.GetHighSeries(),
		low:        bars.GetLowSeries(),
		close:      bars.GetCloseSeries(),
		atr:        NewATR(bars, window),
		multiplier: multiplier,
		fuBand:     []float64{0},
		flBand:     []float64{0},
		Data:       &[]float64{},
	}
}

func (st *SuperTrend) Update() {
	st.atr.Update()
	atrData := *st.atr.Data
	till := len(atrData)
	if till <= st.processFrom {
		return
	}
	for i := st.processFrom; i < till; i++ {
		h := (*st.high)[len(*st.high)-(till-i)]
		l := (*st.low)[len(*st.low)-(till-i)]
		atr := atrData[i]
		ub := (h+l)/2 + (st.multiplier * atr)
		lb := (h+l)/2 - (st.multiplier * atr)

		pfub := st.fuBand[len(st.fuBand)-1]
		pclose := (*st.close)[len(*st.low)-(till-i)-1]
		fub := pfub
		if ub < pfub || pclose > pfub {
			st.fuBand = append(st.fuBand, ub)
			fub = ub
		} else {
			st.fuBand = append(st.fuBand, pfub)
		}

		pflb := st.flBand[len(st.flBand)-1]
		flb := pflb
		if lb > pflb || pclose < pflb {
			st.flBand = append(st.flBand, lb)
			flb = lb
		} else {
			st.flBand = append(st.flBand, pflb)
		}

		if (*st.close)[len(*st.low)-(till-i)] <= fub {
			*st.Data = append(*st.Data, fub)
		} else {
			*st.Data = append(*st.Data, flb)
		}
	}
	st.processFrom = till
}
