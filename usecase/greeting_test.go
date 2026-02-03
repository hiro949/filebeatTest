package usecase

import (
	"log/slog"
	"strings"
	"testing"

	"filebeatTest/pkg/logger"
)

var validGreetings = []string{
	"Good morning",
	"Good afternoon",
	"Good evening",
	"Good night",
}

func TestGreetingUseCase_Greet(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedName string
	}{
		{
			name:         "greet with name",
			input:        "Bob",
			expectedName: "Bob",
		},
		{
			name:         "greet with empty name",
			input:        "",
			expectedName: "Guest",
		},
		{
			name:         "greet with Japanese name",
			input:        "田中",
			expectedName: "田中",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logger.NewWithLevel(slog.LevelError)
			uc := NewGreetingUseCase(log)
			got := uc.Greet(tt.input)

			assertContains(t, got, tt.expectedName)
			assertHasValidGreeting(t, got)
			assertContains(t, got, "Welcome to our API.")
		})
	}
}

func assertContains(t *testing.T, got, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Errorf("got %q, want message containing %q", got, want)
	}
}

func assertHasValidGreeting(t *testing.T, got string) {
	t.Helper()
	for _, greeting := range validGreetings {
		if strings.HasPrefix(got, greeting) {
			return
		}
	}
	t.Errorf("got %q, want message starting with one of %v", got, validGreetings)
}
