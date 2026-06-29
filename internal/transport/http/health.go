package http

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	db *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{
		db: db,
	}
}

type healthResponse struct {
	Status   string `json:"status"`
	Postgres string `json:"postgres"`
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := h.db.Ping(r.Context()); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		if err := json.NewEncoder(w).Encode(healthResponse{
			Status:   "error",
			Postgres: "unavailable",
		}); err != nil {
			http.Error(w, "failed to encode error response", http.StatusInternalServerError)
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(healthResponse{
		Status:   "ok",
		Postgres: "ok",
	}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
