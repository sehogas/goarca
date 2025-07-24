package util

import (
	"context"
	"log/slog"
)

// SlogWriter is an io.Writer that directs output to a slog.Logger.
type SlogWriter struct {
	logger *slog.Logger
	level  slog.Level
}

// NewSlogWriter creates a new SlogWriter.
func NewSlogWriter(logger *slog.Logger, level slog.Level) *SlogWriter {
	return &SlogWriter{
		logger: logger,
		level:  level,
	}
}

// Write implements the io.Writer interface.
func (sw *SlogWriter) Write(p []byte) (n int, err error) {
	sw.logger.Log(context.Background(), sw.level, string(p)) // Log the message at the specified level
	return len(p), nil
}
