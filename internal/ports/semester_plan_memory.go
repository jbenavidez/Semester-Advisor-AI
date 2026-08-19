package ports

import (
	"context"

	"semester-advisor-ai/internal/domain"
)

type SemesterPlanMemory interface {
	StorePlan(ctx context.Context, planID string, courses []domain.Course) error
	GetPlan(ctx context.Context, planID string) ([]domain.Course, error)
}
