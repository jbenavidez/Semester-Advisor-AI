package handlers

import (
	"net/http"
	"semester-advisor-ai/internal/domain"
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
	courses := make([]domain.Course, len(courseIDs))

	for i := range courseIDs {
		credit, _ := strconv.ParseFloat(credits[i], 64)
		course := domain.Course{
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
	plan, planID, err := h.semesterAdvisor.AnalyzeSemester(c.Request.Context(), courses)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "semester_analysis.html", gin.H{
		"PageTitle": "Semester Analysis",
		"Plan":      plan,
		"PlanID":    planID,
	})

}

func (h *Handlers) PlanSemesterAlternative(c *gin.Context) {
	planID := c.PostForm("plan_id")
	if planID == "" {
		c.String(http.StatusBadRequest, "plan id is required")
		return
	}
	plan, err := h.plannerService.PlanSemester(c.Request.Context(), planID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.HTML(http.StatusOK, "semester_alternatives.html", gin.H{
		"PageTitle": "Better Semester Alternatives",
		"Plan":      plan,
		"PlanID":    planID,
	})
}
