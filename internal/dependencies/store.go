package dependencies

import (
	"context"
	store_types "vitalik_backend/internal/pkg/services/store/types"
	"vitalik_backend/internal/pkg/types"
)

type IStore interface {
	CreateWallet(ctx context.Context, wallet types.Wallet) error
	ListWallets(ctx context.Context, args store_types.ListWalletsArgs) ([]*types.Wallet, error)

	Deposit(ctx context.Context, args store_types.DepositArgs) (*types.Transaction, error)

	ListTransactions(ctx context.Context, args store_types.ListTransactionsArgs) ([]*types.Transaction, error)
	Transfer(ctx context.Context, args store_types.TransferArgs) (*types.Transaction, error)

	SaveOrder(ctx context.Context, order types.Order) error
}
