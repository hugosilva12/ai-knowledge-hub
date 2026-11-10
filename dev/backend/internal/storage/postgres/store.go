package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"ai-knowledge-hub/document"
)

type Store struct {
	pool *pgxpool.Pool
}

func Open(ctx context.Context, databaseURL string) (*Store, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database configuration: %w", err)
	}
	config.ConnConfig.RuntimeParams["application_name"] = "ai-knowledge-hub"

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	return &Store{pool: pool}, nil
}

func (s *Store) Close() {
	s.pool.Close()
}

func (s *Store) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

func (s *Store) CreateDocument(ctx context.Context, title, content string) (document.Document, error) {
	if err := document.Validate(title, content); err != nil {
		return document.Document{}, err
	}

	const query = `
		INSERT INTO documents (title, content)
		VALUES ($1, $2)
		RETURNING id::text, title, content`

	var doc document.Document
	if err := s.pool.QueryRow(ctx, query, title, content).Scan(&doc.ID, &doc.Title, &doc.Content); err != nil {
		return document.Document{}, fmt.Errorf("create document: %w", err)
	}
	return doc, nil
}

func (s *Store) ListDocuments(ctx context.Context) ([]document.Document, error) {
	const query = `
		SELECT id::text, title, content
		FROM documents
		ORDER BY id`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	defer rows.Close()

	documents := make([]document.Document, 0)
	for rows.Next() {
		var doc document.Document
		if err := rows.Scan(&doc.ID, &doc.Title, &doc.Content); err != nil {
			return nil, fmt.Errorf("scan document: %w", err)
		}
		documents = append(documents, doc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate documents: %w", err)
	}
	return documents, nil
}
