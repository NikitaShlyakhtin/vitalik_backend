package servicesfx

import (
	"go.uber.org/fx"
	"vitalik_backend/internal/app"
	"vitalik_backend/internal/dependencies"
	"vitalik_backend/internal/pkg/services/order_book_manager"
	"vitalik_backend/internal/pkg/services/wallet_service"
)

var Module = fx.Module("servicesfx",
	fx.Provide(
		fx.Annotate(app.NewApplication, fx.As(new(dependencies.IHandler))),
		fx.Annotate(wallet_service.NewWalletService, fx.As(new(dependencies.IWalletService))),
		fx.Annotate(order_book_manager.NewOrderBookManager, fx.As(new(dependencies.IOrderBookManager))),
	),
)
