package dto

type UploadFileResponse struct {
	UploadedFiles []UploadedFileDTO
}

type UploadedFileDTO struct {
	ID               string
	OriginalFileName string
	StoredFileName   string
	FilePath         string
	Description      string
	ContentType      string
	Size             int64
	DatasetName      string
	SourcePeriod     string
	Status           string
	ErrorMessage     string
	TotalReviews     int
	ProcessedReviews int
	FailedReviews    int
	CreatedAt        string
	UpdatedAt        string
}
