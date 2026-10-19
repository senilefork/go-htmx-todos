package private

import (
	gctx "context"
	"net/http"

	appctx "go-htmx-jobs/internal/context"
	"go-htmx-jobs/internal/builders"
	"go-htmx-jobs/internal/constant"
)

type Dashboard struct {
	ctx *appctx.AppContext
}

func NewDashboard(ctx *appctx.AppContext) *Dashboard {
	return &Dashboard{ctx: ctx}
}

// Path is the mount prefix for this controller.
func (d *Dashboard) Path() string {
	return "/dashboard"
}

// PrivateRoute registers private handlers for dashboard.
func (d *Dashboard) PrivateRoute(ctx gctx.Context, b *builders.PrivateRouteBuilder) {
	b.NewHandler().
		Register(http.MethodGet, "/", d.getDashboard).
		SetRequiredPrivileges(constant.DashboardView).
		Build(ctx)
}

func (d *Dashboard) getDashboard(w http.ResponseWriter, r *http.Request) {
	lgr := appctx.GetLogger(r.Context(), "Dashboard.getDashboard")
	lgr.Info("dashboard private endpoint hit")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("dashboard: private endpoint (stub)"))
}
