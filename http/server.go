package http

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Config struct {
	Host string
	Port string
}

func NewServer(
	config *Config,
	logger *slog.Logger, // Ensure logger is passed
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		logger,
		*config, // Pass the config as a value
	)
	var handler http.Handler = mux
	// Add any middleware here
	return handler
}

func setupServer(ctx context.Context, config *Config) error {
	srv := NewServer(config, slog.Default()) // Use the default logger or pass one
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
	}

	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}

func runServer(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Mock Config (replace with actual config)
	config := &Config{
		Host: "0.0.0.0",
		Port: "8000",
	}

	// Start the server
	return setupServer(ctx, config)
}

func BasemapApi() {
	ctx := context.Background()

	// Environment variables function
	getenv := func(key string) string {
		switch key {
		case "BASEMAP_DOMAIN":
			return "localhost"
		case "BASEMAP_LOG_LEVEL":
			return "DEBUG"
		default:
			return ""
		}
	}

	if err := runServer(ctx, os.Args, getenv, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
