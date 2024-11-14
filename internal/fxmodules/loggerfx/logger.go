package loggerfx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("loggerfx",
	fx.WithLogger(getEventLogger),
	fx.Provide(
		zap.NewDevelopment,
	),
)
