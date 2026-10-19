package builders

import (
	gctx "context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RoutePath is the mount prefix for a controller.
type RoutePath interface {
	Path() string
}

// PublicRoute registers public HTTP routes on the given chi router.
type PublicRoute interface {
	PublicRoute(r chi.Router)
}

// AppPublicRoute is the combined interface for public controllers.
type AppPublicRoute interface {
	RoutePath
	PublicRoute
}

type PrivateRouteBuilder struct {
	router chi.Router
	appCtx interface{}
}

func NewPrivateRouteBuilder(appCtx interface{}) PrivateRouteBuilder {
	return PrivateRouteBuilder{
		router: chi.NewRouter(),
		appCtx: appCtx,
	}
}

// PrivateRoute registers private handlers via the builder.
type PrivateRoute interface {
	PrivateRoute(ctx gctx.Context, b *PrivateRouteBuilder)
}

// AppPrivateRoute is the combined interface for private controllers.
type AppPrivateRoute interface {
	RoutePath
	PrivateRoute
}

func (rb *PrivateRouteBuilder) NewHandler() *privateHandlerBuilder {
	return &privateHandlerBuilder{
		router: rb.router,
		appCtx: rb.appCtx,
	}
}

// NewGroup creates a nested group scope on the current router.
func (rb *PrivateRouteBuilder) NewGroup(f func(g *PrivateRouteBuilder)) {
	builder := PrivateRouteBuilder{
		router: nil,
		appCtx: rb.appCtx,
	}
	rb.router.Group(func(g chi.Router) {
		builder.router = g
	})
	f(&builder)
}

// NewRoute creates a route scope at path on the current router.
func (rb *PrivateRouteBuilder) NewRoute(path string) PrivateRouteBuilder {
	builder := PrivateRouteBuilder{
		router: nil,
		appCtx: rb.appCtx,
	}
	rb.router.Route(path, func(r chi.Router) {
		builder.router = r
	})
	return builder
}

// AddMiddleware applies middleware to the current scoped router.
func (rb *PrivateRouteBuilder) AddMiddleware(middleware func(next http.Handler) http.Handler) {
	rb.router.Use(middleware)
}

// RawRouter returns the built chi router for this builder.
func (rb *PrivateRouteBuilder) RawRouter() chi.Router {
	return rb.router
}
