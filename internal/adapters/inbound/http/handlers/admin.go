package handlers

import (
	"fmt"
	"net/http"
	"semester-advisor-ai/internal/adapters/inbound/http/dto"

	"github.com/gin-gonic/gin"
)

const (
	maxUploadSize = 10 << 20 // 10 MB
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

func (h *Handlers) ProcesssDoc(c *gin.Context) {
	//limit body to 10 MB.
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	if err := c.Request.ParseMultipartForm(maxUploadSize); err != nil {
		c.String(http.StatusBadRequest, "file too large.")
		return
	}
	// get file
	file, header, err := c.Request.FormFile("dataset")
	if err != nil {
		c.String(http.StatusBadRequest, "file is required.")
		return
	}
	defer func() {
		_ = file.Close()
	}()
	datasetName := c.PostForm("dataset_name")
	sourcePeriod := c.PostForm("source_period")
	description := c.PostForm("description")
	//set dto
	uploadDatasetInput := dto.UploadDatasetInput{
		DatasetName:  datasetName,
		SourcePeriod: sourcePeriod,
		Description:  description,
		File:         file,
		Header:       header,
	}

	// save file
	_, err = h.uploadService.SaveUploadedFile(c.Request.Context(), uploadDatasetInput)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("file Uploaded")
	c.Redirect(http.StatusSeeOther, "/admin/documents")

}

func (h *Handlers) GetDocuments(c *gin.Context) {

	files, err := h.uploadService.GetAllUploadedFile()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("gondor")
	c.HTML(http.StatusOK, "documents.html", gin.H{
		"PageTitle":     "Uploaded documents",
		"UploadedFiles": files.UploadedFiles,
	})
}
