package main

import (
	"log"
	"net/http"

	httptransport "github.com/lipkerton/subscription-service1/internal/transport/http"
)

func main() {
	r := httptransport.NewRouter()
	log.Println("server started at :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
