// Package main is the entry point of the greeting API application.
package main

import (
	"log/slog"
	"net/http"
	"time"

	"filebeatTest/handler"
	"filebeatTest/middleware"
	"filebeatTest/pkg/logger"
	"filebeatTest/usecase"
)

const (
	// Application metadata
	appVersion    = "1.0.0"
	appEnvironment = "development"

	// Server configuration
	serverPort     = ":8080"
	greetEndpoint  = "/greet"

	// Server timeouts
	readHeaderTimeout = 10 * time.Second
	readTimeout       = 30 * time.Second
	writeTimeout      = 30 * time.Second
	idleTimeout       = 60 * time.Second
)

func main() {
	// Initialize structured logger
	log := logger.New()

	log.Info("Starting greeting API server",
		slog.String("version", appVersion),
		slog.String("environment", appEnvironment))

	// Dependency injection
	greetingUseCase := usecase.NewGreetingUseCase(log)
	greetingHandler := handler.NewGreetingHandler(greetingUseCase, log)

	// Create mux for middleware support
	mux := http.NewServeMux()
	mux.HandleFunc(greetEndpoint, greetingHandler.Handle)

	// Apply logging middleware
	httpHandler := middleware.LoggingMiddleware(log)(mux)

	log.Info("Server configured",
		slog.String("port", serverPort),
		slog.String("endpoint", greetEndpoint))

	server := &http.Server{
		Addr:              serverPort,
		Handler:           httpHandler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	log.Info("Server starting", slog.String("addr", server.Addr))

	if err := server.ListenAndServe(); err != nil {
		log.Error("Server failed to start",
			slog.String("error", err.Error()))
	}
}
