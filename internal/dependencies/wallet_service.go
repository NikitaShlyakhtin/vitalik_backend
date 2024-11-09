package dependencies

import (
	"context"
	walet_service_types "vitalik_backend/internal/pkg/services/wallet_service/types"
	"vitalik_backend/internal/pkg/types"
)

type IWalletService interface {
	CreateWallet(ctx context.Context, args walet_service_types.CreateWalletArgs) (types.Wallet, error)
}
