package routes

import (
	"semester-advisor-ai/internal/adapters/inbound/http/handlers"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine, appHandlers *handlers.Handlers) {
	router.Static("/static", "./static")
	router.GET("/", appHandlers.SemesterPlan)
	router.POST("/plan-semester", appHandlers.PlanSemester)
	router.POST("/plan-semester/alternative", appHandlers.PlanSemesterAlternative)
	router.GET("/admin/documents", appHandlers.GetDocuments)
	router.GET("/admin/documents/upload", appHandlers.UploadDoc)
	router.POST("/admin/documents/upload", appHandlers.ProcesssDoc)

}
