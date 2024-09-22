package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	server *http.Server
	logger *slog.Logger
}

func NewApp(server *http.Server, logger *slog.Logger) *App {
	return &App{
		server: server,
		logger: logger,
	}
}

func (a *App) Run() {
	ctx := context.Background()

	errChan := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				errChan <- fmt.Errorf("server listen and serve, %w", err)
			}
		}
	}()

	select {
	case err := <-errChan:
		a.logger.Error("component failed", slog.String("error", err.Error()))
	case sig := <-sigChan:
		a.logger.Info("received stopped signal", slog.String("signal", sig.String()))
	}

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("server shutdown failed", slog.String("error", err.Error()))
	}
}
