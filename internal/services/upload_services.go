package services

import (
	"semester-advisor-ai/internal/models"
	"semester-advisor-ai/internal/repository"
	"sync"
)

type UploadJob struct {
	UploadedFile *models.UploadedFile
	SessionID    string
}

type UploadStatusEvent struct {
	SessionID    string
	UploadedFile *models.UploadedFile
	Status       string
	Message      string
	ErrorMessage string
}

type UploadService struct {
	repo             repository.DatabaseRepo
	uploadJobChan    chan *UploadJob
	UploadStatusChan chan *UploadStatusEvent
	wg               sync.WaitGroup
}

func NewUploadServices(r repository.DatabaseRepo, chunkSize int) *UploadService {

	service := &UploadService{
		repo: r,
	}
	return service
}
