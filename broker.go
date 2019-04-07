package trader

import (
	"strconv"

	kt "github.com/zerodhatech/gokiteconnect"
)

type Broker interface {
	PlaceOrder(variety string, orderParams kt.OrderParams) (string, error)
	GetAvailableFunds() float64
}

type PaperBroker struct {
	orders []kt.OrderParams
	funds  float64
}

func (pb *PaperBroker) GetOrders() []kt.OrderParams {
	return pb.orders
}

func (pb *PaperBroker) PlaceOrder(variety string, orderParams kt.OrderParams) (string, error) {
	pb.orders = append(pb.orders, orderParams)
	return strconv.Itoa(len(pb.orders)), nil
}

func (pb *PaperBroker) Buy(inst kt.Instrument, price float64, qunatity int) (string, error) {
	return pb.PlaceOrder(kt.VarietyRegular, kt.OrderParams{
		Exchange:        inst.Exchange,
		Tradingsymbol:   inst.Tradingsymbol,
		TransactionType: kt.TransactionTypeBuy,
		OrderType:       kt.OrderTypeLimit,
		Product:         kt.ProductMIS,
		Validity:        kt.ValidityIOC,
		Quantity:        qunatity,
		Price:           price,
	})
}

func (pb *PaperBroker) Sell(inst kt.Instrument, price float64, qunatity int) (string, error) {
	return pb.PlaceOrder(kt.VarietyRegular, kt.OrderParams{
		Exchange:        inst.Exchange,
		Tradingsymbol:   inst.Tradingsymbol,
		TransactionType: kt.TransactionTypeSell,
		OrderType:       kt.OrderTypeLimit,
		Product:         kt.ProductMIS,
		Validity:        kt.ValidityIOC,
		Quantity:        qunatity,
		Price:           price,
	})
}

func (pb *PaperBroker) GetAvailableFunds() float64 {
	return pb.funds
}
