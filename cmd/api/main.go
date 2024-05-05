package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"release-manager/config"
	githubx "release-manager/github"
	"release-manager/repository"
	resendx "release-manager/resend"
	"release-manager/service"
	"release-manager/transport/handler"

	"github.com/nedpals/supabase-go"
	"go.strv.io/background"
	"go.strv.io/background/observer"
	httpx "go.strv.io/net/http"
	timex "go.strv.io/time"
)

func main() {
	ctx := context.Background()
	cfg := config.Load(ctx)

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	taskManager := background.NewManagerWithOptions(background.Options{
		Observer: observer.Slog{},
	})

	supaClient := supabase.CreateClient(cfg.Supabase.APIURL, cfg.Supabase.SecretKey)
	githubClient := githubx.NewClient()
	resendClient := resendx.NewClient(taskManager, cfg.Resend)

	repo := repository.NewRepository(supaClient)
	svc := service.NewService(
		repo.Auth,
		repo.User,
		repo.Project,
		repo.Settings,
		repo.Release,
		githubClient,
		resendClient,
	)
	h := handler.NewHandler(svc.Auth, svc.User, svc.Project, svc.Settings, svc.Release)

	serverConfig := httpx.ServerConfig{
		Addr: fmt.Sprintf(":%d", cfg.Port),
		Hooks: httpx.ServerHooks{
			BeforeShutdown: []httpx.ServerHookFunc{
				func(_ context.Context) {
					slog.Info("waiting for tasks to finish")
					taskManager.Close()
				},
			},
		},
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
