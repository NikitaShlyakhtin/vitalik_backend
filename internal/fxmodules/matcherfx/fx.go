package matcherfx

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"vitalik_backend/internal/pkg/services/matcher"
)

func startMatcher(m *matcher.Matcher, lc fx.Lifecycle, l *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			l.Info("Starting matcher...")
			go m.Start(ctx)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			l.Info("Stopping matcher...")
			return nil
		},
	})
}
