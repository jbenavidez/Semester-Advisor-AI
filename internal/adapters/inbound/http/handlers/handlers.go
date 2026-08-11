package handlers

import "semester-advisor-ai/internal/application/services"

type Handlers struct {
	uploadService *services.UploadService
	plannerSevice *services.SemesterPlannerServices
}

func New(uploadService *services.UploadService, plannerSevice *services.SemesterPlannerServices) *Handlers {
	h := &Handlers{
		uploadService: uploadService,
		plannerSevice: plannerSevice,
	}

	return h
}
