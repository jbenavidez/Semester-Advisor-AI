package dto

import "mime/multipart"

type UploadDatasetInput struct {
	DatasetName  string
	SourcePeriod string
	Description  string
	File         multipart.File
	Header       *multipart.FileHeader
}
