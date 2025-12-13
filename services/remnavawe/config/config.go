package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BindUrl     string
	RemnaAPIKey string
	RemnaURL    string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	bindUrl := os.Getenv("REMNAVAWE_REST_BIND_URL")
	if len(bindUrl) == 0 {
		bindUrl = ":8080"
	}
	remnaAPIKey := os.Getenv("REMNAVAWE_API_KEY")
	if len(remnaAPIKey) == 0 {
		remnaAPIKey = ""
	}
	remnaURL := os.Getenv("REMNAVAWE_URL")
	if len(remnaURL) == 0 {
		remnaURL = ""
	}
	return &Config{
		BindUrl:     bindUrl,
		RemnaAPIKey: remnaAPIKey,
		RemnaURL:    remnaURL,
	}
}
