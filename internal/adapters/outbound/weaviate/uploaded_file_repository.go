package weaviate

import (
	"context"
	"fmt"
	"semester-advisor-ai/internal/domain"
	"semester-advisor-ai/internal/ports"
	"time"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
)

const (
	timeout                   = time.Second * 3
	UploadedFilesClassName    = "UploadedFile"
	DocumentClassName         = "Document"
	ProfessorReviewsClassName = "ProfessorReview"
)

type WeaviateDBRepo struct {
	DB *weaviate.Client
}

func NewWeaviateDBRepo(db *weaviate.Client) ports.UploadedFileRepository {
	return &WeaviateDBRepo{
		DB: db,
	}

}

func (m *WeaviateDBRepo) SaveFile() {
	// to be continue
}

func (m *WeaviateDBRepo) GetFileByName(fileName string) (*domain.UploadedFile, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	where := filters.Where().WithPath([]string{"originalFileName"}).WithOperator(filters.Equal).WithValueString(fileName)

	res, err := m.DB.GraphQL().Get().
		WithClassName(UploadedFilesClassName).
		WithFields(
			graphql.Field{Name: "originalFileName"},
			graphql.Field{Name: "storedFileName"},
			graphql.Field{Name: "filePath"},
			graphql.Field{Name: "description"},
			graphql.Field{Name: "contentType"},
			graphql.Field{Name: "size"},
			graphql.Field{Name: "datasetName"},
			graphql.Field{Name: "sourcePeriod"},
			graphql.Field{Name: "status"},
			graphql.Field{Name: "errorMessage"},
			graphql.Field{Name: "totalReviews"},
			graphql.Field{Name: "processedReviews"},
			graphql.Field{Name: "failedReviews"},
			graphql.Field{Name: "createdAt"},
			graphql.Field{Name: "updatedAt"},
			graphql.Field{Name: "_additional", Fields: []graphql.Field{{Name: "id"}}},
		).
		WithWhere(where).
		WithLimit(1).
		Do(ctx)

	if err != nil {
		return nil, false, err
	}

	data, ok := res.Data["Get"].(map[string]interface{})
	if !ok {
		return nil, false, fmt.Errorf("invalid weaviate response: missing Get")
	}

	items, ok := data[UploadedFilesClassName].([]interface{})
	if !ok || len(items) == 0 {
		return nil, false, nil
	}

	props, ok := items[0].(map[string]interface{})
	if !ok {
		return nil, false, fmt.Errorf("invalid uploaded file properties")
	}

	uploadedFile := &domain.UploadedFile{
		ID:               getGraphQLID(props),
		OriginalFileName: getStringProperty(props, "originalFileName"),
		StoredFileName:   getStringProperty(props, "storedFileName"),
		FilePath:         getStringProperty(props, "filePath"),
		Description:      getStringProperty(props, "description"),
		ContentType:      getStringProperty(props, "contentType"),
		Size:             getInt64Property(props, "size"),
		DatasetName:      getStringProperty(props, "datasetName"),
		SourcePeriod:     getStringProperty(props, "sourcePeriod"),
		Status:           getStringProperty(props, "status"),
		ErrorMessage:     getStringProperty(props, "errorMessage"),
		TotalReviews:     int(getInt64Property(props, "totalReviews")),
		ProcessedReviews: int(getInt64Property(props, "processedReviews")),
		FailedReviews:    int(getInt64Property(props, "failedReviews")),
		CreatedAt:        getTimeProperty(props, "createdAt"),
		UpdatedAt:        getTimeProperty(props, "updatedAt"),
	}

	return uploadedFile, true, nil
}

func (m *WeaviateDBRepo) SaveFileMetaData(uploadedFile *domain.UploadedFile) (*domain.UploadedFile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fileProperties := map[string]interface{}{
		"originalFileName": uploadedFile.OriginalFileName,
		"storedFileName":   uploadedFile.StoredFileName,
		"filePath":         uploadedFile.FilePath,
		"description":      uploadedFile.Description,
		"contentType":      uploadedFile.ContentType,
		"size":             uploadedFile.Size,
		"datasetName":      uploadedFile.DatasetName,
		"sourcePeriod":     uploadedFile.SourcePeriod,
		"status":           uploadedFile.Status,
		"errorMessage":     uploadedFile.ErrorMessage,
		"totalReviews":     uploadedFile.TotalReviews,
		"processedReviews": uploadedFile.ProcessedReviews,
		"failedReviews":    uploadedFile.FailedReviews,
		"createdAt":        uploadedFile.CreatedAt.UTC().Format(time.RFC3339),
		"updatedAt":        uploadedFile.UpdatedAt.UTC().Format(time.RFC3339),
	}

	resp, err := m.DB.Data().Creator().
		WithClassName(UploadedFilesClassName).
		WithProperties(fileProperties).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to save uploaded file metadata in Weaviate: %w", err)
	}

	uploadedFile.ID = resp.Object.ID.String()

	return uploadedFile, nil
}

