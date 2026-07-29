package weaviate

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

func NewWeaviateClient() (*weaviate.Client, error) {
	ctx := context.Background()

	weaviateURL := os.Getenv("WEAVIATE_URL")
	if weaviateURL == "" {
		weaviateURL = "http://localhost:8081"
	}

	parsedURL, err := url.Parse(weaviateURL)
	if err != nil {
		return nil, fmt.Errorf("invalid WEAVIATE_URL: %w", err)
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("invalid WEAVIATE_URL: %s", weaviateURL)
	}

	client := weaviate.New(weaviate.Config{
		Scheme: parsedURL.Scheme,
		Host:   parsedURL.Host,
	})

	const maxRetries = 10

	var readinessErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		_, readinessErr = client.Schema().Getter().Do(ctx)
		if readinessErr == nil {
			break
		}

		fmt.Printf("Waiting for Weaviate to be ready... attempt %d/%d\n", attempt, maxRetries)
		time.Sleep(2 * time.Second)
	}

	if readinessErr != nil {
		return nil, fmt.Errorf("weaviate not ready after %d attempts: %w", maxRetries, readinessErr)
	}

	professorReviewClassName := "ProfessorReview"

	professorReviewExists, err := client.Schema().
		ClassExistenceChecker().
		WithClassName(professorReviewClassName).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check ProfessorReview class existence: %w", err)
	}

	if !professorReviewExists {
		fmt.Println("Creating class for professor reviews")

		professorReviewClass := &models.Class{
			Class:       professorReviewClassName,
			Description: "Professor and course reviews used by Semester Advisor AI",
			Vectorizer:  "text2vec-openai",
			Properties: []*models.Property{
				{Name: "professor", DataType: schema.DataTypeText.PropString()},
				{Name: "normalizedProfessor", DataType: schema.DataTypeText.PropString()},
				{Name: "department", DataType: schema.DataTypeText.PropString()},
				{Name: "courseId", DataType: schema.DataTypeText.PropString()},
				{Name: "quality", DataType: schema.DataTypeNumber.PropString()},
				{Name: "difficulty", DataType: schema.DataTypeNumber.PropString()},
				{Name: "forCredit", DataType: schema.DataTypeBoolean.PropString()},
				{Name: "wouldTakeAgain", DataType: schema.DataTypeBoolean.PropString()},
				{Name: "grade", DataType: schema.DataTypeText.PropString()},
				{Name: "textbook", DataType: schema.DataTypeBoolean.PropString()},
				{Name: "comment", DataType: schema.DataTypeText.PropString()},
			},
		}

		if err := client.Schema().
			ClassCreator().
			WithClass(professorReviewClass).
			Do(ctx); err != nil {
			return nil, fmt.Errorf("failed to create ProfessorReview class: %w", err)
		}

		fmt.Println("ProfessorReview class created")
	}

	uploadedFileClassName := "UploadedFile"

	uploadedFileExists, err := client.Schema().
		ClassExistenceChecker().
		WithClassName(uploadedFileClassName).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check UploadedFile class existence: %w", err)
	}

	if !uploadedFileExists {
		fmt.Println("Creating class for uploaded files")
		uploadedFileClass := &models.Class{
			Class:       uploadedFileClassName,
			Description: "Uploaded dataset metadata and import processing status",
			Vectorizer:  "none",
			Properties: []*models.Property{
				{Name: "originalFileName", DataType: schema.DataTypeText.PropString()},
				{Name: "storedFileName", DataType: schema.DataTypeText.PropString()},
				{Name: "filePath", DataType: schema.DataTypeText.PropString()},
				{Name: "description", DataType: schema.DataTypeText.PropString()},
				{Name: "contentType", DataType: schema.DataTypeText.PropString()},
				{Name: "size", DataType: schema.DataTypeInt.PropString()},
				{Name: "datasetName", DataType: schema.DataTypeText.PropString()},
				{Name: "sourcePeriod", DataType: schema.DataTypeText.PropString()},
				{Name: "status", DataType: schema.DataTypeText.PropString()},
				{Name: "errorMessage", DataType: schema.DataTypeText.PropString()},
				{Name: "totalReviews", DataType: schema.DataTypeInt.PropString()},
				{Name: "processedReviews", DataType: schema.DataTypeInt.PropString()},
				{Name: "failedReviews", DataType: schema.DataTypeInt.PropString()},
				{Name: "createdAt", DataType: schema.DataTypeDate.PropString()},
				{Name: "updatedAt", DataType: schema.DataTypeDate.PropString()},
			},
		}

		if err := client.Schema().
			ClassCreator().
			WithClass(uploadedFileClass).
			Do(ctx); err != nil {
			return nil, fmt.Errorf("failed to create UploadedFile class: %w", err)
		}
		fmt.Println("UploadedFile class created")
	}

	fmt.Println("Weaviate DB is ready")
	return client, nil
}
