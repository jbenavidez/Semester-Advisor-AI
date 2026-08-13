package ports

import "context"

type LanguageModel interface {
	Generate(ctx context.Context, prompt string) (string, error)
}
