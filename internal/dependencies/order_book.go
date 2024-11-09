package dependencies

import (
	"context"
	"github.com/google/uuid"
	order_book_types "vitalik_backend/internal/pkg/services/order_book/types"
	"vitalik_backend/internal/pkg/types"
)

// IOrderBook defines the methods required for managing and interacting with an order book
type IOrderBook interface {
	CreateOrder(ctx context.Context, args order_book_types.CreateOrderArgs) (*types.Order, error)
	MatchOrders(ctx context.Context) error
	CancelOrder(ctx context.Context, orderID uuid.UUID) error
	ListOrders(ctx context.Context, args order_book_types.ListOrdersArgs) ([]*types.Order, error)
}
