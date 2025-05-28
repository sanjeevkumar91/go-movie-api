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
	userRespository := repository.NewUserRepository(dbInstance)

	client := client.NewClient(config)
	userService := service.NewUserService(userRespository)
	movieService := service.NewMovieService(client, movieRepository)
	moviesController := controllers.NewMoviesController(movieService)
	userController := controllers.NewUserController(userService)

	router.GET("/", moviesController.SendMessage)
	port := config.GetPort()

	usersGroup := router.Group("/users")
	{
		usersGroup.POST("/", userController.CreateUser)
		usersGroup.GET("/", userController.GetUsers)
	}

	moviesGroup := router.Group("/movies")
	{
		moviesGroup.POST("/search", moviesController.SearchMovies)
		moviesGroup.POST("/", moviesController.GetMovieDetails)
		moviesGroup.POST("/cart/add", moviesController.AddToMovieCart)
		moviesGroup.POST("/cart/list", moviesController.GetMoviesInCart)
	}

	if port == "" {
		port = "8080" // if port is not defined in config fallback to default port
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
		panic(err)
	}
}
