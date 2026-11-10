package storage

import (
	"context"
	"strconv"
	"sync"

	"ai-knowledge-hub/document"
)

type MemoryStore struct {
	mu        sync.RWMutex
	documents []document.Document
	nextID    uint64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		documents: make([]document.Document, 0),
		nextID:    1,
	}
}

func (s *MemoryStore) CreateDocument(ctx context.Context, title, content string) (document.Document, error) {
	if err := ctx.Err(); err != nil {
		return document.Document{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	doc, err := document.NewDocument(strconv.FormatUint(s.nextID, 10), title, content)
	if err != nil {
		return document.Document{}, err
	}
	s.nextID++
	s.documents = append(s.documents, doc)
	return doc, nil
}

func (s *MemoryStore) ListDocuments(ctx context.Context) ([]document.Document, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	documents := make([]document.Document, len(s.documents))
	copy(documents, s.documents)
	return documents, nil
}

func (s *MemoryStore) Ping(ctx context.Context) error {
	return ctx.Err()
}
