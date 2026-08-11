package main

import (
	"os"
	"semester-advisor-ai/internal/adapters/inbound/http/handlers"
	"semester-advisor-ai/internal/adapters/inbound/http/routes"
	"semester-advisor-ai/internal/adapters/outbound/processing"
	storageadapter "semester-advisor-ai/internal/adapters/outbound/storage"
	weaviateadapter "semester-advisor-ai/internal/adapters/outbound/weaviate"
	"semester-advisor-ai/internal/application/services"

	"github.com/gin-gonic/gin"
)

const (
	portNumber = ":8080"
	chunkSize  = 10
)

func main() {

	server := gin.Default()
	server.LoadHTMLGlob("templates/*.html")
	// connect to db
	weaviateClient, err := weaviateadapter.NewWeaviateClient()
	if err != nil {
		panic(err)
	}
	// wire evertyhing up
	weaviateRepo := weaviateadapter.NewWeaviateDBRepo(weaviateClient)
	uploadDir := os.Getenv("UPLOAD_DIR")
	fileStorage := storageadapter.NewLocalStorage(uploadDir)
	fileProcessor := processing.NewJSONProcessor(weaviateRepo)
	uploadService := services.NewUploadServices(weaviateRepo, fileStorage, fileProcessor, chunkSize)
	plannerServices := services.NewSemesterPlannerServices(weaviateRepo)
	handlers := handlers.New(uploadService, plannerServices)
	routes.SetUpRoutes(server, handlers)

	if err := server.Run(":8080"); err != nil {
		panic(err)
	}

}
