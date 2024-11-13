package types

import (
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"time"
)

type OrderType string

const (
	Buy  OrderType = "BUY"
	Sell OrderType = "SELL"
)

type OrderStatus string

const (
	OrderOpen      OrderStatus = "ORDER_OPEN"
	OrderClosed    OrderStatus = "ORDER_CLOSED"
	OrderCancelled OrderStatus = "ORDER_CANCELLED"
)

type Order struct {
	ID uuid.UUID `json:"id"`

	Type OrderType `json:"type"`

	SellCurrency   Currency   `json:"sell_currency"`
	SellQuantity   null.Float `json:"sell_quantity"`
	SellRequisites Requisites `json:"sell_requisites"`

	Price float64 `json:"price"`

	BuyCurrency   Currency   `json:"buy_currency"`
	BuyQuantity   null.Float `json:"buy_quantity"`
	BuyRequisites Requisites `json:"buy_requisites"`

	Status OrderStatus `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	RemovedAt null.Time `json:"removed_at,omitempty"`
	ClosedAt  null.Time `json:"closed_at,omitempty"`
}
