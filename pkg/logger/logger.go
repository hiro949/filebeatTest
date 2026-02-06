// Package logger provides structured logging for the application.
package logger

import (
	"log/slog"
	"os"
)

// New creates a new JSON logger for structured logging.
// This logger outputs logs in JSON format which can be easily consumed by Filebeat.
func New() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			// Rename timestamp field to @timestamp for Elasticsearch/Kibana compatibility
			if a.Key == slog.TimeKey {
				a.Key = "@timestamp"
			}
			return a
		},
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}

// NewWithLevel creates a new JSON logger with specified log level.
func NewWithLevel(level slog.Level) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Key = "@timestamp"
			}
			return a
		},
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}
