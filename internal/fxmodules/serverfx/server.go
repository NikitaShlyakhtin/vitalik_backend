package serverfx

import (
	"go.uber.org/fx"
	"vitalik_backend/internal/server"
)

var Module = fx.Module("serverfx",
	fx.Provide(server.NewServer),
	fx.Invoke(startServer),
)
