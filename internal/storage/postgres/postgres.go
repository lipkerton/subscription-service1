package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lipkerton/subscription-service1/internal/domain"
)

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 1
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return pool, nil
}

func (r *SubscriptionRepository) List(ctx context.Context, filter domain.SubscriptionFilter) ([]domain.Subscription, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *SubscriptionRepository) CalculateSummary(ctx context.Context, filter domain.SummaryFilter) (int, error) {
	return 0, fmt.Errorf("not implemented")
}
