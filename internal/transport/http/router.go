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

	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", subscriptionHandler.Create)
	})

	return r
}
