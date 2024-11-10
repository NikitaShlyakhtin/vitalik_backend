package order_book_types

import (
	"github.com/guregu/null/v5"
	"time"
	"vitalik_backend/internal/pkg/types"
)

type CreateOrderArgs struct {
	Type types.OrderType

	SellCurrency   types.Currency
	SellQuantity   null.Float
	SellRequisites types.Requisites

	Price float64

	BuyCurrency   types.Currency
	BuyQuantity   null.Float
	BuyRequisites types.Requisites

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListOrdersArgs struct {
	CurrencyPair types.CurrencyPair

	UserIDIn    []string
	OrderType   null.Value[types.OrderType]
	OrderStatus null.Value[types.OrderStatus]
}
