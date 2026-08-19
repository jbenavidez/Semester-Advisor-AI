package main

import (
	"log"
	"os"
	"semester-advisor-ai/internal/adapters/inbound/http/handlers"
	"semester-advisor-ai/internal/adapters/inbound/http/routes"
	"semester-advisor-ai/internal/adapters/outbound/llm"
	"semester-advisor-ai/internal/adapters/outbound/processing"
	"semester-advisor-ai/internal/adapters/outbound/redis"
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
	llmClient, err := llm.NewOllamaClient()
	if err != nil {
		panic(err)
	}
	ragEngine := llm.NewRag(llmClient)
	// init redis
	redisClient, err := redis.NewRedisClient()
	if err != nil {
		log.Fatal(err)
	}
	semesterPlanMemory := redis.NewSemesterPlanMemory(redisClient)
	// wire evertyhing up
	weaviateRepo := weaviateadapter.NewWeaviateDBRepo(weaviateClient)
	uploadDir := os.Getenv("UPLOAD_DIR")
	fileStorage := storageadapter.NewLocalStorage(uploadDir)
	fileProcessor := processing.NewJSONProcessor(weaviateRepo)
	uploadService := services.NewUploadServices(weaviateRepo, fileStorage, fileProcessor, chunkSize)
	semesterAdvisorService := services.NewSemesterAdvisorService(ragEngine, weaviateRepo, semesterPlanMemory)
	plannerServices := services.NewSemesterPlannerServices(weaviateRepo, ragEngine, semesterPlanMemory)
	handlers := handlers.New(uploadService, plannerServices, semesterAdvisorService)
	routes.SetUpRoutes(server, handlers)

	if err := server.Run(":8080"); err != nil {
		panic(err)
	}

}
