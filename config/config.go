package config

import (
	"context"
	"log/slog"
	"time"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type ResendConfig struct {
	APIKey               string `env:"API_KEY, required"`
	Sender               string `env:"SENDER, required"`
	TestRecipient        string `env:"TEST_RECIPIENT, required"`
	SendToRealRecipients bool   `env:"SEND_TO_REAL_RECIPIENTS, required"`
}

type ServerConfig struct {
	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT, default=60s"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT, default=5s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT, default=10s"`
}

type SupabaseConfig struct {
	DatabaseURL  string `env:"DATABASE_URL, required"`
	APIURL       string `env:"API_URL, required"`
	APISecretKey string `env:"API_SECRET_KEY, required"`
}

type ServiceConfig struct {
	Port     uint           `env:"PORT, default=8080"`
	LogLevel slog.Level     `env:"LOG_LEVEL, default=INFO"`
	Supabase SupabaseConfig `env:", prefix=SUPABASE_"`
	Server   ServerConfig   `env:", prefix=SERVER_"`
	Resend   ResendConfig   `env:", prefix=RESEND_"`
}

func Load(ctx context.Context) ServiceConfig {
	var cfg ServiceConfig

	_ = godotenv.Load()
	if err := envconfig.Process(ctx, &cfg); err != nil {
		panic(err)
	}

	return cfg
}
