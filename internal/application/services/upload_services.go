package services

import (
	"semester-advisor-ai/internal/domain"
	"semester-advisor-ai/internal/ports"
	"sync"
)

type UploadJob struct {
	UploadedFile *domain.UploadedFile
	SessionID    string
}

type UploadStatusEvent struct {
	SessionID    string
	UploadedFile *domain.UploadedFile
	Status       string
	Message      string
	ErrorMessage string
}

type UploadService struct {
	repo             ports.UploadedFileRepository
	uploadJobChan    chan *UploadJob
	UploadStatusChan chan *UploadStatusEvent
	wg               sync.WaitGroup
}

func NewUploadServices(r ports.UploadedFileRepository, chunkSize int) *UploadService {

	service := &UploadService{
		repo: r,
	}
	return service
}
