package store_types

import (
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"time"
	"vitalik_backend/internal/pkg/types"
)

type Order struct {
	ID uuid.UUID `db:"orders.id"`

	Type string `db:"orders.type"`

	SellCurrency string     `db:"orders.sell_currency"`
	SellQuantity null.Float `db:"orders.sell_quantity"`
	SellAddress  string     `db:"orders.sell_address"`
	SellUserID   string     `db:"orders.sell_user_id"`

	Price float64 `db:"orders.price"`

	BuyCurrency string     `db:"orders.buy_currency"`
	BuyQuantity null.Float `db:"orders.buy_quantity"`
	BuyAddress  string     `db:"orders.buy_address"`
	BuyUserID   string     `db:"orders.buy_user_id"`

	Status string `db:"orders.status"`

	CreatedAt time.Time `db:"orders.created_at"`
	UpdatedAt time.Time `db:"orders.updated_at"`
	RemovedAt null.Time `db:"orders.removed_at"`
	ClosedAt  null.Time `db:"orders.closed_at"`
}

func MapToOrderStore(order *types.Order) *Order {
	return &Order{
		ID:           order.ID,
		Type:         string(order.Type),
		SellCurrency: string(order.SellCurrency),
		SellQuantity: order.SellQuantity,
		SellAddress:  order.SellRequisites.Address,
		SellUserID:   order.SellRequisites.UserID,
		Price:        order.Price,
		BuyCurrency:  string(order.BuyCurrency),
		BuyQuantity:  order.BuyQuantity,
		BuyAddress:   order.BuyRequisites.Address,
		BuyUserID:    order.BuyRequisites.UserID,
		Status:       string(order.Status),
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
		RemovedAt:    order.RemovedAt,
		ClosedAt:     order.ClosedAt,
	}
}

func MapToOrder(orderStore *Order) *types.Order {
	return &types.Order{
		ID:           orderStore.ID,
		Type:         types.OrderType(orderStore.Type),
		SellCurrency: types.Currency(orderStore.SellCurrency),
		SellQuantity: orderStore.SellQuantity,
		SellRequisites: types.Requisites{
			Address: orderStore.SellAddress,
			UserID:  orderStore.SellUserID,
		},
		Price:       orderStore.Price,
		BuyCurrency: types.Currency(orderStore.BuyCurrency),
		BuyQuantity: orderStore.BuyQuantity,
		BuyRequisites: types.Requisites{
			Address: orderStore.BuyAddress,
			UserID:  orderStore.BuyUserID,
		},
		Status:    types.OrderStatus(orderStore.Status),
		CreatedAt: orderStore.CreatedAt,
		UpdatedAt: orderStore.UpdatedAt,
		RemovedAt: orderStore.RemovedAt,
		ClosedAt:  orderStore.ClosedAt,
	}
}
