package wallet_service_types

import "vitalik_backend/internal/pkg/types"

type CreateWalletArgs struct {
	UserID   string
	Currency types.Currency
}
