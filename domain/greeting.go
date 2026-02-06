// Package domain contains the domain models and business logic.
package domain

import "time"

const (
	// Time boundaries for different greetings (in hours)
	morningStartHour   = 5
	afternoonStartHour = 12
	eveningStartHour   = 18
	nightStartHour     = 22

	// Default recipient name
	defaultRecipient = "Guest"

	// Greeting messages
	greetingMorning   = "Good morning"
	greetingAfternoon = "Good afternoon"
	greetingEvening   = "Good evening"
	greetingNight     = "Good night"

	// Message template suffix
	messageTemplate = "! Welcome to our API."
)

// Greeting represents the domain model for greeting
type Greeting struct {
	Recipient string
}

// NewGreeting creates a new Greeting
func NewGreeting(recipient string) *Greeting {
	if recipient == "" {
		recipient = defaultRecipient
	}
	return &Greeting{
		Recipient: recipient,
	}
}

// getTimeBasedGreeting returns appropriate greeting based on current time
func getTimeBasedGreeting(t time.Time) string {
	hour := t.Hour()

	switch {
	case hour >= morningStartHour && hour < afternoonStartHour:
		return greetingMorning
	case hour >= afternoonStartHour && hour < eveningStartHour:
		return greetingAfternoon
	case hour >= eveningStartHour && hour < nightStartHour:
		return greetingEvening
	default:
		return greetingNight
	}
}

// GenerateMessage generates a greeting message based on time of day
func (g *Greeting) GenerateMessage() string {
	return g.GenerateMessageWithTime(time.Now())
}

// GenerateMessageWithTime generates a greeting message based on specified time
func (g *Greeting) GenerateMessageWithTime(t time.Time) string {
	greeting := getTimeBasedGreeting(t)
	return greeting + ", " + g.Recipient + messageTemplate
}
