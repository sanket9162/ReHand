package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	readTimeout       = 10 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 1 * time.Minute
	readHeaderTimeout = 10 * time.Second
)

func (app *application) server() error {

	// define the server with timeout
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", app.config.Port),
		Handler:           app.router(),
		ErrorLog:          slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	// define graceful shutdown
	shutdownError := make(chan error, 1)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		app.logger.Info("received signal", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	app.logger.Info("start serving", "addr", srv.Addr)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed: %w", err)
	}

	err := <-shutdownError
	if err != nil {
		return err
	}

	return nil
}
