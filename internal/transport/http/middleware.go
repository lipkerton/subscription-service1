package http

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func RequestLogger(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrappedWriter := newResponseWriter(w)

			next.ServeHTTP(wrappedWriter, r)

			duration := time.Since(start)

			log.Info(
				"http request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.RawQuery,
				"status", wrappedWriter.statusCode,
				"duration", duration.String(),
				"remote_addr", r.RemoteAddr,
			)
		})
	}
}

func Recoverer(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					log.Error(
						"panic recovered",
						"method", r.Method,
						"path", r.URL.Path,
						"query", r.URL.RawQuery,
						"panic", recovered,
					)

					writeError(w, http.StatusInternalServerError, "internal server error")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
