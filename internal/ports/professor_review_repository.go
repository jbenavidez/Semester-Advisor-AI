package ports

import (
	"semester-advisor-ai/internal/domain"
)

type ProfessorReviewRepository interface {
	GetReviewsForCourses(courses []domain.Course) ([]domain.ProfessorReview, error)
}
