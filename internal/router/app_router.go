package router

import (
	gctx "context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"go-htmx-jobs/internal/context"
	"go-htmx-jobs/internal/controllers/public"
	"go-htmx-jobs/internal/middlewares"

	"github.com/go-chi/chi/v5"
)

// PublicRoute is the contract every public controller implements —
// same shape as the starter's builders.AppPublicRoute: a mount path
// plus a function that registers its sub-routes.
type PublicRoute interface {
	Path() string
	PublicRoute(chi.Router)
}

type AppRouter struct {
	public []PublicRoute
	appCtx *context.AppContext
	router chi.Router
}

func NewAppRouter(appCtx *context.AppContext) AppRouter {
	return AppRouter{
		appCtx: appCtx,
		public: []PublicRoute{
			public.NewWebPublic(appCtx), // mount /public first...
			public.NewTodos(appCtx),     // ...then "/" (the catch-all)
		},
	}
}

func (a *AppRouter) BuildRouter(ctx gctx.Context) {
	a.router = chi.NewRouter()
	a.router.Use(middlewares.Recover)

	for _, r := range a.public {
		sub := chi.NewRouter()
		r.PublicRoute(sub)
		a.router.Mount(r.Path(), sub)
		slog.Info("mounted route", "path", r.Path())
	}
}

func (a *AppRouter) Start() error {
	if a.router == nil {
		return errors.New("AppRouter has no router: did you forget to call BuildRouter?")
	}

	addr := ":9001" // override with PORT=...
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	fmt.Printf("\nListening on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, a.router)
}
