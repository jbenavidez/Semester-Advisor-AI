package routes

import (
	"semester-advisor-ai/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine, appHandlers *handlers.Handlers) {

	router.GET("/admin/documents/upload", appHandlers.UploadDoc)

}
