package matcher

import (
	"context"
	"go.uber.org/zap"
	"time"
	"vitalik_backend/internal/dependencies"
)

type Matcher struct {
	orderBookManager dependencies.IOrderBookManager
	logger           *zap.Logger
	interval         time.Duration
}

func NewMatcher(orderBookManager dependencies.IOrderBookManager, logger *zap.Logger) *Matcher {
	return &Matcher{
		orderBookManager: orderBookManager,
		logger:           logger,
		interval:         3 * time.Second,
	}
}

func (m *Matcher) Start(ctx context.Context) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	m.logger.Info("Matcher started, will call handler every", zap.Duration("interval", m.interval))

	ctx = context.Background()

	for {
		select {
		case <-ticker.C:
			err := m.orderBookManager.MatchOrders(ctx)
			if err != nil {
				m.logger.Error("Error calling handler", zap.Error(err))
			}
		case <-ctx.Done():
			m.logger.Info("Worker stopped")
			return
		}

	}
}
