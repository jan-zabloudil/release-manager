package main

import (
	"context"
	"log/slog"

	"release-manager/mailer"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type supabaseConfig struct {
	ApiURL    string `env:"API_URL, required"`
	SecretKey string `env:"SECRET_KEY, required"`
}

type mailerConfig struct {
	ResendApiKey  string `env:"RESEND_API_KEY, required"`
	Sender        string `env:"SENDER, default=onboarding@resend.dev"`
	TestingMode   bool   `env:"TESTING_MODE, default=true"`
	TestRecipient string `env:"TEST_RECIPIENT, default=delivered@resend.dev"`
}

type serviceConfig struct {
	Port     uint           `env:"PORT, default=8080"`
	LogLevel slog.Level     `env:"LOG_LEVEL, default=INFO"`
	Supabase supabaseConfig `env:", prefix=SUPABASE_"`
	Mailer   mailerConfig   `env:", prefix=MAILER_"`
}

func loadConfig(ctx context.Context) serviceConfig {
	var cfg serviceConfig

	_ = godotenv.Load()
	if err := envconfig.Process(ctx, &cfg); err != nil {
		panic(err)
	}

	return cfg
}

func (c *serviceConfig) ToMailerConfig() mailer.Config {
	return mailer.Config{
		TestingMode:   c.Mailer.TestingMode,
		ApiKey:        c.Mailer.ResendApiKey,
		TestRecipient: c.Mailer.TestRecipient,
		Sender:        c.Mailer.Sender,
	}
}
