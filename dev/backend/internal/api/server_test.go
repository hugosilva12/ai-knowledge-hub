package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"ai-knowledge-hub/document"
	"ai-knowledge-hub/internal/storage"
)

func newTestHandler() http.Handler {
	return NewServer(storage.NewMemoryStore()).Handler()
}

func TestHealth(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	newTestHandler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if response.Body.String() != "ok\n" {
		t.Fatalf("body = %q, want %q", response.Body.String(), "ok\n")
	}
}

func TestDocumentsLifecycle(t *testing.T) {
	handler := newTestHandler()
	empty := httptest.NewRecorder()
	handler.ServeHTTP(empty, httptest.NewRequest(http.MethodGet, "/documents", nil))
	if empty.Code != http.StatusOK || strings.TrimSpace(empty.Body.String()) != "[]" {
		t.Fatalf("empty list response = %d %q, want 200 []", empty.Code, empty.Body.String())
	}

	invalid := httptest.NewRecorder()
	handler.ServeHTTP(invalid, httptest.NewRequest(http.MethodPost, "/documents", strings.NewReader(`{"title":"","content":"content"}`)))
	if invalid.Code != http.StatusBadRequest {
		t.Fatalf("invalid document status = %d, want %d", invalid.Code, http.StatusBadRequest)
	}

	created := httptest.NewRecorder()
	handler.ServeHTTP(created, httptest.NewRequest(http.MethodPost, "/documents", strings.NewReader(`{"title":"Go","content":"A programming language"}`)))
	if created.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d; body = %s", created.Code, http.StatusCreated, created.Body.String())
	}

	var createdDocument document.Document
	if err := json.NewDecoder(created.Body).Decode(&createdDocument); err != nil {
		t.Fatalf("decode created document: %v", err)
	}
	if createdDocument.ID != "1" {
		t.Fatalf("created id = %q, want %q", createdDocument.ID, "1")
	}

	listed := httptest.NewRecorder()
	handler.ServeHTTP(listed, httptest.NewRequest(http.MethodGet, "/documents", nil))
	if listed.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", listed.Code, http.StatusOK)
	}

	var documents []document.Document
	if err := json.NewDecoder(listed.Body).Decode(&documents); err != nil {
		t.Fatalf("decode document list: %v", err)
	}
	if len(documents) != 1 || documents[0].ID != "1" {
		t.Fatalf("documents = %#v, want one document with id 1", documents)
	}
}

func TestConcurrentDocumentCreationUsesUniqueIDs(t *testing.T) {
	const total = 50
	handler := newTestHandler()

	var wg sync.WaitGroup
	statuses := make(chan int, total)
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			body := strings.NewReader(fmt.Sprintf(`{"title":"Document %d","content":"Content"}`, i))
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/documents", body))
			statuses <- response.Code
		}(i)
	}
	wg.Wait()
	close(statuses)

	for status := range statuses {
		if status != http.StatusCreated {
			t.Fatalf("create status = %d, want %d", status, http.StatusCreated)
		}
	}

	listed := httptest.NewRecorder()
	handler.ServeHTTP(listed, httptest.NewRequest(http.MethodGet, "/documents", nil))

	var documents []document.Document
	if err := json.NewDecoder(listed.Body).Decode(&documents); err != nil {
		t.Fatalf("decode document list: %v", err)
	}
	if len(documents) != total {
		t.Fatalf("document count = %d, want %d", len(documents), total)
	}

	ids := make(map[string]struct{}, total)
	for _, doc := range documents {
		if _, exists := ids[doc.ID]; exists {
			t.Fatalf("duplicate document id %q", doc.ID)
		}
		ids[doc.ID] = struct{}{}
	}
}
