package configs

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	Port          string `json:"port"`
	ApiKey        string `json:"api_key"`
	MoviesListUrl string `json:"get_movie_list_url"`
}

type Config interface {
	GetPort() string
	GetApiKey() string
	SearchMoviesUrl() string
}

func NewConfig() *config {
	return &config{}
}

func (c *config) GetPort() string {
	return c.Port
}

func (c *config) GetApiKey() string {
	return c.ApiKey
}

func (c *config) SearchMoviesUrl() string {
	return c.MoviesListUrl
}

func LoadConfig(config *config, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Could not open the config.json", err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatal("Could not load the config", err)
	}
}
