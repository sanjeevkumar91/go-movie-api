package controllers

import (
	"go-movie-api/movies/model"
	"go-movie-api/movies/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type moviesController struct {
	movieService service.MovieService
}

type MoviesController interface {
	SendMessage(c *gin.Context)
	SearchMovies(c *gin.Context)
	GetMovieDetails(c *gin.Context)
	AddToMovieCart(c *gin.Context)
	GetMoviesInCart(c *gin.Context)
}

func NewMoviesController(movieService service.MovieService) MoviesController {
	return moviesController{movieService: movieService}
}

func (mc moviesController) SendMessage(c *gin.Context) {
	c.String(200, "Hello, World!")
}

func (mc moviesController) SearchMovies(ctx *gin.Context) {
	var movieReq model.SearchMovieRequest
	if err := ctx.ShouldBindJSON(&movieReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("request is valid")

	resp, err := mc.movieService.SearchMovies(ctx, movieReq)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

func (mc moviesController) GetMovieDetails(ctx *gin.Context) {
	var movieReq model.GetMovieDetailsRequest
	if err := ctx.ShouldBindJSON(&movieReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("request is valid")

	resp, err := mc.movieService.GetMovieDetails(ctx, movieReq)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

func (mc moviesController) AddToMovieCart(ctx *gin.Context) {
	var addMovieToCartReq model.AddMovieToCartRequest
	if err := ctx.ShouldBindJSON(&addMovieToCartReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("request is valid")

	err := mc.movieService.AddMovieToCart(ctx, addMovieToCartReq)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, model.AddMovieToCartResponse{Status: "Success"})
}

func (mc moviesController) GetMoviesInCart(ctx *gin.Context) {
	resp, err := mc.movieService.GetMoviesInCart(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}
