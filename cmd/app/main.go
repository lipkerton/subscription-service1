package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lipkerton/subscription-service1/internal/config"
	"github.com/lipkerton/subscription-service1/internal/service"
	storagepostgres "github.com/lipkerton/subscription-service1/internal/storage/postgres"
	httptransport "github.com/lipkerton/subscription-service1/internal/transport/http"
)

func main() {
	log := config.NewLogger()
	cfg, err := config.Load()
	if err != nil {
		log.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	dbPool, err := storagepostgres.NewPool(ctx, cfg.Postgres.DSN())
	if err != nil {
		log.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	log.Info("connected to postgres")

	subscriptionRepo := storagepostgres.NewSubscriptionRepository(dbPool)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)

	r := httptransport.NewRouter(dbPool, subscriptionService)
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		log.Info("server started", "addr", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}
	log.Info("server stopped gracefully")
}
