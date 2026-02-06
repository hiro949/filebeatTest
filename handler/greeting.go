// Package handler contains the HTTP request handlers.
package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"filebeatTest/model"
	"filebeatTest/usecase"
)

// GreetingHandler handles HTTP requests for greeting
type GreetingHandler struct {
	greetingUseCase usecase.GreetingUseCaseInterface
	logger          *slog.Logger
}

// NewGreetingHandler creates a new GreetingHandler
func NewGreetingHandler(greetingUseCase usecase.GreetingUseCaseInterface, logger *slog.Logger) *GreetingHandler {
	return &GreetingHandler{
		greetingUseCase: greetingUseCase,
		logger:          logger,
	}
}

// Handle processes the greeting HTTP request
func (h *GreetingHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Warn("Method not allowed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	h.logger.Debug("Processing greeting request",
		slog.String("name", name))

	message := h.greetingUseCase.Greet(name)

	response := model.NewGreetingResponse(message)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Failed to encode response",
			slog.String("error", err.Error()))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Greeting response sent successfully",
		slog.String("name", name))
}
