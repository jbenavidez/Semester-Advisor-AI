package ports

import (
	"context"
	"semester-advisor-ai/internal/domain"
)

type UploadedFileRepository interface {
	SaveFile()
	SaveFileMetaData(storedFile *domain.UploadedFile) (*domain.UploadedFile, error)
	GetFileByName(string) (*domain.UploadedFile, bool, error)
	GetAllUploadFiles() ([]domain.UploadedFile, error)
	UpdateFile(uploadedFile *domain.UploadedFile) error
	SaveReview(ctx context.Context, review *domain.ProfessorReview) error
}
