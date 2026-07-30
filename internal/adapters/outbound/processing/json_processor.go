package processing

import (
	"context"
	"fmt"
	"io"
	"semester-advisor-ai/internal/domain"
	"semester-advisor-ai/internal/ports"
)

type JSONProcessor struct {
	repo ports.UploadedFileRepository
}

func NewJSONProcessor(repo ports.UploadedFileRepository) ports.DatasetProcessor {
	return &JSONProcessor{
		repo: repo,
	}
}

func (j *JSONProcessor) Process(ctx context.Context, reader io.Reader, uploadedFile *domain.UploadedFile) error {
	fmt.Println("About to proces the file")
	return nil
}
