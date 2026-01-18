package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string     `yaml:"env" env-deafault:"local"`
	DSN  string     `yaml:"dsn" env-required:"true"`
	GRPC GRPCConfig `yaml:"grpc" env-required:"true"`
}

type GRPCConfig struct {
	Host    string        `yaml:"host"`
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad(configPath string) *Config {
	if configPath == "" {
		panic("config path required")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("file does not exist: " + configPath)
	}
	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		panic(err)
	}
	return &config
}
