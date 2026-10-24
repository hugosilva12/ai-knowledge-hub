package document

import (
	"errors"
	"strings"
)

var (
	ErrEmptyID      = errors.New("id cannot be empty")
	ErrEmptyTitle   = errors.New("title cannot be empty")
	ErrEmptyContent = errors.New("content cannot be empty")
)

type Document struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func NewDocument(id, title, content string) (Document, error) {
	if strings.TrimSpace(id) == "" {
		return Document{}, ErrEmptyID
	}
	if strings.TrimSpace(title) == "" {
		return Document{}, ErrEmptyTitle
	}
	if strings.TrimSpace(content) == "" {
		return Document{}, ErrEmptyContent
	}
	return Document{ID: id, Title: title, Content: content}, nil
}

func (d Document) Summary() string {
	if len(d.Content) <= 50 {
		return d.Content
	}
	return d.Content[:50] + "..."
}

func (d *Document) SetTitle(title string) {
	d.Title = title
}
