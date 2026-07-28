package routes

import (
	"semester-advisor-ai/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetUpReoutes(router *gin.Engine) {

	router.GET("/admin/documents/upload", handlers.UploadDoc)

}
