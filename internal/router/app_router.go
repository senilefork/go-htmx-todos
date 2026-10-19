package router

import (
	gctx "context"
	"errors"
	"fmt"
	"net/http"

	"go-htmx-jobs/internal/builders"
	"go-htmx-jobs/internal/cfg"
	appctx "go-htmx-jobs/internal/context"
	"go-htmx-jobs/internal/controllers/private"
	"go-htmx-jobs/internal/controllers/public"
	"go-htmx-jobs/internal/middlewares"

	"github.com/go-chi/chi/v5"
)

type AppRouter struct {
	public  []builders.AppPublicRoute
	private []builders.AppPrivateRoute
	addr    string
	router  chi.Router
	appCtx  *appctx.AppContext
}

func NewAppRouter(appCtx *appctx.AppContext) AppRouter {
	return AppRouter{
		appCtx: appCtx,
		public: []builders.AppPublicRoute{
			public.NewWebPublic(appCtx),
			public.NewTodos(appCtx),
			// INSERT PUBLIC
		},
		private: []builders.AppPrivateRoute{
			private.NewDashboard(appCtx),
			// INSERT PRIVATE
		},
	}
}

func (a *AppRouter) BuildRouter(ctx gctx.Context) {
	lgr := appctx.GetLogger(ctx, "AppRouter.BuildRouter")
	a.router = chi.NewRouter()
	a.router.Use(middlewares.Recover(ctx, a.appCtx))

	// Public group (no auth enforcement)
	a.router.Group(func(g chi.Router) {
		g.Use(middlewares.Auth(a.appCtx, false))
		for _, r := range a.public {
			router := chi.NewRouter()
			r.PublicRoute(router)
			g.Mount(r.Path(), router)
			lgr.Info("mounted public", "path", r.Path())
		}
	})

	// Private group (auth enforced)
	a.router.Group(func(g chi.Router) {
		g.Use(middlewares.Auth(a.appCtx, true))
		for _, r := range a.private {
			builder := builders.NewPrivateRouteBuilder(a.appCtx)
			r.PrivateRoute(ctx, &builder)
			g.Mount(r.Path(), builder.RawRouter())
			lgr.Info("mounted private", "path", r.Path())
		}
	})
}

func (a *AppRouter) Start(cfg cfg.Config) error {
	if a.router == nil {
		return errors.New("AppRouter has no router. Did you forget to call BuildRouter?")
	}
	a.addr = fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)
	fmt.Printf("\nListening on http://%s\n", a.addr)
	return http.ListenAndServe(a.addr, a.router)
}
