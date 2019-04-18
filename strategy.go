package trader

type Signal int

const (
	BUY Signal = iota
	SELL
	SHORT
	COVER
	HOLD
)

type StrategyBuilder func(*Bars) Strategy

type Strategy interface {
	OnBar(b *Bar) Signal
	End()
}
