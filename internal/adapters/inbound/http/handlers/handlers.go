package handlers

import "semester-advisor-ai/internal/application/services"

type Handlers struct {
	uploadService   *services.UploadService
	plannerSevice   *services.SemesterPlannerServices
	semesterAdvisor *services.SemesterAdvisorService
}

func New(uploadService *services.UploadService, plannerSevice *services.SemesterPlannerServices, semesterAdvisor *services.SemesterAdvisorService) *Handlers {
	h := &Handlers{
		uploadService:   uploadService,
		plannerSevice:   plannerSevice,
		semesterAdvisor: semesterAdvisor,
	}

	return h
}
