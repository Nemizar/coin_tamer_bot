package sl

import (
	"context"
	"log/slog"

	"github.com/Nemizar/coin_tamer_bot/internal/core/ports"
)

var _ ports.Logger = &Logger{}

type Logger struct {
	l *slog.Logger
}

func NewSlogLogger(l *slog.Logger) *Logger {
	return &Logger{l: l}
}

func (s *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	s.l.InfoContext(ctx, msg, args...)
}

func (s *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	s.l.ErrorContext(ctx, msg, args...)
}

func (s *Logger) Info(msg string, args ...any) {
	s.l.Info(msg, args...)
}

func (s *Logger) Error(msg string, args ...any) {
	s.l.Error(msg, args...)
}
