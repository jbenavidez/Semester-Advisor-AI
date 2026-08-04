package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) SemesterPlan(c *gin.Context) {
	c.HTML(http.StatusOK, "semester_plan.html", gin.H{
		"PageTitle": "Build Your Semester Plan",
	})

}
