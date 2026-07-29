package handlers

import "semester-advisor-ai/internal/application/services"

type Handlers struct {
	uploadService *services.UploadService
}

func New(uploadService *services.UploadService) *Handlers {
	h := &Handlers{
		uploadService: uploadService,
	}

	return h
}
