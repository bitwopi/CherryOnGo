package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RESTBindUrl string
	GRPCBindUrl string
	RemnaAPIKey string
	RemnaURL    string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	rBindUrl := os.Getenv("REMNAVAWE_REST_BIND_URL")
	if len(rBindUrl) == 0 {
		rBindUrl = ":8080"
	}
	gBindUrl := os.Getenv("REMNAVAWE_GRPC_BIND_URL")
	if len(gBindUrl) == 0 {
		gBindUrl = ":8080"
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
		RESTBindUrl: rBindUrl,
		GRPCBindUrl: gBindUrl,
		RemnaAPIKey: remnaAPIKey,
		RemnaURL:    remnaURL,
	}
}
