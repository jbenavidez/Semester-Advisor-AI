package ports

import (
	"context"
	"io"
	"semester-advisor-ai/internal/domain"
)

type DatasetProcessor interface {
	Process(ctx context.Context, reader io.Reader, uploadedFile *domain.UploadedFile) error
}
