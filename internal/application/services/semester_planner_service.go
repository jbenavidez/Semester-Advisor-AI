package services

import (
	"semester-advisor-ai/internal/ports"
)

type SemesterPlannerServices struct {
	repo ports.ProfessorReviewRepository
}

func NewSemesterPlannerServices(repo ports.ProfessorReviewRepository) *SemesterPlannerServices {

	return &SemesterPlannerServices{
		repo: repo,
	}
}
