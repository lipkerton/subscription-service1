package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          int64
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartMonth  time.Time
	EndMonth    *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SubscriptionFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
}

type SummaryFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	FromMonth   time.Time
	ToMonth     time.Time
}
