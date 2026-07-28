package models

import "time"

type UploadedFile struct {
	ID               string    `json:"id"`
	OriginalFileName string    `json:"original_file_name"`
	StoredFileName   string    `json:"stored_file_name"`
	FilePath         string    `json:"file_path"`
	Description      string    `json:"description"`
	ContentType      string    `json:"content_type"`
	Size             int64     `json:"size"`
	DatasetName      string    `json:"dataset_name"`
	SourcePeriod     string    `json:"source_period"`
	Status           string    `json:"status"`
	ErrorMessage     string    `json:"error_message,omitempty"`
	TotalReviews     int       `json:"total_reviews"`
	ProcessedReviews int       `json:"processed_reviews"`
	FailedReviews    int       `json:"failed_reviews"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
