package storage

import (
	"context"

	"ai-knowledge-hub/document"
)

type Store interface {
	CreateDocument(ctx context.Context, title, content string) (document.Document, error)
	ListDocuments(ctx context.Context) ([]document.Document, error)
	Ping(ctx context.Context) error
}
