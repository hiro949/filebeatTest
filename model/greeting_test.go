package model

import "testing"

func TestNewGreetingResponse(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "create greeting response",
			message:  "Hello, World!",
			expected: "Hello, World!",
		},
		{
			name:     "create empty greeting response",
			message:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := NewGreetingResponse(tt.message)
			if response.Message != tt.expected {
				t.Errorf("expected message %q, got %q", tt.expected, response.Message)
			}
		})
	}
}
