package configs_test

import (
	"go-movie-api/configs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	conf := configs.NewConfig()
	conf.Port = "8080"
	conf.ApiKey = "my-secret"
	conf.MoviesListUrl = "http://www.omdbapi.com/search"

	assert.Equal(t, "8080", conf.GetPort())
	assert.Equal(t, "my-secret", conf.GetApiKey())
	assert.Equal(t, "http://www.omdbapi.com/search", conf.SearchMoviesUrl())
}

func TestLoadConfig(t *testing.T) {
	tempFile := filepath.Join(os.TempDir(), "config_test.json")
	configJSON := `{
		"port": "9000",
		"api_key": "dummy-key",
		"get_movie_list_url": "http://mock-api/movies"
	}`

	err := os.WriteFile(tempFile, []byte(configJSON), 0644)
	assert.NoError(t, err)

	defer func() {
		_ = os.Remove(tempFile)
	}()

	conf := configs.NewConfig()
	configs.LoadConfig(conf, tempFile)

	assert.Equal(t, "9000", conf.GetPort())
	assert.Equal(t, "dummy-key", conf.GetApiKey())
	assert.Equal(t, "http://mock-api/movies", conf.SearchMoviesUrl())
}
