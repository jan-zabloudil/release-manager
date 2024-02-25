package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jan-zabloudil/release-manager/transport"
	envx "go.strv.io/env"
)

type supabaseConfig struct {
	ApiURL    string `env:"SUPABASE_API_URL"`
	SecretKey string `env:"SUPABASE_SECRET_KEY"`
}

type serviceConfig struct {
	Port     uint           `env:"PORT"`
	LogLevel slog.Level     `env:"LOG_LEVEL"`
	Supabase supabaseConfig `env:",dive"`
}

type server struct {
	server *http.Server
	wg     sync.WaitGroup
}

func main() {
	cfg := serviceConfig{
		Port:     8080,
		LogLevel: slog.LevelDebug,
	}

	err := os.Setenv("APP_PREFIX", "RELEASE_MANAGER")
	if err != nil {
		panic(err)
	}

	envx.MustApply(&cfg)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	h := transport.NewHandler()

	srv := &server{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      h.Mux,
			IdleTimeout:  time.Minute,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	err = srv.run()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func (srv *server) run() error {
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		slog.Info("caught signal", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := srv.server.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		slog.Info("completing background tasks", "addr", srv.server.Addr)

		srv.wg.Wait()
		shutdownError <- nil
	}()

	slog.Info("starting srv", "addr", srv.server.Addr)

	err := srv.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	slog.Info("stopped srv", "addr", srv.server.Addr)

	return nil
}
