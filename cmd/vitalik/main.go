package main

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"vitalik_backend/internal/app"
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
			app.NewApplication,
			server.NewServer,
		),
		fx.Invoke(startServer),
	)
}
