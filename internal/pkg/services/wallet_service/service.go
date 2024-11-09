package wallet_service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
	"vitalik_backend/internal/dependencies"
	walet_service_types "vitalik_backend/internal/pkg/services/wallet_service/types"
	"vitalik_backend/internal/pkg/types"
)

type WalletService struct {
	store dependencies.IStore
}

func NewWalletService(store dependencies.IStore) (*WalletService, error) {
	if store == nil {
		return nil, fmt.Errorf("failed to initialize wallet service")
	}

	return &WalletService{
		store: store,
	}, nil
}

var _ dependencies.IWalletService = (*WalletService)(nil)

func (s *WalletService) CreateWallet(ctx context.Context, args walet_service_types.CreateWalletArgs) (types.Wallet, error) {
	hash := make([]byte, 32)
	_, err := rand.Read(hash)
	if err != nil {
		return types.Wallet{}, fmt.Errorf("failed to generate wallet: %w", err)
	}

	wallet := types.Wallet{
		Requisites: types.Requisites{
			Address: hex.EncodeToString(hash),
			UserID:  args.UserID,
		},
		Currency:  args.Currency,
		Balance:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.store.CreateWallet(ctx, wallet); err != nil {
		return types.Wallet{}, fmt.Errorf("Store CreateWallet failed: %w", err)
	}

	return wallet, nil
}
