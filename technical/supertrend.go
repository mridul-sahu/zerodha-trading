package technical

import trader "github.com/mridul-sahu/zerodha-trading"

type SuperTrend struct {
	bars       *trader.Bars
	atr        *ATR
	multiplier float64

	fuBand []float64
	flBand []float64
	Data   []float64

	processFrom int
}

func NewSuperTrend(bars *Bars, window int, multiplier float64) *SuperTrend {
	return &SuperTrend{
		bars:       bars,
		atr:        NewATR(bars, window),
		multiplier: multiplier,
		fuBand:     []float64{0},
		flBand:     []float64{0},
	}
}

func (st *SuperTrend) Update() {
	st.atr.Update()
	high := st.bars.GetHighSeries()
	low := st.bars.GetCloseSeries()
	close := st.bars.GetCloseSeries()

	atrData := st.atr.Data
	till := len(atrData)
	if till <= st.processFrom {
		return
	}
	for i := st.processFrom; i < till; i++ {
		h := high[len(high)-(till-i)]
		l := low[len(low)-(till-i)]
		atr := atrData[i]
		ub := (h+l)/2 + (st.multiplier * atr)
		lb := (h+l)/2 - (st.multiplier * atr)

		pfub := st.fuBand[len(st.fuBand)-1]
		pclose := close[len(low)-(till-i)-1]
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

		if close[len(low)-(till-i)] <= fub {
			st.Data = append(st.Data, fub)
		} else {
			st.Data = append(st.Data, flb)
		}
	}
	st.processFrom = till
}
