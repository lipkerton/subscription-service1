package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(db *pgxpool.Pool, subscriptionService SubscriptionService) http.Handler {
	r := chi.NewRouter()
	h := NewHandler(db)
	subscriptionHandler := NewSubscriptionHandler(subscriptionService)

	r.Get("/health", h.HealthCheck)

	r.Post("/subscriptions", subscriptionHandler.Create)
	r.Get("/subscriptions", subscriptionHandler.List)
	r.Get("/subscriptions/summary", subscriptionHandler.Summary)
	r.Get("/subscriptions/{id}", subscriptionHandler.GetByID)
	r.Put("/subscriptions/{id}", subscriptionHandler.Update)
	r.Delete("/subscriptions/{id}", subscriptionHandler.Delete)

	return r
}
