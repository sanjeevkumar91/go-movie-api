package client

import (
	"encoding/json"
	"go-movie-api/configs"
	"go-movie-api/movies/model"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type Client interface {
	SearchMovies(ctx *gin.Context, request model.SearchMovieRequest) (movie model.SearchMovieResponse, err error)
	GetMovieDetails(ctx *gin.Context, request model.GetMovieDetailsRequest) (movie model.GetMovieDetailsResponse, err error)
	GetMovieDetailsById(ctx *gin.Context, request model.AddMovieToCartRequest) (movie model.GetMovieDetailsResponse, err error)
}

type client struct {
	appConfig configs.Config
}

func NewClient(appConfig configs.Config) client {
	return client{appConfig: appConfig}
}

func (c client) SearchMovies(ctx *gin.Context, request model.SearchMovieRequest) (model.SearchMovieResponse, error) {
	log.Println("initiating movies search req", request)

	queryParams := constructParamsForSearchMovies(request, c.appConfig)

	var searchMovieResponse model.SearchMovieResponse
	if err := makeGetRequest(c.appConfig.SearchMoviesUrl(), queryParams, &searchMovieResponse); err != nil {
		log.Println(err)
		return model.SearchMovieResponse{}, err
	}

	log.Println("fetched response from movies search")

	return searchMovieResponse, nil
}

func (c client) GetMovieDetails(ctx *gin.Context, request model.GetMovieDetailsRequest) (model.GetMovieDetailsResponse, error) {
	log.Println("initiating movie search req", request)

	queryParams := constructParamsForGetMovieDetails(request, c.appConfig)

	var movieDetailsResponse model.GetMovieDetailsResponse
	if err := makeGetRequest(c.appConfig.SearchMoviesUrl(), queryParams, &movieDetailsResponse); err != nil {
		log.Println(err)
		return model.GetMovieDetailsResponse{}, err
	}

	log.Println("fetched the movie details successfully")

	return movieDetailsResponse, nil
}

func (c client) GetMovieDetailsById(ctx *gin.Context, request model.AddMovieToCartRequest) (model.GetMovieDetailsResponse, error) {
	log.Println("initiating movie search req", request)

	queryParams := constructParamsForGetMovieDetailsById(request, c.appConfig)
	var movieDetailsResponse model.GetMovieDetailsResponse
	if err := makeGetRequest(c.appConfig.SearchMoviesUrl(), queryParams, &movieDetailsResponse); err != nil {
		log.Println(err)
		return model.GetMovieDetailsResponse{}, err
	}

	log.Println("fetched movie details for id", request.MovieID)

	return movieDetailsResponse, nil
}

func constructParamsForSearchMovies(request model.SearchMovieRequest, appConfig configs.Config) url.Values {
	params := url.Values{}
	params.Add("apikey", appConfig.GetApiKey())
	params.Add("s", request.SearchQuery)

	log.Println("request", request)

	if request.Title != "" {
		params.Add("t", request.Title)
	}
	if request.Year != "" {
		params.Add("y", request.Year)
	}
	if request.Type != "" {
		params.Add("type", request.Type)
	}
	if request.Page != "" {
		params.Add("page", request.Page)
	}
	return params
}

func constructParamsForGetMovieDetails(request model.GetMovieDetailsRequest, appConfig configs.Config) url.Values {
	params := url.Values{}
	params.Add("apikey", appConfig.GetApiKey())

	if request.MovieID != "" {
		params.Add("i", request.MovieID)
	}
	if request.Title != "" {
		params.Add("t", request.Title)
	}
	if request.Year != "" {
		params.Add("y", request.Year)
	}
	if request.Type != "" {
		params.Add("type", request.Type)
	}
	return params
}

func constructParamsForGetMovieDetailsById(request model.AddMovieToCartRequest, appConfig configs.Config) url.Values {
	params := url.Values{}
	params.Add("apikey", appConfig.GetApiKey())

	if request.MovieID != "" {
		params.Add("i", request.MovieID)
	}
	return params
}

func makeGetRequest(apiUrl string, queryParams url.Values, out any) (err error) {
	u, err := url.Parse(apiUrl)
	if err != nil {
		log.Println(err)
		return err
	}
	u.RawQuery = queryParams.Encode()

	log.Println("movies req str", u.String())

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return err
	}

	return nil
}
