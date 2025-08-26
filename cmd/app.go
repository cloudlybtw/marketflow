package cmd

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"marketflow/internal/adapter/config"
	handler "marketflow/internal/adapter/http"
	"marketflow/internal/adapter/logger"
	"marketflow/internal/adapter/source"
	"marketflow/internal/adapter/storage"
	"marketflow/internal/adapter/storage/cache"
	"marketflow/internal/adapter/storage/db"
	"marketflow/internal/core/service"
)

var PortN int

func NewApp() {
	flag.IntVar(&PortN, "port", 8080, "Port number")
	flag.Usage = func() {
		fmt.Print("Marketflow\n\n",
			"Usage:\n",
			"\tmarketflow [-port <N>]\n\tmarketflow --help\n\n",
			"Options:\n",
			"- --port N\tPort number\n")
	}
	flag.Parse()

	config := config.New()

	logger.Set(config.App)
	slog.Info("Staring application", "app", config.App.Name, "on port", PortN)

	redis, err := cache.NewRedisClient(config.Redis)
	if err != nil {
		slog.Error("Failed to connect to the redis", "error", err)
		os.Exit(1)
	}

	database, err := db.OpenDB(config.DB)
	if err != nil {
		slog.Error("Failed to connect to the database", "error", err)
		os.Exit(1)
	}
	defer database.Close()
	pgs := db.NewPostgresClient(database)

	priceRep := storage.NewPriceRepository(pgs, redis)
	priceService := service.NewPriceService(priceRep)
	m := source.NewModeAdapter(config.Exchanges, priceService)

	// Context with cancellation on SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start ModeAdapter with cancelable context
	go m.Start(ctx)

	// Setup HTTP server with shutdown support
	req := handler.NewHTTPHandler(priceService, m)
	mux := handler.NewRouter(req)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", PortN),
		Handler: mux,
	}

	go func() {
		slog.Info("HTTP server starting", "port", PortN)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server failed", "error", err)
			stop() // trigger shutdown
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	slog.Info("Shutting down gracefully...")

	// Gracefully shutdown HTTP server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server shutdown error", "error", err)
	} else {
		slog.Info("HTTP server stopped cleanly")
	}

	slog.Info("Application shutdown complete")
}
