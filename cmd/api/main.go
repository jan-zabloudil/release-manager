package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"release-manager/repository"
	"release-manager/service"
	"release-manager/transport"
	"release-manager/transport/utils"

	"github.com/nedpals/supabase-go"
	httpx "go.strv.io/net/http"
	timex "go.strv.io/time"
)

func main() {
	ctx := context.Background()
	cfg := loadConfig(ctx)

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	supaClient := supabase.CreateClient(cfg.Supabase.ApiURL, cfg.Supabase.SecretKey)

	repo := repository.NewRepository(supaClient)
	svc := service.NewService(repo.User)
	h := transport.NewHandler(svc.User)

	serverConfig := httpx.ServerConfig{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: h.Mux,
		Limits: &httpx.Limits{
			Timeouts: &httpx.Timeouts{
				IdleTimeout:  timex.Duration(time.Minute),
				ReadTimeout:  timex.Duration(5 * time.Second),
				WriteTimeout: timex.Duration(10 * time.Second),
			},
		},
		Logger: utils.NewServerLogger("server"),
	}

	server := httpx.NewServer(&serverConfig)
	slog.Info("starting server", "addr", serverConfig.Addr)
	if err := server.Run(ctx); err != nil {
		slog.Error("running server", "error", err)
		os.Exit(1)
	}
}
