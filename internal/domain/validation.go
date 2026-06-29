package domain

import (
	"errors"
	"strings"
)

var (
	ErrServiceNameRequired  = errors.New("service name is required")
	ErrPriceMustBePositive  = errors.New("price must be positive")
	ErrInvalidPeriod        = errors.New("end month must be greater than or equal to start month")
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

func ValidateSubscription(sub Subscription) error {
	if strings.TrimSpace(sub.ServiceName) == "" {
		return ErrServiceNameRequired
	}

	if sub.Price <= 0 {
		return ErrPriceMustBePositive
	}

	if sub.EndMonth != nil && sub.EndMonth.Before(sub.StartMonth) {
		return ErrInvalidPeriod
	}

	return nil
}
