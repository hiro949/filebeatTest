// Package model contains the data transfer objects (DTOs).
package model

// GreetingResponse represents the HTTP response for greeting
type GreetingResponse struct {
	Message string `json:"message"`
}

// NewGreetingResponse creates a new GreetingResponse
func NewGreetingResponse(message string) *GreetingResponse {
	return &GreetingResponse{
		Message: message,
	}
}
