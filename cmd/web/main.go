package main

import (
	"semester-advisor-ai/internal/adapters/inbound/http/handlers"
	"semester-advisor-ai/internal/adapters/inbound/http/routes"
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
	uploadService := services.NewUploadServices(weaviateRepo, chunkSize)
	handlers := handlers.New(uploadService)
	routes.SetUpRoutes(server, handlers)

	if err := server.Run(":8080"); err != nil {
		panic(err)
	}

}
