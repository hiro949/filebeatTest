package domain

import (
	"testing"
	"time"
)

func TestNewGreeting(t *testing.T) {
	tests := []struct {
		name      string
		recipient string
		want      string
	}{
		{
			name:      "with name",
			recipient: "Alice",
			want:      "Alice",
		},
		{
			name:      "empty name defaults to Guest",
			recipient: "",
			want:      "Guest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			greeting := NewGreeting(tt.recipient)
			if greeting.Recipient != tt.want {
				t.Errorf("NewGreeting() recipient = %v, want %v", greeting.Recipient, tt.want)
			}
		})
	}
}

func TestGreeting_GenerateMessageWithTime(t *testing.T) {
	tests := []struct {
		name      string
		recipient string
		hour      int
		want      string
	}{
		{
			name:      "good morning for Alice at 8am",
			recipient: "Alice",
			hour:      8,
			want:      "Good morning, Alice! Welcome to our API.",
		},
		{
			name:      "good afternoon for Bob at 3pm",
			recipient: "Bob",
			hour:      15,
			want:      "Good afternoon, Bob! Welcome to our API.",
		},
		{
			name:      "good evening for Charlie at 7pm",
			recipient: "Charlie",
			hour:      19,
			want:      "Good evening, Charlie! Welcome to our API.",
		},
		{
			name:      "good night for Guest at 11pm",
			recipient: "Guest",
			hour:      23,
			want:      "Good night, Guest! Welcome to our API.",
		},
		{
			name:      "good night for Guest at 2am",
			recipient: "Guest",
			hour:      2,
			want:      "Good night, Guest! Welcome to our API.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			greeting := &Greeting{Recipient: tt.recipient}
			testTime := time.Date(2024, 1, 1, tt.hour, 0, 0, 0, time.UTC)
			got := greeting.GenerateMessageWithTime(testTime)
			if got != tt.want {
				t.Errorf("Greeting.GenerateMessageWithTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
