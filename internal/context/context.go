package context

import (
	"context"
	"log/slog"
)

// Request-scoped values ride on the stdlib context.Context — the same
// mechanism the starter uses for db, logger, user_id, and tx.

type ctxKey int

const loggerKey ctxKey = 0

// WithLogger attaches a logger pre-loaded with the request id.
func WithLogger(ctx context.Context, lgr *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, lgr)
}

// GetLogger returns the request logger, or a default one if the
// Recover middleware hasn't run (e.g. in tests).
func GetLogger(ctx context.Context, scope string) *slog.Logger {
	if lgr, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return lgr.With("scope", scope)
	}
	return slog.Default().With("scope", scope)
}
