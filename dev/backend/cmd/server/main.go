package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"ai-knowledge-hub/internal/api"
	"ai-knowledge-hub/internal/storage/postgres"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	connectCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	store, err := postgres.Open(connectCtx, databaseURL)
	cancel()
	if err != nil {
		log.Fatalf("database startup failed: %v", err)
	}
	defer store.Close()

	server := &http.Server{
		Addr:              ":8081",
		Handler:           api.NewServer(store).Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("server listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server stopped: %v", err)
	}
}
