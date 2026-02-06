package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_usecase "filebeatTest/mock/usecase"
	"filebeatTest/model"
	"filebeatTest/pkg/logger"

	"github.com/golang/mock/gomock"
)

func TestGreetingHandler_Handle(t *testing.T) {
	tests := []struct {
		name           string
		queryParam     string
		mockReturn     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful greeting with name",
			queryParam:     "name=Alice",
			mockReturn:     "Good morning, Alice! Welcome to our API.",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Good morning, Alice! Welcome to our API."}`,
		},
		{
			name:           "successful greeting without name",
			queryParam:     "",
			mockReturn:     "Good morning, Guest! Welcome to our API.",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Good morning, Guest! Welcome to our API."}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mock_usecase.NewMockGreetingUseCaseInterface(ctrl)

			req := httptest.NewRequest(http.MethodGet, "/greet?"+tt.queryParam, http.NoBody)
			rec := httptest.NewRecorder()

			if tt.queryParam != "" {
				mockUseCase.EXPECT().Greet("Alice").Return(tt.mockReturn)
			} else {
				mockUseCase.EXPECT().Greet("").Return(tt.mockReturn)
			}

			log := logger.NewWithLevel(slog.LevelError) // Use error level to suppress logs in tests
			handler := NewGreetingHandler(mockUseCase, log)
			handler.Handle(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			var response model.GreetingResponse
			if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			var expectedResponse model.GreetingResponse
			if err := json.Unmarshal([]byte(tt.expectedBody), &expectedResponse); err != nil {
				t.Fatalf("failed to unmarshal expected response: %v", err)
			}

			if response.Message != expectedResponse.Message {
				t.Errorf("expected message %q, got %q", expectedResponse.Message, response.Message)
			}
		})
	}
}

func TestGreetingHandler_Handle_MethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock_usecase.NewMockGreetingUseCaseInterface(ctrl)

	req := httptest.NewRequest(http.MethodPost, "/greet", http.NoBody)
	rec := httptest.NewRecorder()

	log := logger.NewWithLevel(slog.LevelError)
	handler := NewGreetingHandler(mockUseCase, log)
	handler.Handle(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}
