package main

import (
	"semester-advisor-ai/internal/routes"

	"github.com/gin-gonic/gin"
)

const (
	portNumber = ":8080"
)

func main() {

	server := gin.Default()
	server.LoadHTMLGlob("templates/*.html")
	routes.SetUpReoutes(server)

	if err := server.Run(":8080"); err != nil {
		panic(err)
	}

}
