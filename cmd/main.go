package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sanjeevkumar91/go-movie-api/configs"
	"github.com/sanjeevkumar91/go-movie-api/pkg/movies"
)

func main() {
	router := gin.Default()
	moviesController := movies.NewMoviesController()

	router.GET("/", moviesController.SendMessage)
	port := configs.LoadConfig().Port

	if port == "" {
		port = "8080" // if port is not defined in config fallback to default port
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
