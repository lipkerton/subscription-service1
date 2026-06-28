package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lipkerton/subscription-service1/internal/config"
	httptransport "github.com/lipkerton/subscription-service1/internal/transport/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	r := httptransport.NewRouter()
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	log.Printf("server started at %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
