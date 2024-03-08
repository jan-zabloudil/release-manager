package main

import (
	"log/slog"
	"os"

	"github.com/jan-zabloudil/release-manager/utils"
	"github.com/joho/godotenv"
	envx "go.strv.io/env"
)

type supabaseConfig struct {
	ApiURL    string `env:"SUPABASE_API_URL" validate:"required"`
	SecretKey string `env:"SUPABASE_SECRET_KEY" validate:"required"`
}

type serviceConfig struct {
	Port     uint           `env:"PORT"`
	LogLevel slog.Level     `env:"LOG_LEVEL"`
	Supabase supabaseConfig `env:",dive"`
}

func loadConfig() serviceConfig {
	cfg := serviceConfig{
		Port:     8080,
		LogLevel: slog.LevelInfo,
	}

	if err := os.Setenv("APP_PREFIX", "RELEASE_MANAGER"); err != nil {
		panic(err)
	}

	_ = godotenv.Load()
	envx.MustApply(&cfg)
	if err := utils.Validate.Struct(&cfg); err != nil {
		panic(err)
	}

	return cfg
}
