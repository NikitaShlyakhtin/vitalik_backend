package main

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"vitalik_backend/internal/app"
	"vitalik_backend/internal/dependencies"
	"vitalik_backend/internal/pkg/services/order_book_manager"
	"vitalik_backend/internal/pkg/services/store"
	"vitalik_backend/internal/pkg/services/wallet_service"
	"vitalik_backend/internal/server"
)

func main() {
	fx.New(buildFxOptions()).
		Run()
}

func buildFxOptions() fx.Option {
	return fx.Options(
		fx.WithLogger(getEventLogger),
		fx.Provide(
			zap.NewDevelopment,

			fx.Annotate(app.NewApplication, fx.As(new(dependencies.IHandler))),
			fx.Annotate(wallet_service.NewWalletService, fx.As(new(dependencies.IWalletService))),
			fx.Annotate(order_book_manager.NewOrderBookManager, fx.As(new(dependencies.IOrderBookManager))),
			fx.Annotate(store.NewStore, fx.As(new(dependencies.IStore))),

			server.NewServer,
		),
		fx.Invoke(startServer),
	)
}
