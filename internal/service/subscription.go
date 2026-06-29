package service

import (
	"context"
	"fmt"

	"github.com/lipkerton/subscription-service1/internal/domain"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error)
	GetByID(ctx context.Context, id int64) (domain.Subscription, error)
	Update(ctx context.Context, sub domain.Subscription) (domain.Subscription, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filter domain.SubscriptionFilter) ([]domain.Subscription, error)
	CalculateSummary(ctx context.Context, filter domain.SummaryFilter) (int, error)
}

type SubscriptionService struct {
	repo SubscriptionRepository
}

func NewSubscriptionService(repo SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
	}
}

func (s *SubscriptionService) Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	if err := domain.ValidateSubscription(sub); err != nil {
		return domain.Subscription{}, fmt.Errorf("create subscription: %w", err)
	}
	createdSub, err := s.repo.Create(ctx, sub)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("create subscription: %w", err)
	}
	return createdSub, nil
}

func (s *SubscriptionService) GetByID(ctx context.Context, id int64) (domain.Subscription, error) {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("get subscription by id: %w", err)
	}
	return sub, nil
}

func (s *SubscriptionService) Update(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	if err := domain.ValidateSubscription(sub); err != nil {
		return domain.Subscription{}, fmt.Errorf("update subscription: %w", err)
	}
	updatedSub, err := s.repo.Update(ctx, sub)
	if err != nil {
		return domain.Subscription{}, fmt.Errorf("update subscription: %w", err)
	}
	return updatedSub, nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete subsctiption: %w", err)
	}
	return nil
}

func (s *SubscriptionService) List(ctx context.Context, filter domain.SubscriptionFilter) ([]domain.Subscription, error) {
	subs, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}

	return subs, nil
}

func (s *SubscriptionService) CalculateSummary(ctx context.Context, filter domain.SummaryFilter) (int, error) {
	if filter.ToMonth.Before(filter.FromMonth) {
		return 0, domain.ErrInvalidPeriod
	}

	total, err := s.repo.CalculateSummary(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("calculate summary: %w", err)
	}

	return total, nil
}
