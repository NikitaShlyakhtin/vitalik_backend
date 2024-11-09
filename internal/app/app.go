package app

import (
	"errors"
	"go.uber.org/zap"
	"vitalik_backend/internal/dependencies"
)

// Application holds the application state and dependencies
type Application struct {
	Logger        *zap.SugaredLogger
	WalletService dependencies.IWalletService
	OrderBook     dependencies.IOrderBook
	Store         dependencies.IStore
}

// NewApplication initializes a new Application instance
func NewApplication(
	logger *zap.Logger,
	walletService dependencies.IWalletService,
	orderBook dependencies.IOrderBook,
	store dependencies.IStore,
) (*Application, error) {
	if logger == nil ||
		walletService == nil ||
		orderBook == nil ||
		store == nil {
		return nil, errors.New("failed to initialize application")
	}

	return &Application{
		Logger:        logger.Sugar(),
		WalletService: walletService,
		OrderBook:     orderBook,
		Store:         store,
	}, nil
}

var _ dependencies.IHandler = (*Application)(nil)
