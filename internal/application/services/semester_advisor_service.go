package services

import (
	"context"
	"encoding/json"
	"fmt"

	"semester-advisor-ai/internal/application/prompts"
	"semester-advisor-ai/internal/domain"
	"semester-advisor-ai/internal/ports"

	"github.com/google/uuid"
)

type SemesterAdvisorService struct {
	languageModel      ports.LanguageModel
	repo               ports.ProfessorReviewRepository
	semesterPlanMemory ports.SemesterPlanMemory
}

func NewSemesterAdvisorService(languageModel ports.LanguageModel, repo ports.ProfessorReviewRepository, memoryplan ports.SemesterPlanMemory) *SemesterAdvisorService {
	return &SemesterAdvisorService{
		languageModel:      languageModel,
		repo:               repo,
		semesterPlanMemory: memoryplan,
	}
}

func (s *SemesterAdvisorService) FormatProfessorReviews(reviews []domain.ProfessorReview) string {
	if len(reviews) == 0 {
		return ""
	}

	var formatted string

	for _, review := range reviews {
		formatted += fmt.Sprintf("Professor: %s\n", review.Professor)
		formatted += fmt.Sprintf("Course: %s\n", review.CourseID)
		formatted += fmt.Sprintf("Department: %s\n", review.Department)
		formatted += fmt.Sprintf("Quality: %.1f\n", review.Quality)
		formatted += fmt.Sprintf("Difficulty: %.1f\n", review.Difficulty)
		formatted += fmt.Sprintf("Grade: %s\n", review.Grade)

		if review.ForCredit != nil {
			formatted += fmt.Sprintf("For Credit: %t\n", *review.ForCredit)
		}

		if review.WouldTakeAgain != nil {
			formatted += fmt.Sprintf("Would Take Again: %t\n", *review.WouldTakeAgain)
		}

		if review.Textbook != nil {
			formatted += fmt.Sprintf("Textbook: %t\n", *review.Textbook)
		}

		formatted += fmt.Sprintf("Comment: %s\n\n", review.Comment)
	}

	return formatted
}

func (s *SemesterAdvisorService) AnalyzeSemester(ctx context.Context, courses []domain.Course) (string, string, error) {

	// get feedback from db
	reviews, err := s.repo.GetReviewsForCourses(courses)
	if err != nil {
		return "", "", err
	}
	// format reviews
	formmattedReviews := s.FormatProfessorReviews(reviews)
	courseData, err := json.Marshal(courses)
	if err != nil {
		return "", "", err
	}
	// build prompt
	prompt, err := prompts.BuildSemesterAdvisorPrompt(courseData, []byte(formmattedReviews))

	answer, err := s.languageModel.Generate(ctx, prompt)
	if err != nil {
		return "", "", err
	}
	// save anwer on redis
	planID := uuid.NewString()
	if err := s.semesterPlanMemory.StorePlan(ctx, planID, courses); err != nil {
		return "", "", err
	}
	return answer, planID, nil
}
