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
	UploadedFile *domain.UploadedFile
	Status       string
	Message      string
	ErrorMessage string
}

type UploadService struct {
	repo             ports.UploadedFileRepository
	storage          ports.FileStorage
	jsonProcessor    ports.DatasetProcessor
	uploadJobChan    chan *UploadJob
	UploadStatusChan chan *UploadStatusEvent
	wg               sync.WaitGroup
}

func NewUploadServices(r ports.UploadedFileRepository, s ports.FileStorage, jsonProcessor ports.DatasetProcessor, chunkSize int) *UploadService {

	service := &UploadService{
		repo:             r,
		storage:          s,
		jsonProcessor:    jsonProcessor,
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
				UploadedFile: uploadJob.UploadedFile,
				Status:       uploadJob.UploadedFile.Status,
				Message:      "failed to process file",
				ErrorMessage: err.Error(),
			}

			continue
		}
		fmt.Println("valinor_03", uploadJob.UploadedFile.Status)
		// Send success event to notification chan.
		s.UploadStatusChan <- &UploadStatusEvent{
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
		return fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	ext := strings.ToLower(filepath.Ext(uploadedFile.FilePath))

	switch ext {
	case ".json":
		uploadedFile.Status = "processing"
		uploadedFile.ErrorMessage = ""
		uploadedFile.UpdatedAt = time.Now().UTC()

		if err := s.repo.UpdateFile(uploadedFile); err != nil {
			return fmt.Errorf("failed to update uploaded file status to processing: %w", err)
		}

		if err := s.jsonProcessor.Process(ctx, file, uploadedFile); err != nil {
			uploadedFile.Status = "failed"
			uploadedFile.ErrorMessage = err.Error()
			uploadedFile.UpdatedAt = time.Now().UTC()

			if updateErr := s.repo.UpdateFile(uploadedFile); updateErr != nil {
				return fmt.Errorf("failed to process JSON file: %w; failed to update file status: %v", err, updateErr)
			}

			return fmt.Errorf("failed to process JSON file: %w", err)
		}

		uploadedFile.Status = "processed"
		uploadedFile.ErrorMessage = ""
		uploadedFile.UpdatedAt = time.Now().UTC()
		fmt.Println("gondor_update_status", uploadedFile.Status)
		if err := s.repo.UpdateFile(uploadedFile); err != nil {
			return fmt.Errorf("failed to update uploaded file status to processed: %w", err)
		}

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
	//set uploadJob
	uploadJob := &UploadJob{
		UploadedFile: uploadedFile,
	}
	//send upload-job to chan
	s.uploadJobChan <- uploadJob

	return storedFile, nil

}

func (s *UploadService) GetAllUploadedFile() (*dto.UploadFileResponse, error) {
	// Get all dfiles
	files, err := s.repo.GetAllUploadFiles()
	if err != nil {
		return nil, err
	}

	uploadedFiles := make([]dto.UploadedFileDTO, len(files))

	for i := range files {
		file := files[i]
		uploadedFiles[i] = dto.UploadedFileDTO{
			ID:               file.ID,
			OriginalFileName: file.OriginalFileName,
			StoredFileName:   file.StoredFileName,
			FilePath:         file.FilePath,
			Description:      file.Description,
			ContentType:      file.ContentType,
			Size:             file.Size,
			DatasetName:      file.DatasetName,
			SourcePeriod:     file.SourcePeriod,
			Status:           file.Status,
			ErrorMessage:     file.ErrorMessage,
			TotalReviews:     file.TotalReviews,
			ProcessedReviews: file.ProcessedReviews,
			FailedReviews:    file.FailedReviews,
			CreatedAt:        file.CreatedAt.Format("2006-01-02 15:04"),
			UpdatedAt:        file.UpdatedAt.Format("2006-01-02 15:04"),
		}
	}
	//set response
	res := &dto.UploadFileResponse{
		UploadedFiles: uploadedFiles,
	}
	return res, nil
}
