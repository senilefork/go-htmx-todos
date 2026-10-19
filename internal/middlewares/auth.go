package middlewares

import (
	"net/http"

	appctx "go-htmx-jobs/internal/context"
)

// Auth is a minimal stub matching starter signature.
func Auth(appCtx *appctx.AppContext, enforce bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}
