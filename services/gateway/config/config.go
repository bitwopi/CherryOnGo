package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env             string        `yaml:"env" env-deafault:"local"`
	UserService     GRPCConfig    `yaml:"user_service" env-required:"true"`
	OrderService    GRPCConfig    `yaml:"order_service" env-required:"true"`
	RemnaService    GRPCConfig    `yaml:"remna_service" env-required:"true"`
	REST            RESTConfig    `yaml:"rest" env-required:"true"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env-deafault:"730h"`
	JWTSecret       string        `yaml:"jwt_secret" env-required:"true"`
	BotToken        string        `yaml:"bot_token" env-required:"true"`
}

type GRPCConfig struct {
	Addr       string        `yaml:"addr"`
	Timeout    time.Duration `yaml:"timeout"`
	MaxRetries uint          `yaml:"max_retries"`
}

type RESTConfig struct {
	Addr           string        `yaml:"addr"`
	RequestTimeout time.Duration `yaml:"request_timeout"`
	IdleTimeout    time.Duration `yaml:"idle_timout"`
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
