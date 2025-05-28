package service

import (
	"errors"
	"go-movie-api/movies/client"
	"go-movie-api/movies/model"
	"go-movie-api/movies/repository"
	"log"

	"github.com/gin-gonic/gin"
)

type movieService struct {
	client     client.Client
	repository repository.MovieRespository
}

type MovieService interface {
	SearchMovies(ctx *gin.Context, req model.SearchMovieRequest) (resp []model.Movie, err error)
	GetMovieDetails(ctx *gin.Context, req model.GetMovieDetailsRequest) (resp model.GetMovieDetailsResponse, err error)
	AddMovieToCart(ctx *gin.Context, req model.AddMovieToCartRequest) (err error)
	GetMoviesInCart(ctx *gin.Context, req model.GetMoviesInCartReq) (movies []model.MovieDetailsInCart, err error)
}

func NewMovieService(client client.Client, repository repository.MovieRespository) movieService {
	return movieService{client: client, repository: repository}
}

func (ms movieService) SearchMovies(ctx *gin.Context, req model.SearchMovieRequest) (movies []model.Movie, err error) {
	resp, err := ms.client.SearchMovies(ctx, req)

	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}

	return resp.Movies, nil
}

func (ms movieService) GetMovieDetails(ctx *gin.Context, req model.GetMovieDetailsRequest) (movieDetails model.GetMovieDetailsResponse, err error) {
	resp, err := ms.client.GetMovieDetails(ctx, req)

	if err != nil {
		return model.GetMovieDetailsResponse{}, err
	}

	if resp.Error != "" {
		return model.GetMovieDetailsResponse{}, errors.New(resp.Error)
	}

	return resp, nil
}

func (ms movieService) AddMovieToCart(ctx *gin.Context, req model.AddMovieToCartRequest) (err error) {
	resp, err := ms.client.GetMovieDetailsById(ctx, req)
	if err != nil {
		return err
	}

	if err := ms.repository.AddToMovieCart(resp, req.UserID); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (ms movieService) GetMoviesInCart(ctx *gin.Context, req model.GetMoviesInCartReq) (movies []model.MovieDetailsInCart, err error) {
	movies, dbErr := ms.repository.GetMoviesInCart(req.UserID)
	if dbErr != nil {
		log.Println(dbErr)
		return nil, dbErr
	}

	return movies, nil
}
