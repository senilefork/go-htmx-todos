package context

import (
	"context"
	"database/sql"
	"log/slog"
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

// DBTX is the subset of database/sql that both *sql.DB (the pool) and
// *sql.Tx (a transaction) implement. Storing this interface — not a
// concrete *sql.DB — means the same GetDB call works inside and outside
// a transaction. (The starter gets the same effect from go-jet's qrm.DB.)
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

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
