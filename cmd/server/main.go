package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HV-Hung/family-svc/internal/config"
	"github.com/HV-Hung/family-svc/internal/database"
	"github.com/HV-Hung/family-svc/internal/handler"
	"github.com/HV-Hung/family-svc/internal/middleware"
	"github.com/HV-Hung/family-svc/internal/telemetry"
)

func main() {
	// Structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// Load configuration from environment
	cfg := config.Load()
	slog.Info("configuration loaded",
		"http_port", cfg.HTTPPort,
		"db_host", cfg.DBHost,
		"db_port", cfg.DBPort,
		"db_name", cfg.DBName,
	)

	// Connect to database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := database.NewPool(ctx, cfg.DSN())
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	slog.Info("database connected")

	// Initialise Prometheus metrics registry
	reg := telemetry.NewRegistry()

	// Register routes
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/hello", handler.HelloHandler())
	mux.HandleFunc("GET /healthz/live", handler.LivenessHandler())
	mux.HandleFunc("GET /healthz/ready", handler.ReadinessHandler(pool))
	mux.HandleFunc("GET /metrics", handler.MetricsHandler(reg))

	// Create server.
	// Middleware chain (outer → inner): LogRequest → InstrumentHandler → mux
	// Health-check probes (/healthz/*) are excluded from both logging and metrics.
	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      middleware.LogRequest(middleware.InstrumentHandler(reg, mux)),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		slog.Info("server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown on SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	slog.Info("shutting down", "signal", sig.String())

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}
