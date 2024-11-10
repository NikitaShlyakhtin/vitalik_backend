package order_book

import (
	"vitalik_backend/internal/pkg/types"
)

type OrderBook struct {
	CurrencyPair types.CurrencyPair

	SellOrders []*types.Order
	BuyOrders  []*types.Order
}

func NewOrderBook(currencyPair types.CurrencyPair) (*OrderBook, error) {
	return &OrderBook{
		CurrencyPair: currencyPair,

		SellOrders: make([]*types.Order, 0),
		BuyOrders:  make([]*types.Order, 0),
	}, nil
}
