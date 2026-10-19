package builders

import (
	gctx "context"
	"net/http"

	"go-htmx-jobs/internal/constant"

	"github.com/go-chi/chi/v5"
)

type privateHandlerBuilder struct {
	appCtx     interface{}
	router     chi.Router
	mw         []func(next http.Handler) http.Handler
	method     string
	pattern    string
	handle     http.HandlerFunc
	privileges []constant.PrivilegeName // privileges required to access private endpoint
}

func (mb *privateHandlerBuilder) Register(method string, pattern string, handle http.HandlerFunc) *privateHandlerBuilder {
	mb.method = method
	mb.pattern = pattern
	mb.handle = handle
	return mb
}

// SetRequiredPrivileges lists privilege names required to access this endpoint.
func (mb *privateHandlerBuilder) SetRequiredPrivileges(privileges ...constant.PrivilegeName) *privateHandlerBuilder {
	mb.privileges = privileges
	return mb
}

// SetMiddlewares attaches per-handler middlewares.
func (mb *privateHandlerBuilder) SetMiddlewares(middlewares ...func(next http.Handler) http.Handler) *privateHandlerBuilder {
	mb.mw = middlewares
	return mb
}

// Build registers the handler on the scoped router.
// No-DB adaptation: skips privilege DAO upserts and just applies privilege middleware if present.
func (mb *privateHandlerBuilder) Build(ctx gctx.Context) {
	r := mb.router

	if len(mb.privileges) > 0 {
		r = mb.router.With(ApplyPrivilegesStub())
	}

	if len(mb.mw) > 0 {
		r = r.With(mb.mw...)
	}

	r.MethodFunc(mb.method, mb.pattern, mb.handle)
}

// ApplyPrivilegesStub is a no-op passthrough for privilege middleware.
func ApplyPrivilegesStub() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}
