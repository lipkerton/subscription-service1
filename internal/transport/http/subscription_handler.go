package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lipkerton/subscription-service1/internal/domain"
	"github.com/lipkerton/subscription-service1/internal/transport/dto"
)

type SubscriptionService interface {
	Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error)
}

type SubscriptionHandler struct {
	service SubscriptionService
}

func NewSubscriptionHandler(service SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
	}
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request dto.CreateSubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	sub, err := request.ToDomain()
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	createdSub, err := h.service.Create(r.Context(), sub)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrServiceNameRequired),
			errors.Is(err, domain.ErrPriceMustBePositive),
			errors.Is(err, domain.ErrInvalidPeriod):
			writeError(w, http.StatusBadRequest, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(dto.NewSubscriptionResponse(createdSub))
}
