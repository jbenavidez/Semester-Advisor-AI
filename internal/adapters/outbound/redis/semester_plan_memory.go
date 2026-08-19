package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"semester-advisor-ai/internal/domain"

	redisclient "github.com/redis/go-redis/v9"
)

const semesterPlanExpireTime = 24 * time.Hour

type SemesterPlanMemory struct {
	client *redisclient.Client
}

func NewSemesterPlanMemory(client *redisclient.Client) *SemesterPlanMemory {
	return &SemesterPlanMemory{
		client: client,
	}
}

func (m *SemesterPlanMemory) StorePlan(ctx context.Context, planID string, courses []domain.Course) error {
	key := fmt.Sprintf("semester-plan:%s", planID)

	data, err := json.Marshal(courses)
	if err != nil {
		return fmt.Errorf("failed to serialize semester plan: %w", err)
	}

	if err := m.client.Set(ctx, key, data, semesterPlanExpireTime).Err(); err != nil {
		return fmt.Errorf("failed to store semester plan in redis: %w", err)
	}

	return nil
}

func (m *SemesterPlanMemory) GetPlan(ctx context.Context, planID string) ([]domain.Course, error) {
	key := fmt.Sprintf("semester-plan:%s", planID)

	data, err := m.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redisclient.Nil {
			return nil, fmt.Errorf("semester plan not found")
		}

		return nil, fmt.Errorf("failed to retrieve semester plan from redis: %w", err)
	}

	var courses []domain.Course

	if err := json.Unmarshal(data, &courses); err != nil {
		return nil, fmt.Errorf("failed to deserialize semester plan: %w", err)
	}

	return courses, nil
}
