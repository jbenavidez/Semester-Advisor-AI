package services

import (
	"context"
	"encoding/json"
	"semester-advisor-ai/internal/application/prompts"
	"semester-advisor-ai/internal/ports"
)

type SemesterPlannerServices struct {
	repo               ports.ProfessorReviewRepository
	languageModel      ports.LanguageModel
	semesterPlanMemory ports.SemesterPlanMemory
}

func NewSemesterPlannerServices(repo ports.ProfessorReviewRepository, languageModel ports.LanguageModel, memoryplan ports.SemesterPlanMemory) *SemesterPlannerServices {

	return &SemesterPlannerServices{
		repo:               repo,
		languageModel:      languageModel,
		semesterPlanMemory: memoryplan,
	}
}

func (s *SemesterPlannerServices) PlanSemester(ctx context.Context, planID string) (string, error) {
	courses, err := s.semesterPlanMemory.GetPlan(ctx, planID)
	if err != nil {
		return "", err
	}

	reviews, err := s.repo.GetReviewsForCourses(courses)
	if err != nil {
		return "", err
	}

	courseData, err := json.Marshal(courses)
	if err != nil {
		return "", err
	}

	reviewData, err := json.Marshal(reviews)
	if err != nil {
		return "", err
	}

	prompt, err := prompts.BuildSemesterPlannerPrompt(courseData, reviewData)
	if err != nil {
		return "", err
	}

	answer, err := s.languageModel.Generate(ctx, prompt)
	if err != nil {
		return "", err
	}

	return answer, nil
}
