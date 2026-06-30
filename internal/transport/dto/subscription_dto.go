package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/lipkerton/subscription-service1/internal/domain"
)

type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}

type SubscriptionResponse struct {
	ID          int64  `json:"id"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (r CreateSubscriptionRequest) ToDomain() (domain.Subscription, error) {
	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		return domain.Subscription{}, err
	}

	startMonth, err := domain.ParseMonth(r.StartDate)
	if err != nil {
		return domain.Subscription{}, err
	}

	var endMonth *time.Time

	if r.EndDate != "" {
		parsedEndMonth, err := domain.ParseMonth(r.EndDate)
		if err != nil {
			return domain.Subscription{}, err
		}

		endMonth = &parsedEndMonth
	}

	return domain.Subscription{
		ServiceName: r.ServiceName,
		Price:       r.Price,
		UserID:      userID,
		StartMonth:  startMonth,
		EndMonth:    endMonth,
	}, nil
}

func NewSubscriptionResponse(sub domain.Subscription) SubscriptionResponse {
	response := SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID.String(),
		StartDate:   domain.FormatMonth(sub.StartMonth),
		CreatedAt:   sub.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   sub.UpdatedAt.Format(time.RFC3339),
	}

	if sub.EndMonth != nil {
		response.EndDate = domain.FormatMonth(*sub.EndMonth)
	}

	return response
}
