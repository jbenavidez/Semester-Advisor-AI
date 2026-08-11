package handlers

import (
	"fmt"
	"net/http"
	"semester-advisor-ai/internal/adapters/inbound/http/dto"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) SemesterPlan(c *gin.Context) {
	c.HTML(http.StatusOK, "semester_plan.html", gin.H{
		"PageTitle": "Build Your Semester Plan",
	})

}

func (h *Handlers) PlanSemester(c *gin.Context) {
	// get courses
	courseIDs := c.PostFormArray("course_id")
	courseNames := c.PostFormArray("course_name")
	professors := c.PostFormArray("professor")
	departments := c.PostFormArray("department")
	credits := c.PostFormArray("credits")
	courseTypes := c.PostFormArray("course_type")
	notes := c.PostFormArray("course_notes")

	// create list
	courses := make([]dto.Course, len(courseIDs))

	for i := range courseIDs {
		credit, _ := strconv.ParseFloat(credits[i], 64)
		course := dto.Course{
			CourseID:      courseIDs[i],
			CourseName:    courseNames[i],
			ProfessorName: professors[i],
			Department:    departments[i],
			Credits:       credit,
			CourseType:    courseTypes[i],
			Notes:         notes[i],
		}
		courses[i] = course
	}
	//
	plan := h.plannerSevice.PlanSemester(courses)
	fmt.Println("complete_gondor", plan)

	// Todo: return result

}
