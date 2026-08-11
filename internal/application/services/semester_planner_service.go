package services

import (
	"fmt"
	"semester-advisor-ai/internal/adapters/inbound/http/dto"
	"semester-advisor-ai/internal/domain"
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

func (sm *SemesterPlannerServices) PlanSemester(courses []dto.Course) error {
	fmt.Println("hello_gondor", len(courses))
	//convert dto to domain
	selectedCourses := make([]domain.Course, len(courses))
	for i, c := range courses {
		course := domain.Course{
			CourseID:      c.CourseID,
			CourseName:    c.CourseName,
			ProfessorName: c.ProfessorName,
			Department:    c.Department,
			Credits:       c.Credits,
			CourseType:    c.CourseType,
			Notes:         c.Notes,
		}
		selectedCourses[i] = course
	}
	// Get reviews
	reviews, err := sm.repo.GetReviewsForCourses(selectedCourses)
	if err != nil {
		return err
	}
	fmt.Println(reviews)
	// Todo: Pass Reviews to LLM
	return nil
}
