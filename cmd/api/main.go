package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"release-manager/repository"
	"release-manager/service"
	"release-manager/transport"

	"github.com/nedpals/supabase-go"
	httpx "go.strv.io/net/http"
	timex "go.strv.io/time"
)

func main() {
	ctx := context.Background()
	cfg := loadConfig(ctx)

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	supaClient := supabase.CreateClient(cfg.Supabase.APIURL, cfg.Supabase.SecretKey)

	repo := repository.NewRepository(supaClient)
	svc := service.NewService(repo.Auth, repo.User)
	h := transport.NewHandler(svc.Auth, svc.User)

	serverConfig := httpx.ServerConfig{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: h.Mux,
		Limits: &httpx.Limits{
			Timeouts: &httpx.Timeouts{
				IdleTimeout:  timex.Duration(cfg.Server.IdleTimeout),
				ReadTimeout:  timex.Duration(cfg.Server.ReadTimeout),
				WriteTimeout: timex.Duration(cfg.Server.WriteTimeout),
			},
		},
		Logger: logger.WithGroup("server"),
	}

	server := httpx.NewServer(&serverConfig)
	slog.Info("starting server", "addr", serverConfig.Addr)
	if err := server.Run(ctx); err != nil {
		slog.Error("running server", "error", err)
		os.Exit(1)
	}
}
