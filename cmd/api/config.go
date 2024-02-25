package main

import (
	"context"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type supabaseConfig struct {
	ApiURL    string `env:"API_URL, required"`
	SecretKey string `env:"SECRET_KEY, required"`
}

type serviceConfig struct {
	Port     uint           `env:"PORT, default=8080"`
	LogLevel slog.Level     `env:"LOG_LEVEL, default=INFO"`
	Supabase supabaseConfig `env:", prefix=SUPABASE_"`
}

func loadConfig(ctx context.Context) serviceConfig {
	var cfg serviceConfig

	_ = godotenv.Load()
	if err := envconfig.Process(ctx, &cfg); err != nil {
		panic(err)
	}

	return cfg
}
