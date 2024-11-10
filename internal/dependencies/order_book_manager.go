package dependencies

import (
	"context"
	"github.com/google/uuid"
	order_book_types "vitalik_backend/internal/pkg/services/order_book/types"
	"vitalik_backend/internal/pkg/types"
)

// IOrderBookManager defines methods for managing multiple order books.
type IOrderBookManager interface {
	ListAvailableCurrencyPairs(ctx context.Context) ([]types.CurrencyPair, error)
	MatchOrders(ctx context.Context) error

	CreateOrder(ctx context.Context, args order_book_types.CreateOrderArgs) (*types.Order, error)
	CancelOrder(ctx context.Context, currencyPair types.CurrencyPair, orderID uuid.UUID) error
	ListOrders(ctx context.Context, args order_book_types.ListOrdersArgs) ([]*types.Order, error)
}
