package app

import (
	"errors"
	"go.uber.org/zap"
	"vitalik_backend/internal/dependencies"
)

// Application holds the application state and dependencies
type Application struct {
	Logger           *zap.SugaredLogger
	WalletService    dependencies.IWalletService
	OrderBookManager dependencies.IOrderBookManager
	Store            dependencies.IStore
	AuthService      dependencies.IAuthService
}

// NewApplication initializes a new Application instance
func NewApplication(
	logger *zap.Logger,
	walletService dependencies.IWalletService,
	orderBookManager dependencies.IOrderBookManager,
	store dependencies.IStore,
	authService dependencies.IAuthService,
) (*Application, error) {
	if logger == nil ||
		walletService == nil ||
		orderBookManager == nil ||
		store == nil ||
		authService == nil {
		return nil, errors.New("failed to initialize application")
	}

	return &Application{
		Logger:           logger.Sugar(),
		WalletService:    walletService,
		OrderBookManager: orderBookManager,
		Store:            store,
		AuthService:      authService,
	}, nil
}

var _ dependencies.IHandler = (*Application)(nil)
