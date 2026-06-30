package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lipkerton/subscription-service1/internal/domain"
	"github.com/lipkerton/subscription-service1/internal/transport/dto"
)

type SubscriptionService interface {
	Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error)
	GetByID(ctx context.Context, id int64) (domain.Subscription, error)
	Update(ctx context.Context, sub domain.Subscription) (domain.Subscription, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filter domain.SubscriptionFilter) ([]domain.Subscription, error)
	CalculateSummary(ctx context.Context, filter domain.SummaryFilter) (int, error)
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

func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid subscription id")
		return
	}

	err = h.service.Delete(r.Context(), id)
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

	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	filter := domain.SubscriptionFilter{}

	userIDParam := r.URL.Query().Get("user_id")
	if userIDParam != "" {
		userID, err := uuid.Parse(userIDParam)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid user_id")
			return
		}

		filter.UserID = &userID
	}

	serviceNameParam := r.URL.Query().Get("service_name")
	if serviceNameParam != "" {
		filter.ServiceName = &serviceNameParam
	}

	subscriptions, err := h.service.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(dto.NewListSubscriptionsResponse(subscriptions))
}

func (h *SubscriptionHandler) Summary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fromParam := r.URL.Query().Get("from")
	if fromParam == "" {
		writeError(w, http.StatusBadRequest, "from is required")
		return
	}

	toParam := r.URL.Query().Get("to")
	if toParam == "" {
		writeError(w, http.StatusBadRequest, "to is required")
		return
	}

	fromMonth, err := domain.ParseMonth(fromParam)
	if err != nil {
		writeError(w, http.StatusBadRequest, "from must be in MM-YYYY format")
		return
	}

	toMonth, err := domain.ParseMonth(toParam)
	if err != nil {
		writeError(w, http.StatusBadRequest, "to must be in MM-YYYY format")
		return
	}

	filter := domain.SummaryFilter{
		FromMonth: fromMonth,
		ToMonth:   toMonth,
	}

	userIDParam := r.URL.Query().Get("user_id")
	if userIDParam != "" {
		userID, err := uuid.Parse(userIDParam)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid user_id")
			return
		}

		filter.UserID = &userID
	}

	serviceNameParam := r.URL.Query().Get("service_name")
	if serviceNameParam != "" {
		filter.ServiceName = &serviceNameParam
	}

	total, err := h.service.CalculateSummary(r.Context(), filter)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidPeriod):
			writeError(w, http.StatusBadRequest, "to must be greater than or equal to from")
			return
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(dto.SummaryResponse{
		Total: total,
	})
}
