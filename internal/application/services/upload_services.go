package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"semester-advisor-ai/internal/adapters/inbound/http/dto"
	"semester-advisor-ai/internal/domain"
	"semester-advisor-ai/internal/ports"
	"strings"
	"sync"
	"time"
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
	storage          ports.FileStorage
	uploadJobChan    chan *UploadJob
	UploadStatusChan chan *UploadStatusEvent
	wg               sync.WaitGroup
}

func NewUploadServices(r ports.UploadedFileRepository, s ports.FileStorage, chunkSize int) *UploadService {

	service := &UploadService{
		repo:             r,
		storage:          s,
		uploadJobChan:    make(chan *UploadJob),
		UploadStatusChan: make(chan *UploadStatusEvent),
	}
	service.StartWorker()
	return service
}

func (s *UploadService) StartWorker() {
	s.wg.Add(1)
	//spin go rutine
	go s.processFileWorker()
}

func (s *UploadService) processFileWorker() {
	defer s.wg.Done()

	// Wait for upload jobs to be added to the channel.
	for uploadJob := range s.uploadJobChan {
		if uploadJob == nil || uploadJob.UploadedFile == nil {
			continue
		}
		if err := s.ProcessFile(context.Background(), uploadJob.UploadedFile); err != nil {
			fmt.Println("failed to process file:", err)
			// Send failed event to notification chan.
			s.UploadStatusChan <- &UploadStatusEvent{
				SessionID:    uploadJob.SessionID,
				UploadedFile: uploadJob.UploadedFile,
				Status:       uploadJob.UploadedFile.Status,
				Message:      "failed to process file",
				ErrorMessage: err.Error(),
			}

			continue
		}
		// Send success event to notification chan.
		s.UploadStatusChan <- &UploadStatusEvent{
			SessionID:    uploadJob.SessionID,
			UploadedFile: uploadJob.UploadedFile,
			Status:       uploadJob.UploadedFile.Status,
			Message:      "file processed successfully",
			ErrorMessage: "",
		}
	}
}

func (s *UploadService) ProcessFile(ctx context.Context, uploadedFile *domain.UploadedFile) error {
	file, err := os.Open(uploadedFile.FilePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	ext := strings.ToLower(filepath.Ext(uploadedFile.FilePath))

	switch ext {
	case ".csv":
		// Mark file as processing
		uploadedFile.Status = "processing"
		uploadedFile.ErrorMessage = ""
		// to be continue
		fmt.Println("gondor", uploadedFile)

	default:
		return fmt.Errorf("unsupported file extension: %s", ext)
	}
	return nil
}

func (s *UploadService) SaveUploadedFile(ctx context.Context, uploadDataset dto.UploadDatasetInput) (*ports.StoredFile, error) {
	// Check if file already exists before saving it locally.
	_, found, err := s.repo.GetFileByName(uploadDataset.Header.Filename)
	if err != nil {
		return nil, err
	}
	if found {
		return nil, fmt.Errorf("file already uploaded")
	}

	//save file
	storedFile, err := s.storage.Save(ctx, uploadDataset.File, uploadDataset.Header)
	if err != nil {
		return nil, err
	}
	// safe file metatada in DB.
	now := time.Now()
	uploadedFile := &domain.UploadedFile{
		OriginalFileName: storedFile.OriginalFileName,
		StoredFileName:   storedFile.StoredFileName,
		FilePath:         storedFile.FilePath,
		Description:      uploadDataset.Description,
		ContentType:      storedFile.ContentType,
		Size:             storedFile.Size,
		DatasetName:      uploadDataset.DatasetName,
		SourcePeriod:     uploadDataset.SourcePeriod,
		Status:           "pending",
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	savedFile, err := s.repo.SaveFileMetaData(uploadedFile)
	if err != nil {
		return nil, err
	}
	fmt.Println("valinor", savedFile)

	return nil, nil

}
