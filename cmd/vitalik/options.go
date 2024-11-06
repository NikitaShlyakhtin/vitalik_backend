package main

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"vitalik_backend/internal/server"
)

func startServer(lc fx.Lifecycle, s *server.Server, l *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				port := 8080
				if err := s.Start(fmt.Sprintf(":%d", port)); err != nil {
					l.Error(err.Error())
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.Shutdown(ctx)
		},
	})
}

func getEventLogger(l *zap.Logger) fxevent.Logger {
	return &fxevent.ZapLogger{Logger: l}
}
