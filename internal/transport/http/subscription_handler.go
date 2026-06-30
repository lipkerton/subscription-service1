package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lipkerton/subscription-service1/internal/domain"
	"github.com/lipkerton/subscription-service1/internal/transport/dto"
)

type SubscriptionService interface {
	Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error)
	GetByID(ctx context.Context, id int64) (domain.Subscription, error)
	Update(ctx context.Context, sub domain.Subscription) (domain.Subscription, error)
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

func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid subscription id")
		return
	}

	sub, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSubscriptionNotFound):
			writeError(w, http.StatusNotFound, "subscription not found")
			return
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(dto.NewSubscriptionResponse(sub))
}

func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid subscription id")
		return
	}

	var request dto.UpdateSubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	sub, err := request.ToDomain(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	updatedSub, err := h.service.Update(r.Context(), sub)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrServiceNameRequired),
			errors.Is(err, domain.ErrPriceMustBePositive),
			errors.Is(err, domain.ErrInvalidPeriod):
			writeError(w, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, domain.ErrSubscriptionNotFound):
			writeError(w, http.StatusNotFound, "subscription not found")
			return
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(dto.NewSubscriptionResponse(updatedSub))
}
