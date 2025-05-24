package configs

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Port string `json:"port"`
}

func LoadConfig() Config {
	file, err := os.Open("configs/config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatal(err)
	}
	return config
}
