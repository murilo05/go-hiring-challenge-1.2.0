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

	"github.com/mytheresa/go-hiring-challenge/app/api/handler/http/catalog"
	"github.com/mytheresa/go-hiring-challenge/app/api/handler/http/middleware"
	"github.com/mytheresa/go-hiring-challenge/app/database"

	"github.com/mytheresa/go-hiring-challenge/app/repository"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Starting Logs
	var lg *zap.Logger
	switch os.Getenv("ENVIRONMENT") {
	case "DEV":
		lg, _ = zap.NewDevelopment()
	default:
		lg, _ = zap.NewProduction()
	}

	defer lg.Sync()
	logger := lg.Sugar()

	logger.Info("Starting the application: ", os.Getenv("APP"), "-", os.Getenv("ENVIRONMENT"))

	// signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize database connection
	db, close := database.New(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		logger,
	)
	defer close()

	// Initialize handlers
	productsRepository := repository.NewProductsRepository(db, logger)
	catalogHandler := catalog.NewCatalogHandler(productsRepository, logger)

	// Set up routing
	mux := http.NewServeMux()
	mux.Handle(
		"/catalog",
		middleware.LoggingMiddleware(logger)(
			http.HandlerFunc(catalogHandler.ListProducts),
		),
	)
	mux.Handle(
		"/catalog/{code}",
		middleware.LoggingMiddleware(logger)(
			http.HandlerFunc(catalogHandler.GetProduct),
		),
	)

	// Set up the HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", os.Getenv("HTTP_PORT")),
		Handler: mux,
	}

	// Start the server
	go func() {
		logger.Info("Starting server on http://", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed: %s", err)
		}
		logger.Info("Server stopped gracefully")
	}()

	<-ctx.Done()
	logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	srv.Shutdown(shutdownCtx)
}
