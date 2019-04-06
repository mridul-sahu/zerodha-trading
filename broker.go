package trader

import (
	"strconv"

	kt "github.com/zerodhatech/gokiteconnect"
)

type Broker interface {
	PlaceOrder(orderParams kt.OrderParams) (string, error)
}

type PaperBroker struct {
	orders []kt.OrderParams
}

func (pb *PaperBroker) GetOrders() []kt.OrderParams {
	return pb.orders
}

func (pb *PaperBroker) PlaceOrders(orderParams kt.OrderParams) (string, error) {
	pb.orders = append(pb.orders, orderParams)
	return strconv.Itoa(len(pb.orders)), nil
}
