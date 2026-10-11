package middlewares

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"runtime/debug"

	appctx "go-htmx-jobs/internal/context"
)

// Recover wraps every request: it recovers panics -> 500, and injects a
// logger carrying a random request id into the context. Handlers then do
//   lgr := context.GetLogger(req.Context(), "todos.list")
// (The starter does the same, plus puts the DB handle on the context.)
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := newRequestID()
		lgr := slog.Default().With("req_id", reqID)
		ctx := appctx.WithLogger(r.Context(), lgr)

		defer func() {
			if rec := recover(); rec != nil {
				lgr.Error("panic recovered",
					"panic", rec,
					"path", r.URL.Path,
					"stack", string(debug.Stack()),
				)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func newRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}
