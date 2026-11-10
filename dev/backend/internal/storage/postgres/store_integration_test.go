//go:build integration

package postgres

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestStorePersistsDocumentsAcrossConnections(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	first, err := Open(ctx, databaseURL)
	if err != nil {
		t.Fatalf("open first connection: %v", err)
	}
	doc, err := first.CreateDocument(ctx, "Persistence test", "Stored in PostgreSQL")
	if err != nil {
		first.Close()
		t.Fatalf("create document: %v", err)
	}
	first.Close()

	second, err := Open(ctx, databaseURL)
	if err != nil {
		t.Fatalf("open second connection: %v", err)
	}
	defer second.Close()
	defer func() {
		_, _ = second.pool.Exec(context.Background(), "DELETE FROM documents WHERE id = $1", doc.ID)
	}()

	documents, err := second.ListDocuments(ctx)
	if err != nil {
		t.Fatalf("list documents: %v", err)
	}
	for _, stored := range documents {
		if stored.ID == doc.ID && stored.Title == doc.Title && stored.Content == doc.Content {
			return
		}
	}
	t.Fatalf("document %s was not persisted", doc.ID)
}
