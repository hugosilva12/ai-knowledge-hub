package ollama

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestGenerate(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/generate" {
			t.Fatalf("request = %s %s, want POST /api/generate", r.Method, r.URL.Path)
		}

		var request struct {
			Model  string `json:"model"`
			Prompt string `json:"prompt"`
			Stream bool   `json:"stream"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if request.Model != "test-model" || request.Prompt != "What is Go?" || request.Stream {
			t.Fatalf("unexpected request: %#v", request)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"response":"Go is a programming language.","done":true}`)),
		}, nil
	})}

	client, err := New("http://ollama.test", "test-model", httpClient)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	answer, err := client.Generate(context.Background(), "What is Go?")
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	if answer != "Go is a programming language." {
		t.Fatalf("answer = %q", answer)
	}
}

func TestGenerateReportsOllamaError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"error":"model not found"}`)),
		}, nil
	})}

	client, err := New("http://ollama.test", "missing-model", httpClient)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	_, err = client.Generate(context.Background(), "question")
	if err == nil || !strings.Contains(err.Error(), "model not found") {
		t.Fatalf("Generate() error = %v, want model not found", err)
	}
}

func TestNewValidatesConfiguration(t *testing.T) {
	if _, err := New("not-a-url", "model", nil); err == nil {
		t.Fatal("New() accepted an invalid URL")
	}
	if _, err := New("http://localhost:11434", "", nil); err == nil {
		t.Fatal("New() accepted an empty model")
	}
}
