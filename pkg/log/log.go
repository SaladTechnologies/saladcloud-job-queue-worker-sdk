package log

import (
	"context"
	"log/slog"
)

type loggerKey struct{}

// Adds a `slog.Logger` to the context used to log the SaladCloud Job Queue
// worker activity.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// Gets the `slog.Logger` from the context used to log the SaladCloud Job Queue
// worker activity.
func FromContext(ctx context.Context) *slog.Logger {
	value := ctx.Value(loggerKey{})
	if value == nil {
		return slog.New(nopHandler{})
	}

	logger, ok := value.(*slog.Logger)
	if !ok {
		return slog.New(nopHandler{})
	}

	return logger
}

// A no-op `slog.Handler`.
type nopHandler struct{}

func (nopHandler) Enabled(context.Context, slog.Level) bool {
	return false
}

func (nopHandler) Handle(context.Context, slog.Record) error {
	return nil
}

func (h nopHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h nopHandler) WithGroup(string) slog.Handler {
	return h
}

// A leveled `slog.Handler`.
type LeveledHandler struct {
	handler slog.Handler
	leveler slog.Leveler
}

// Creates a new `slog.Handler` with a configurable level.
func NewLeveledHandler(handler slog.Handler, leveler slog.Leveler) slog.Handler {
	return &LeveledHandler{
		handler: handler,
		leveler: leveler,
	}
}

func (h *LeveledHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.leveler.Level()
}

func (h *LeveledHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.handler.Handle(ctx, r)
}

func (h *LeveledHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewLeveledHandler(h.handler.WithAttrs(attrs), h.leveler)
}

func (h *LeveledHandler) WithGroup(name string) slog.Handler {
	return NewLeveledHandler(h.handler.WithGroup(name), h.leveler)
}
