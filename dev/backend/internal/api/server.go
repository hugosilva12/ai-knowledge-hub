package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"ai-knowledge-hub/document"
	"ai-knowledge-hub/internal/llm"
	"ai-knowledge-hub/internal/storage"
)

const (
	maxRequestBodySize = 1 << 20
	maxQuestionLength  = 8_000
)

type Server struct {
	store     storage.Store
	generator llm.Generator
}

func NewServer(store storage.Store, generator llm.Generator) *Server {
	return &Server{store: store, generator: generator}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/documents", s.documentsHandler)
	mux.HandleFunc("/chat", s.chatHandler)
	return mux
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := s.store.Ping(ctx); err != nil {
		slog.Error("health check failed", "error", err)
		writeError(w, http.StatusServiceUnavailable, "service unavailable")
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "ok\n")
}

func (s *Server) documentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listDocuments(w, r)
	case http.MethodPost:
		s.createDocument(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) listDocuments(w http.ResponseWriter, r *http.Request) {
	documents, err := s.store.ListDocuments(r.Context())
	if err != nil {
		slog.Error("list documents", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, documents)
}

func (s *Server) createDocument(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if !decodeJSONBody(w, r, &input) {
		return
	}

	doc, err := s.store.CreateDocument(r.Context(), input.Title, input.Content)
	if err != nil {
		if errors.Is(err, document.ErrEmptyTitle) || errors.Is(err, document.ErrEmptyContent) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		slog.Error("create document", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, doc)
}

func (s *Server) chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var input struct {
		Question string `json:"question"`
	}
	if !decodeJSONBody(w, r, &input) {
		return
	}
	if strings.TrimSpace(input.Question) == "" {
		writeError(w, http.StatusBadRequest, "question cannot be empty")
		return
	}
	if len(input.Question) > maxQuestionLength {
		writeError(w, http.StatusBadRequest, "question is too long")
		return
	}

	answer, err := s.generator.Generate(r.Context(), input.Question)
	if err != nil {
		slog.Error("generate chat answer", "error", err)
		writeError(w, http.StatusBadGateway, "language model unavailable")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"answer": answer})
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, destination any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return false
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "body must contain one JSON object")
		return false
	}
	return true
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
