package context

import (
	"context"
	"log/slog"

	"github.com/go-jet/jet/v2/qrm"
)

// Request-scoped values ride on the stdlib context.Context — the same
// mechanism the starter uses for db, logger, user_id, and tx.

type ctxKey int

//prevents key collisions
const (
	loggerKey ctxKey = iota
	dbKey
)

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

// DBTX is go-jet's database interface. Both *sql.DB (the pool) and
// *sql.Tx (a transaction) implement it, so GetDB works inside and
// outside a transaction — same reason the starter returns qrm.DB.
type DBTX = qrm.DB

// WithDB attaches the database handle to the context.
func WithDB(ctx context.Context, db DBTX) context.Context {
	return context.WithValue(ctx, dbKey, db)
}

// GetDB returns the database handle that cmd/web.go attached at startup.
// A missing handle panics on purpose: it can only happen if bootstrap
// wiring is broken, which is a programmer error, not a runtime condition.
func GetDB(ctx context.Context) DBTX {
	return ctx.Value(dbKey).(DBTX)
}
