package llm

import (
	"context"
	"errors"
)

var ErrEmptyPrompt = errors.New("prompt cannot be empty")

type Generator interface {
	Generate(ctx context.Context, prompt string) (string, error)
}
