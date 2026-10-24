package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"ai-knowledge-hub/internal/api"
)

func main() {
	server := &http.Server{
		Addr:              ":8081",
		Handler:           api.NewServer().Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("server listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server stopped: %v", err)
	}
}
