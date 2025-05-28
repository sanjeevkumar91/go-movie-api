package main

import (
	"log"

	"go-movie-api/configs"
	"go-movie-api/movies/client"
	"go-movie-api/movies/constants"
	db "go-movie-api/movies/db"
	"go-movie-api/movies/repository"
	"go-movie-api/movies/service"

	"go-movie-api/movies/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	config := configs.NewConfig()
	configs.LoadConfig(config, constants.ConfigFilePath)
	dbInstance := db.InitDB()
	movieRepository := repository.NewMovieRepository(dbInstance)

	client := client.NewClient(config)
	movieService := service.NewMovieService(client, movieRepository)
	moviesController := controllers.NewMoviesController(movieService)

	router.GET("/", moviesController.SendMessage)
	port := config.GetPort()

	moviesGroup := router.Group("/movies")
	{
		moviesGroup.POST("/search", moviesController.SearchMovies)
		moviesGroup.POST("/", moviesController.GetMovieDetails)
		moviesGroup.POST("/cart", moviesController.AddToMovieCart)
		moviesGroup.GET("/cart", moviesController.GetMoviesInCart)
	}

	if port == "" {
		port = "8080" // if port is not defined in config fallback to default port
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
