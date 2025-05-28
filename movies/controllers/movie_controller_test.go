package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	mock_service "go-movie-api/movies/mock"
	"go-movie-api/movies/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupRouter(ctrl *gomock.Controller) (*gin.Engine, *mock_service.MockMovieService) {
	mockService := mock_service.NewMockMovieService(ctrl)
	controller := NewMoviesController(mockService)

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/", controller.SendMessage)
	r.POST("/search", controller.SearchMovies)
	r.POST("/details", controller.GetMovieDetails)
	r.POST("/cart", controller.AddToMovieCart)
	r.GET("/cart", controller.GetMoviesInCart)

	return r, mockService
}

func TestSendMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gin.SetMode(gin.TestMode)

	router, _ := setupRouter(ctrl)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "Hello, World!", resp.Body.String())
}

func TestSearchMovies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router, mockService := setupRouter(ctrl)

	t.Run("should return movies on calling search movies endpoint", func(t *testing.T) {
		reqBody := model.SearchMovieRequest{SearchQuery: "Batman"}
		expectedResp := []model.Movie{{Title: "Batman Begins"}}

		mockService.EXPECT().
			SearchMovies(gomock.Any(), reqBody).
			Return(expectedResp, nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/search", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("should return bad request error when invalid request is passed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/search", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("should return internal server error error when movies end point is failing for any reason", func(t *testing.T) {
		reqBody := model.SearchMovieRequest{SearchQuery: "Batman"}

		mockService.EXPECT().
			SearchMovies(gomock.Any(), reqBody).
			Return(nil, errors.New("service failed"))

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/search", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestGetMovieDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router, mockService := setupRouter(ctrl)

	t.Run("should return movie details on calling get movies details endpoint", func(t *testing.T) {
		reqBody := model.GetMovieDetailsRequest{MovieID: "tt1375666"}
		expectedResp := model.GetMovieDetailsResponse{
			Title: "Inception", Year: "2010", Genre: "Sci-Fi", ImdbID: "tt1375666",
		}

		mockService.EXPECT().
			GetMovieDetails(gomock.Any(), reqBody).
			Return(expectedResp, nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/details", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("should return bad request error when invalid request is passed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/details", bytes.NewBufferString("bad json"))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("should return internal server error error when movies end point is failing for any reason", func(t *testing.T) {
		reqBody := model.GetMovieDetailsRequest{MovieID: "tt1375666"}

		mockService.EXPECT().
			GetMovieDetails(gomock.Any(), reqBody).
			Return(model.GetMovieDetailsResponse{}, errors.New("not found"))

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/details", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestAddToMovieCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router, mockService := setupRouter(ctrl)

	t.Run("should add a movie to the cart when movie id is passed", func(t *testing.T) {
		reqBody := model.AddMovieToCartRequest{MovieID: "tt1375666"}

		mockService.EXPECT().
			AddMovieToCart(gomock.Any(), reqBody).
			Return(nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/cart", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("should return bad request error when invalid request is passed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/cart", bytes.NewBufferString("invalid"))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("should return internal server error error when movies end point is failing for any reason", func(t *testing.T) {
		reqBody := model.AddMovieToCartRequest{MovieID: "tt1234567"}

		mockService.EXPECT().
			AddMovieToCart(gomock.Any(), reqBody).
			Return(errors.New("cart error"))

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/cart", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestGetMoviesInCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router, mockService := setupRouter(ctrl)

	t.Run("should return the movies details in the cart", func(t *testing.T) {
		expected := []model.GetMovieDetailsResponse{
			{Title: "Interstellar", ImdbID: "tt0816692"},
		}

		mockService.EXPECT().
			GetMoviesInCart(gomock.Any()).
			Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/cart", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("should return internal server error error when movies end point is failing for any reason", func(t *testing.T) {
		mockService.EXPECT().
			GetMoviesInCart(gomock.Any()).
			Return(nil, errors.New("db failure"))

		req := httptest.NewRequest(http.MethodGet, "/cart", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}
