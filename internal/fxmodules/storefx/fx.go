package storefx

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func newPgxPool() (*pgxpool.Pool, error) {
	dsn := "postgres://postgres:password@localhost:5452/vitalik?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func startPgxPool(lc fx.Lifecycle, pool *pgxpool.Pool, l *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			l.Info("Connecting to PostgreSQL...")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			l.Info("Closing PostgreSQL connection pool...")
			pool.Close()
			return nil
		},
	})
}
