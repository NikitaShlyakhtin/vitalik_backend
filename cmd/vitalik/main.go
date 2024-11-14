package main

import (
	"go.uber.org/fx"
	"vitalik_backend/internal/fxmodules/loggerfx"
	"vitalik_backend/internal/fxmodules/matcherfx"
	"vitalik_backend/internal/fxmodules/serverfx"
	"vitalik_backend/internal/fxmodules/servicesfx"
	"vitalik_backend/internal/fxmodules/storefx"
)

func main() {
	fx.New(
		loggerfx.Module,
		storefx.Module,
		servicesfx.Module,
		serverfx.Module,
		matcherfx.Module,
	).Run()
}
