package service

import (
	"errors"
	"go-movie-api/movies/model"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock "go-movie-api/movies/mock"
)

func TestSearchMovies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockClient(ctrl)
	mockRepo := mock.NewMockMovieRespository(ctrl)
	svc := NewMovieService(mockClient, mockRepo)

	ctx := &gin.Context{}

	t.Run("should return movies when api returns success", func(t *testing.T) {
		req := model.SearchMovieRequest{SearchQuery: "Inception"}
		resp := model.SearchMovieResponse{
			Movies: []model.Movie{{Title: "Inception", Year: "2010"}},
			Error:  "",
		}

		mockClient.EXPECT().SearchMovies(ctx, req).Return(resp, nil)

		movies, err := svc.SearchMovies(ctx, req)

		assert.NoError(t, err)
		assert.Len(t, movies, 1)
		assert.Equal(t, "Inception", movies[0].Title)
	})

	t.Run("should return client error when there is client failure", func(t *testing.T) {
		req := model.SearchMovieRequest{SearchQuery: "Inception"}

		mockClient.EXPECT().SearchMovies(ctx, req).Return(model.SearchMovieResponse{}, errors.New("client failure"))

		movies, err := svc.SearchMovies(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, movies)
		assert.Equal(t, "client failure", err.Error())
	})

	t.Run("should return api response error when movie is not found", func(t *testing.T) {
		req := model.SearchMovieRequest{SearchQuery: "Inception"}
		resp := model.SearchMovieResponse{
			Movies: nil,
			Error:  "movie not found",
		}

		mockClient.EXPECT().SearchMovies(ctx, req).Return(resp, nil)

		movies, err := svc.SearchMovies(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, movies)
		assert.Equal(t, "movie not found", err.Error())
	})
}

func TestAddMovieToCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockClient(ctrl)
	mockRepo := mock.NewMockMovieRespository(ctrl)
	svc := NewMovieService(mockClient, mockRepo)
	ctx := &gin.Context{}

	t.Run("should add movie to cart successfully", func(t *testing.T) {
		req := model.AddMovieToCartRequest{
			UserID:  "123",
			MovieID: "tt1375666",
		}

		resp := model.GetMovieDetailsResponse{
			Title:  "Inception",
			Year:   "2010",
			Genre:  "Sci-Fi",
			ImdbID: "tt1375666",
		}

		mockClient.EXPECT().GetMovieDetailsById(ctx, req).Return(resp, nil)
		mockRepo.EXPECT().AddToMovieCart(resp, req.UserID).Return(nil)

		err := svc.AddMovieToCart(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("should return error when client fails", func(t *testing.T) {
		req := model.AddMovieToCartRequest{
			UserID:  "123",
			MovieID: "tt1375666",
		}

		mockClient.EXPECT().GetMovieDetailsById(ctx, req).Return(model.GetMovieDetailsResponse{}, errors.New("client error"))

		err := svc.AddMovieToCart(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, "client error", err.Error())
	})

	t.Run("should return error when repo fails", func(t *testing.T) {
		req := model.AddMovieToCartRequest{
			UserID:  "123",
			MovieID: "tt1375666",
		}

		resp := model.GetMovieDetailsResponse{
			Title:  "Inception",
			Year:   "2010",
			Genre:  "Sci-Fi",
			ImdbID: "tt1375666",
		}

		mockClient.EXPECT().GetMovieDetailsById(ctx, req).Return(resp, nil)
		mockRepo.EXPECT().AddToMovieCart(resp, req.UserID).Return(errors.New("repo error"))

		err := svc.AddMovieToCart(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, "repo error", err.Error())
	})
}

func TestGetMovieDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockClient(ctrl)
	mockRepo := mock.NewMockMovieRespository(ctrl)

	svc := NewMovieService(mockClient, mockRepo)
	ctx := &gin.Context{}

	t.Run("should return movie details on api success", func(t *testing.T) {
		req := model.GetMovieDetailsRequest{MovieID: "tt1375666"}
		resp := model.GetMovieDetailsResponse{
			Title:  "Inception",
			Year:   "2010",
			Genre:  "Sci-Fi",
			ImdbID: "tt1375666",
			Error:  "",
		}

		mockClient.EXPECT().GetMovieDetails(ctx, req).Return(resp, nil)

		result, err := svc.GetMovieDetails(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, resp, result)
	})

	t.Run("should return client error when there is client failure", func(t *testing.T) {
		req := model.GetMovieDetailsRequest{MovieID: "tt1375666"}

		mockClient.EXPECT().GetMovieDetails(ctx, req).Return(model.GetMovieDetailsResponse{}, errors.New("client error"))

		result, err := svc.GetMovieDetails(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, result.Title)
		assert.Equal(t, "client error", err.Error())
	})

	t.Run("should return api error when movie is not found", func(t *testing.T) {
		req := model.GetMovieDetailsRequest{MovieID: "tt1375666"}
		resp := model.GetMovieDetailsResponse{
			Error: "movie not found",
		}

		mockClient.EXPECT().GetMovieDetails(ctx, req).Return(resp, nil)

		result, err := svc.GetMovieDetails(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, result.Title)
		assert.Equal(t, "movie not found", err.Error())
	})
}

func TestGetMoviesInCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockMovieRespository(ctrl)
	mockClient := mock.NewMockClient(ctrl)
	svc := NewMovieService(mockClient, mockRepo)
	ctx := &gin.Context{}

	t.Run("should return movies from database", func(t *testing.T) {
		req := model.GetMoviesInCartReq{UserID: "123"}
		expected := []model.MovieDetailsInCart{
			{Title: "Inception", Year: "2010", Genre: "Sci-Fi", ImdbID: "tt1375666"},
		}

		mockRepo.EXPECT().GetMoviesInCart(req.UserID).Return(expected, nil)

		movies, err := svc.GetMoviesInCart(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, expected, movies)
	})

	t.Run("should return db error when fetch fails", func(t *testing.T) {
		req := model.GetMoviesInCartReq{UserID: "123"}

		mockRepo.EXPECT().GetMoviesInCart(req.UserID).Return(nil, errors.New("db error"))

		movies, err := svc.GetMoviesInCart(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, movies)
		assert.Equal(t, "db error", err.Error())
	})
}
