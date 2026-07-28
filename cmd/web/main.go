package main

import (
	"semester-advisor-ai/internal/db"
	"semester-advisor-ai/internal/handlers"
	dbrepo "semester-advisor-ai/internal/repository/db_repo"
	"semester-advisor-ai/internal/routes"
	"semester-advisor-ai/internal/services"

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
	weaviateClient, err := db.NewWeaviateClient()
	if err != nil {
		panic(err)
	}
	// wire evertyhing up
	weaviateRepo := dbrepo.NewWeaviateDBRepo(weaviateClient)
	uploadService := services.NewUploadServices(weaviateRepo, chunkSize)
	handlers := handlers.New(uploadService)
	routes.SetUpRoutes(server, handlers)

	if err := server.Run(":8080"); err != nil {
		panic(err)
	}

}
