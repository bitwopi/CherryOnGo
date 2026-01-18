package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env             string        `yaml:"env" env-deafault:"local"`
	DSN             string        `yaml:"dsn" env-required:"true"`
	GRPC            GRPCConfig    `yaml:"grpc" env-required:"true"`
	Redis           RedisConfig   `yaml:"redis" env-required:"true"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl" env-deafault:"15m"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env-deafault:"30d"`
	JWTSecret       string        `yaml:"jwt_secret" env-required:"true"`
}

type GRPCConfig struct {
	Host    string        `yaml:"host"`
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type RedisConfig struct {
	Addr  string `yaml:"addr"`
	DBNum int    `yaml:"dbnum"`
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
