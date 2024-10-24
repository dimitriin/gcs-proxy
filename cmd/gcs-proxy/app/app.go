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
	"time"

	"github.com/dimitriin/gcs-proxy/cmd/gcs-proxy/config"
)

type App struct {
	server         *http.Server
	logger         *slog.Logger
	shutdownConfig config.ShutdownConfig
}

func NewApp(
	server *http.Server,
	logger *slog.Logger,
	shutdownConfig config.ShutdownConfig,
) *App {
	return &App{
		server:         server,
		logger:         logger,
		shutdownConfig: shutdownConfig,
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

	var sig os.Signal

	select {
	case err := <-errChan:
		a.logger.Error("component failed", slog.String("error", err.Error()))
	case sig = <-sigChan:
		a.logger.Info("received stopped signal", slog.String("signal", sig.String()))

		if a.shutdownConfig.PreStopTimeout > 0 {
			a.logger.Info("sleeping before shutdown",
				slog.Duration("duration", a.shutdownConfig.PreStopTimeout),
			)

			time.Sleep(a.shutdownConfig.PreStopTimeout)
		}
	}

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("server shutdown failed", slog.String("error", err.Error()))
	}

	exitCode, err := a.shutdownConfig.ExitCodes.GetExitCode(sig)
	if err != nil {
		a.logger.Error("get exit code for signal failed",
			slog.String("error", err.Error()),
			slog.String("signal", sig.String()),
		)
	}

	a.logger.Info("exiting", slog.Int("exit_code", exitCode))

	os.Exit(exitCode)
}
