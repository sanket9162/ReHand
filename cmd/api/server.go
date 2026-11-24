package main

import (
	"context"
	"fmt"
	"log"
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
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	// define gracefull shutdown
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		s := <-quit
		log.Println("received signal", s)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	log.Println("start serving on", srv.Addr)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed: %w", err)
	}

	err := <-shutdownError
	if err != nil {
		return err
	}

	return nil
}
