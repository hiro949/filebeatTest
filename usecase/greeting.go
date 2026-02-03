// Package usecase contains the application business logic and use cases.
package usecase

import (
	"log/slog"

	"filebeatTest/domain"
)

//go:generate /home/user/go/bin/mockgen -source=greeting.go -destination=../mock/usecase/greeting_mock.go -package=mock_usecase

// GreetingUseCaseInterface defines the interface for greeting use case
type GreetingUseCaseInterface interface {
	Greet(name string) string
}

// GreetingUseCase handles greeting business operations
type GreetingUseCase struct {
	logger *slog.Logger
}

// NewGreetingUseCase creates a new GreetingUseCase
func NewGreetingUseCase(logger *slog.Logger) *GreetingUseCase {
	return &GreetingUseCase{
		logger: logger,
	}
}

// Greet generates a greeting message for the given name
func (uc *GreetingUseCase) Greet(name string) string {
	uc.logger.Debug("Generating greeting",
		slog.String("name", name))

	greeting := domain.NewGreeting(name)
	message := greeting.GenerateMessage()

	uc.logger.Info("Greeting generated",
		slog.String("name", name),
		slog.String("recipient", greeting.Recipient))

	return message
}
