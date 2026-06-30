package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lipkerton/subscription-service1/internal/domain"
)

type SubscriptionRepository struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{
		db: db,
	}
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	const query = `
		INSERT INTO subscriptions (
			service_name,
			price,
			user_id,
			start_month,
			end_month
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING 
			id,
			service_name,
			price,
			user_id,
			start_month,
			end_month,
			created_at,
			updated_at
	`

	var created domain.Subscription

	err := r.db.QueryRow(
		ctx,
		query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartMonth,
		sub.EndMonth,
	).Scan(
		&created.ID,
		&created.ServiceName,
		&created.Price,
		&created.UserID,
		&created.StartMonth,
		&created.EndMonth,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("insert subscription: %w", err)
	}

	return created, nil
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id int64) (domain.Subscription, error) {
	const query = `
		SELECT
			id,
			service_name,
			price,
			user_id,
			start_month,
			end_month,
			created_at,
			updated_at
		FROM subscriptions
		WHERE id = $1
	`

	var sub domain.Subscription

	err := r.db.QueryRow(ctx, query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartMonth,
		&sub.EndMonth,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Subscription{}, domain.ErrSubscriptionNotFound
		}

		return domain.Subscription{}, fmt.Errorf("select subscription by id: %w", err)
	}

	return sub, nil
}
