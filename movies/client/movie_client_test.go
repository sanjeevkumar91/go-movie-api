package client

import (
	"bytes"
	"go-movie-api/movies/mock"
	"go-movie-api/movies/model"
	"io"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type roundTripperFunc func(*http.Request) *http.Response

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func mockHttpClient(responseBody string, statusCode int) *http.Client {
	return &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
				Header:     make(http.Header),
			}
		}),
	}
}

func TestForSearchMovies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := mock.NewMockConfig(ctrl)

	mockCfg.EXPECT().GetApiKey().AnyTimes().Return("mock-key")

	// Replace http.DefaultClient with a custom client temporarily
	defaultClient := http.DefaultClient
	http.DefaultClient = mockHttpClient(`{"Search":[{"Title":"Inception"}]}`, 200)
	defer func() { http.DefaultClient = defaultClient }()

	c := NewClient(mockCfg)

	t.Run("should return valid response when search movie api is success", func(t *testing.T) {
		mockCfg.EXPECT().SearchMoviesUrl().Return("http://mock-api/movies")
		ctx, _ := gin.CreateTestContext(nil)
		req := model.SearchMovieRequest{SearchQuery: "Inception"}

		resp, err := c.SearchMovies(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.GreaterOrEqual(t, len(resp.Movies), 1)
		assert.Equal(t, "Inception", resp.Movies[0].Title)
	})

	t.Run("should throw error when url is invalid", func(t *testing.T) {
		mockCfg.EXPECT().SearchMoviesUrl().Return("http://%%invalid-url").Times(1)
		ctx, _ := gin.CreateTestContext(nil)
		req := model.SearchMovieRequest{SearchQuery: "fail"}

		resp, err := c.SearchMovies(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, resp.Movies)
	})

	t.Run("should return http request error for invalid url", func(t *testing.T) {
		mockCfg.EXPECT().SearchMoviesUrl().Return("http://[invalid-url").Times(1)
		ctx, _ := gin.CreateTestContext(nil)
		req := model.SearchMovieRequest{SearchQuery: "fail"}

		resp, err := c.SearchMovies(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, resp.Movies)
	})

	t.Run("should return json decode error if api returns invalid json", func(t *testing.T) {
		mockCfg.EXPECT().SearchMoviesUrl().Return("http://mock-api/movies").Times(1)

		http.DefaultClient = mockHttpClient("not json", 200)

		ctx, _ := gin.CreateTestContext(nil)
		req := model.SearchMovieRequest{SearchQuery: "fail"}

		resp, err := c.SearchMovies(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, resp.Movies)
	})
}

func TestForGetMovieDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := mock.NewMockConfig(ctrl)

	mockCfg.EXPECT().GetApiKey().AnyTimes().Return("mock-key")

	defaultClient := http.DefaultClient
	http.DefaultClient = mockHttpClient(`{"Search":[{"Title":"Inception"}]}`, 200)
	defer func() { http.DefaultClient = defaultClient }()

	c := NewClient(mockCfg)

	t.Run("should get the movie details based on the request", func(t *testing.T) {
		mockCfg.EXPECT().SearchMoviesUrl().Return("http://mock-api/movies")
		http.DefaultClient = mockHttpClient(`{"Title":"Inception","Year":"2010","Genre":"Sci-Fi","ImdbID":"tt1375666"}`, 200)
		ctx, _ := gin.CreateTestContext(nil)
		req := model.GetMovieDetailsRequest{MovieID: "tt1375666"}

		resp, err := c.GetMovieDetails(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, "Inception", resp.Title)
		assert.Equal(t, "tt1375666", resp.ImdbID)
	})

	t.Run("should return json decode error if invalid json is received from the api", func(t *testing.T) {
		mockCfg.EXPECT().SearchMoviesUrl().Return("http://mock-api/movies")
		http.DefaultClient = mockHttpClient("invalid json", 200)
		ctx, _ := gin.CreateTestContext(nil)
		req := model.GetMovieDetailsRequest{MovieID: "tt1375666"}

		resp, err := c.GetMovieDetails(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, resp.Title)
	})
}

func TestForGetMovieDetailsById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := mock.NewMockConfig(ctrl)

	mockCfg.EXPECT().GetApiKey().AnyTimes().Return("mock-key")

	// Replace http.DefaultClient with a custom client temporarily
	defaultClient := http.DefaultClient
	http.DefaultClient = mockHttpClient(`{"Search":[{"Title":"Inception"}]}`, 200)
	defer func() { http.DefaultClient = defaultClient }()

	c := NewClient(mockCfg)

	t.Run("should get the movie details based on the movie id", func(t *testing.T) {
		mockCfg.EXPECT().SearchMoviesUrl().Return("http://mock-api/movies")
		http.DefaultClient = mockHttpClient(`{"Title":"Matrix","Year":"1999","Genre":"Action","ImdbID":"tt0133093"}`, 200)
		ctx, _ := gin.CreateTestContext(nil)
		req := model.AddMovieToCartRequest{MovieID: "tt0133093"}

		resp, err := c.GetMovieDetailsById(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, "Matrix", resp.Title)
		assert.Equal(t, "tt0133093", resp.ImdbID)
	})

	t.Run("should return json decode error if invalid json is received from the api", func(t *testing.T) {
		mockCfg.EXPECT().SearchMoviesUrl().Return("http://mock-api/movies")
		http.DefaultClient = mockHttpClient("invalid json", 200)
		ctx, _ := gin.CreateTestContext(nil)
		req := model.AddMovieToCartRequest{MovieID: "tt0133093"}

		resp, err := c.GetMovieDetailsById(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, resp.Title)
	})
}
