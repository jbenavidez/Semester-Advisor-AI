package weaviate

import (
	"context"
	"fmt"
	"strings"

	"semester-advisor-ai/internal/domain"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
)

func (m *WeaviateDBRepo) GetReviewsForCourses(courses []domain.Course) ([]domain.ProfessorReview, error) {
	if len(courses) == 0 {
		return []domain.ProfessorReview{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	courseFilters := make([]*filters.WhereBuilder, 0, len(courses))

	for _, course := range courses {
		courseID := strings.TrimSpace(course.CourseID)
		professor := strings.TrimSpace(course.ProfessorName)

		courseFilter := filters.Where().
			WithOperator(filters.And).
			WithOperands([]*filters.WhereBuilder{
				filters.Where().WithPath([]string{"courseId"}).WithOperator(filters.Equal).WithValueText(courseID),
				filters.Where().WithPath([]string{"professor"}).WithOperator(filters.Equal).WithValueText(professor),
			})

		courseFilters = append(courseFilters, courseFilter)
	}

	var where *filters.WhereBuilder

	if len(courseFilters) == 1 {
		where = courseFilters[0]
	} else {
		where = filters.Where().WithOperator(filters.Or).WithOperands(courseFilters)
	}

	res, err := m.DB.GraphQL().Get().
		WithClassName(ProfessorReviewsClassName).
		WithFields(
			graphql.Field{Name: "uploadedFileId"},
			graphql.Field{Name: "courseId"},
			graphql.Field{Name: "quality"},
			graphql.Field{Name: "difficulty"},
			graphql.Field{Name: "forCredit"},
			graphql.Field{Name: "wouldTakeAgain"},
			graphql.Field{Name: "grade"},
			graphql.Field{Name: "textbook"},
			graphql.Field{Name: "comment"},
			graphql.Field{Name: "professor"},
			graphql.Field{Name: "department"},
			graphql.Field{Name: "_additional", Fields: []graphql.Field{{Name: "id"}}},
		).
		WithWhere(where).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve professor reviews from Weaviate: %w", err)
	}

	data, ok := res.Data["Get"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid Weaviate response: missing Get")
	}

	items, ok := data[ProfessorReviewsClassName].([]interface{})
	if !ok || len(items) == 0 {
		return []domain.ProfessorReview{}, nil
	}

	reviews := make([]domain.ProfessorReview, 0, len(items))

	for _, item := range items {
		props, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		review := domain.ProfessorReview{
			ID:             getGraphQLID(props),
			UploadedFileID: getStringProperty(props, "uploadedFileId"),
			CourseID:       getStringProperty(props, "courseId"),
			Quality:        getFloat64Property(props, "quality"),
			Difficulty:     getFloat64Property(props, "difficulty"),
			ForCredit:      getOptionalBoolProperty(props, "forCredit"),
			WouldTakeAgain: getOptionalBoolProperty(props, "wouldTakeAgain"),
			Grade:          getStringProperty(props, "grade"),
			Textbook:       getOptionalBoolProperty(props, "textbook"),
			Comment:        getStringProperty(props, "comment"),
			Professor:      getStringProperty(props, "professor"),
			Department:     getStringProperty(props, "department"),
		}

		reviews = append(reviews, review)
	}

	return reviews, nil
}

func getFloat64Property(props map[string]interface{}, key string) float64 {
	value, ok := props[key]
	if !ok || value == nil {
		return 0
	}

	switch number := value.(type) {
	case float64:
		return number
	case float32:
		return float64(number)
	case int:
		return float64(number)
	case int64:
		return float64(number)
	default:
		return 0
	}
}

func getOptionalBoolProperty(props map[string]interface{}, key string) *bool {
	value, ok := props[key]
	if !ok || value == nil {
		return nil
	}

	boolValue, ok := value.(bool)
	if !ok {
		return nil
	}

	return &boolValue
}
