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

// ClientServiceConfig contains client service deployment URL and client service routes
// It is used for generating links to the client service (when sending a project invitation email)
type ClientServiceConfig struct {
	URL                   string `env:"URL, required"`
	SignUpRoute           string `env:"SIGN_UP_ROUTE, required"`
	AcceptInvitationRoute string `env:"ACCEPT_INVITATION_ROUTE, required"`
	RejectInvitationRoute string `env:"REJECT_INVITATION_ROUTE, required"`
}

type ServiceConfig struct {
	Port          uint                `env:"APP_PORT, default=8080"`
	LogLevel      slog.Level          `env:"LOG_LEVEL, default=INFO"`
	Supabase      SupabaseConfig      `env:", prefix=SUPABASE_"`
	Server        ServerConfig        `env:", prefix=SERVER_"`
	Resend        ResendConfig        `env:", prefix=RESEND_"`
	ClientService ClientServiceConfig `env:", prefix=CLIENT_SERVICE_"`
}

func Load(ctx context.Context) ServiceConfig {
	var cfg ServiceConfig

	_ = godotenv.Load()
	if err := envconfig.Process(ctx, &cfg); err != nil {
		panic(err)
	}

	return cfg
}
