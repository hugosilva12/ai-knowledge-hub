package document

import "errors"

type Document struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func NewDocument(id, title, content string) (Document, error) {
	if content == "" {
		return Document{}, errors.New("Content cannot be empty")
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
