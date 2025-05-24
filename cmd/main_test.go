package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sanjeevkumar91/go-movie-api/pkg/movies"
	"github.com/stretchr/testify/assert"
)

func TestSendMessageRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	moviesController := movies.NewMoviesController()

	router.GET("/", moviesController.SendMessage)

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err, "Failed to create request")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")
	assert.Equal(t, "Hello, World!", w.Body.String(), "Response body mismatch")
}
