package store_types

import (
	"time"
	"vitalik_backend/internal/pkg/types"
)

type Wallet struct {
	Address string `db:"wallets.address"`
	UserID  string `db:"wallets.user_id"`

	Currency string  `db:"wallets.currency"`
	Balance  float64 `db:"wallets.balance"`

	CreatedAt time.Time `db:"wallets.created_at"`
	UpdatedAt time.Time `db:"wallets.updated_at"`
}

// MapToWalletStore converts a Wallet struct to a WalletStore struct.
func MapToWalletStore(wallet *types.Wallet) *Wallet {
	return &Wallet{
		Address:   wallet.Requisites.Address,
		UserID:    wallet.Requisites.UserID,
		Currency:  string(wallet.Currency),
		Balance:   wallet.Balance,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}
}

// MapToWallet converts a WalletStore struct to a Wallet struct.
func MapToWallet(walletStore Wallet) *types.Wallet {
	return &types.Wallet{
		Requisites: types.Requisites{
			Address: walletStore.Address,
			UserID:  walletStore.UserID,
		},
		Currency:  types.Currency(walletStore.Currency),
		Balance:   walletStore.Balance,
		CreatedAt: walletStore.CreatedAt,
		UpdatedAt: walletStore.UpdatedAt,
	}
}