func (m *WeaviateDBRepo) GetAllUploadFiles() ([]domain.UploadedFile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := m.DB.Data().ObjectsGetter().
		WithClassName(UploadedFilesClassName).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve uploaded files from Weaviate: %w", err)
	}

	uploadedFiles := make([]domain.UploadedFile, 0, len(res))

	for i := range res {
		obj := res[i]

		props, ok := obj.Properties.(map[string]interface{})
		if !ok {
			continue
		}

		uploadedFile := domain.UploadedFile{
			ID:               obj.ID.String(),
			OriginalFileName: getStringProperty(props, "originalFileName"),
			StoredFileName:   getStringProperty(props, "storedFileName"),
			FilePath:         getStringProperty(props, "filePath"),
			Description:      getStringProperty(props, "description"),
			ContentType:      getStringProperty(props, "contentType"),
			Size:             getInt64Property(props, "size"),
			DatasetName:      getStringProperty(props, "datasetName"),
			SourcePeriod:     getStringProperty(props, "sourcePeriod"),
			Status:           getStringProperty(props, "status"),
			ErrorMessage:     getStringProperty(props, "errorMessage"),
			TotalReviews:     getIntProperty(props, "totalReviews"),
			ProcessedReviews: getIntProperty(props, "processedReviews"),
			FailedReviews:    getIntProperty(props, "failedReviews"),
			CreatedAt:        getTimeProperty(props, "createdAt"),
			UpdatedAt:        getTimeProperty(props, "updatedAt"),
		}

		uploadedFiles = append(uploadedFiles, uploadedFile)
	}

	return uploadedFiles, nil
}

func (m *WeaviateDBRepo) UpdateFile(uploadedFile *domain.UploadedFile) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	uploadedFile.UpdatedAt = time.Now().UTC()

	fileProperties := map[string]interface{}{
		"originalFileName": uploadedFile.OriginalFileName,
		"storedFileName":   uploadedFile.StoredFileName,
		"filePath":         uploadedFile.FilePath,
		"description":      uploadedFile.Description,
		"contentType":      uploadedFile.ContentType,
		"size":             uploadedFile.Size,
		"datasetName":      uploadedFile.DatasetName,
		"sourcePeriod":     uploadedFile.SourcePeriod,
		"status":           uploadedFile.Status,
		"errorMessage":     uploadedFile.ErrorMessage,
		"totalReviews":     uploadedFile.TotalReviews,
		"processedReviews": uploadedFile.ProcessedReviews,
		"failedReviews":    uploadedFile.FailedReviews,
		"updatedAt":        uploadedFile.UpdatedAt.Format(time.RFC3339),
	}

	err := m.DB.Data().Updater().
		WithMerge().
		WithID(uploadedFile.ID).
		WithClassName(UploadedFilesClassName).
		WithProperties(fileProperties).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to update uploaded file: %w", err)
	}

	return nil
}

func (m *WeaviateDBRepo) SaveReview(ctx context.Context, review *domain.ProfessorReview) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	reviewProperties := map[string]interface{}{
		"uploadedFileId": review.UploadedFileID,
		"courseId":       review.CourseID,
		"quality":        review.Quality,
		"difficulty":     review.Difficulty,
		"grade":          review.Grade,
		"comment":        review.Comment,
		"professor":      review.Professor,
		"department":     review.Department,
	}

	if review.ForCredit != nil {
		reviewProperties["forCredit"] = *review.ForCredit
	}

	if review.WouldTakeAgain != nil {
		reviewProperties["wouldTakeAgain"] = *review.WouldTakeAgain
	}

	if review.Textbook != nil {
		reviewProperties["textbook"] = *review.Textbook
	}

	response, err := m.DB.Data().Creator().
		WithClassName(ProfessorReviewsClassName).
		WithProperties(reviewProperties).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to save review for professor %q and course %q: %w", review.Professor, review.CourseID, err)
	}

	review.ID = response.Object.ID.String()

	return nil
}

// HELPERS
func getGraphQLID(props map[string]interface{}) string {
	additional, ok := props["_additional"].(map[string]interface{})
	if !ok {
		return ""
	}

	return getStringProperty(additional, "id")
}

func getStringProperty(props map[string]interface{}, key string) string {
	value, ok := props[key]
	if !ok || value == nil {
		return ""
	}

	strValue, ok := value.(string)
	if !ok {
		return ""
	}

	return strValue
}

func getInt64Property(props map[string]interface{}, key string) int64 {
	value, ok := props[key]
	if !ok || value == nil {
		return 0
	}

	switch v := value.(type) {
	case int:
		return int64(v)
	case int64:
		return v
	case float64:
		return int64(v)
	default:
		return 0
	}
}

func getTimeProperty(props map[string]interface{}, key string) time.Time {
	value, ok := props[key]
	if !ok || value == nil {
		return time.Time{}
	}

	strValue, ok := value.(string)
	if !ok {
		return time.Time{}
	}

	parsedTime, err := time.Parse(time.RFC3339, strValue)
	if err != nil {
		return time.Time{}
	}

	return parsedTime
}

func getIntProperty(props map[string]interface{}, key string) int {
	value, ok := props[key]
	if !ok || value == nil {
		return 0
	}

	switch number := value.(type) {
	case int:
		return number
	case int64:
		return int(number)
	case float64:
		return int(number)
	default:
		return 0
	}
}
