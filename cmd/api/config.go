package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type serverConfig struct {
	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT, default=60s"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT, default=5s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT, default=10s"`
}

type supabaseConfig struct {
	APIURL    string `env:"API_URL, required"`
	SecretKey string `env:"SECRET_KEY, required"`
}

type serviceConfig struct {
	Port     uint           `env:"PORT, default=8080"`
	LogLevel slog.Level     `env:"LOG_LEVEL, default=INFO"`
	Supabase supabaseConfig `env:", prefix=SUPABASE_"`
	Server   serverConfig   `env:", prefix=SERVER_"`
}

func loadConfig(ctx context.Context) serviceConfig {
	var cfg serviceConfig

	_ = godotenv.Load()
	if err := envconfig.Process(ctx, &cfg); err != nil {
		panic(err)
	}

	return cfg
}
