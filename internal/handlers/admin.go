package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UploadDoc renders the upload page
func (h *Handlers) UploadDoc(c *gin.Context) {

	c.HTML(http.StatusOK, "dataset_upload.html", gin.H{
		"PageTitle":       "Upload Data",
		"PageDescription": "Import and process professor review data.",
		"ActivePage":      "upload",
		"DataMenuOpen":    true,
	})

}
