package processing

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	processingdto "semester-advisor-ai/internal/adapters/outbound/processing/dto"
	"semester-advisor-ai/internal/domain"
	"semester-advisor-ai/internal/ports"
	"strconv"
	"strings"
)

type JSONProcessor struct {
	repo ports.UploadedFileRepository
}

func NewJSONProcessor(repo ports.UploadedFileRepository) ports.DatasetProcessor {
	return &JSONProcessor{
		repo: repo,
	}
}

func (j *JSONProcessor) Process(ctx context.Context, reader io.Reader, uploadedFile *domain.UploadedFile) error {
	var reviewGroups [][]processingdto.ProfessorReview

	if err := json.NewDecoder(reader).Decode(&reviewGroups); err != nil {
		return fmt.Errorf("failed to decode professor reviews: %w", err)
	}

	uploadedFile.TotalReviews = 0
	uploadedFile.ProcessedReviews = 0
	uploadedFile.FailedReviews = 0

	for _, group := range reviewGroups {
		uploadedFile.TotalReviews += len(group)
	}

	for _, group := range reviewGroups {
		for _, review := range group {
			if err := ctx.Err(); err != nil {
				return fmt.Errorf("professor review processing canceled: %w", err)
			}

			professor := strings.TrimSpace(review.Professor)
			courseID := strings.TrimSpace(review.CourseID)

			quality, err := strconv.ParseFloat(strings.TrimSpace(review.Quality), 64)
			if err != nil {
				uploadedFile.FailedReviews++
				return fmt.Errorf("invalid quality %q for professor %q and course %q: %w", review.Quality, professor, courseID, err)
			}

			difficulty, err := strconv.ParseFloat(strings.TrimSpace(review.Difficulty), 64)
			if err != nil {
				uploadedFile.FailedReviews++
				return fmt.Errorf("invalid difficulty %q for professor %q and course %q: %w", review.Difficulty, professor, courseID, err)
			}

			forCredit, err := parseOptionalBoolean(review.ForCredit)
			if err != nil {
				uploadedFile.FailedReviews++
				return fmt.Errorf("invalid for-credit value for professor %q and course %q: %w", professor, courseID, err)
			}

			wouldTakeAgain, err := parseOptionalBoolean(review.WouldTakeAgain)
			if err != nil {
				uploadedFile.FailedReviews++
				return fmt.Errorf("invalid would-take-again value for professor %q and course %q: %w", professor, courseID, err)
			}

			textbook, err := parseOptionalBoolean(review.Textbook)
			if err != nil {
				uploadedFile.FailedReviews++
				return fmt.Errorf("invalid textbook value for professor %q and course %q: %w", professor, courseID, err)
			}

			professorReview := &domain.ProfessorReview{
				UploadedFileID: uploadedFile.ID, CourseID: courseID, Quality: quality, Difficulty: difficulty,
				ForCredit: forCredit, WouldTakeAgain: wouldTakeAgain, Grade: strings.TrimSpace(review.Grade),
				Textbook: textbook, Comment: strings.TrimSpace(review.Comment), Professor: professor,
				Department: strings.TrimSpace(review.Department),
			}

			if err := j.repo.SaveReview(ctx, professorReview); err != nil {
				uploadedFile.FailedReviews++
				return fmt.Errorf("failed to save review for professor %q and course %q: %w", professor, courseID, err)
			}

			uploadedFile.ProcessedReviews++
		}
	}

	return nil
}

func parseOptionalBoolean(value string) (*bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "":
		return nil, nil
	case "yes", "true", "1":
		result := true
		return &result, nil
	case "no", "false", "0":
		result := false
		return &result, nil
	default:
		return nil, fmt.Errorf("invalid boolean value %q", value)
	}
}
