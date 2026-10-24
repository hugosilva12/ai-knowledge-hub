package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"ai-knowledge-hub/internal/llm"
)

const maxResponseSize = 4 << 20

type Client struct {
	endpoint   string
	model      string
	httpClient *http.Client
}

func New(baseURL, model string, httpClient *http.Client) (*Client, error) {
	parsedURL, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") || parsedURL.Host == "" {
		return nil, errors.New("invalid Ollama URL")
	}
	if strings.TrimSpace(model) == "" {
		return nil, errors.New("Ollama model cannot be empty")
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 90 * time.Second}
	}

	return &Client{
		endpoint:   strings.TrimRight(parsedURL.String(), "/") + "/api/generate",
		model:      strings.TrimSpace(model),
		httpClient: httpClient,
	}, nil
}

func (c *Client) Generate(ctx context.Context, prompt string) (string, error) {
	if strings.TrimSpace(prompt) == "" {
		return "", llm.ErrEmptyPrompt
	}

	payload, err := json.Marshal(struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
		Stream bool   `json:"stream"`
	}{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
	})
	if err != nil {
		return "", fmt.Errorf("encode Ollama request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("create Ollama request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("call Ollama: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseSize+1))
	if err != nil {
		return "", fmt.Errorf("read Ollama response: %w", err)
	}
	if len(body) > maxResponseSize {
		return "", errors.New("Ollama response exceeds size limit")
	}

	var result struct {
		Response string `json:"response"`
		Error    string `json:"error"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("decode Ollama response: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		if result.Error == "" {
			result.Error = http.StatusText(response.StatusCode)
		}
		return "", fmt.Errorf("Ollama returned status %d: %s", response.StatusCode, result.Error)
	}
	if strings.TrimSpace(result.Response) == "" {
		return "", errors.New("Ollama returned an empty response")
	}

	return result.Response, nil
}
