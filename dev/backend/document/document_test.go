package document

import (
	"errors"
	"testing"
)

func TestNewDocumentValidation(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		title   string
		content string
		wantErr error
	}{
		{name: "missing id", title: "Title", content: "Content", wantErr: ErrEmptyID},
		{name: "missing title", id: "1", title: "  ", content: "Content", wantErr: ErrEmptyTitle},
		{name: "missing content", id: "1", title: "Title", content: "\n\t", wantErr: ErrEmptyContent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDocument(tt.id, tt.title, tt.content)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewDocument() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewDocument(t *testing.T) {
	doc, err := NewDocument("1", "Go", "A programming language")
	if err != nil {
		t.Fatalf("NewDocument() unexpected error: %v", err)
	}
	if doc.ID != "1" || doc.Title != "Go" || doc.Content != "A programming language" {
		t.Fatalf("NewDocument() = %#v", doc)
	}
}
