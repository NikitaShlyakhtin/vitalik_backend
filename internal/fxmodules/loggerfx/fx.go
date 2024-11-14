package loggerfx

import (
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func getEventLogger(l *zap.Logger) fxevent.Logger {
	return &fxevent.ZapLogger{Logger: l}
}
