package middlewares

import (
	"net/http"

	appctx "go-htmx-jobs/internal/context"
)

// ApplyPrivileges is a minimal stub matching starter signature.
func ApplyPrivileges(privIDs []int64, appCtx *appctx.AppContext) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}
