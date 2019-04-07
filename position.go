package trader

import "errors"

type PositionType int

const (
	BOUGHT PositionType = iota
	BORROWED
)

type Position struct {
	shares        int
	tp            PositionType
	StoplossPrice float64
}

func NewPosition(tp PositionType, shares int, stoplossPrice float64) *Position {
	return &Position{
		shares:        shares,
		tp:            tp,
		StoplossPrice: stoplossPrice,
	}
}

func (p *Position) Type() PositionType {
	return p.tp
}

func (p *Position) Shares() int {
	return p.shares
}

func (p *Position) AddShares(s int) {
	p.shares += s
}

func (p *Position) RemoveShares(s int) error {
	if p.shares < s {
		return errors.New("Cannot remove more shares than added")
	}
	p.shares -= s
	return nil
}

func (p *Position) StoplossHit(b *Bar) bool {
	if (p.tp == BOUGHT && b.Low <= p.StoplossPrice) || ((p.tp == BORROWED) && b.High >= p.StoplossPrice) {
		return true
	}
	return false
}
