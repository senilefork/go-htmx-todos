package cmd

import (
	gctx "context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	"go-htmx-jobs/internal/cfg"
	"go-htmx-jobs/internal/context"
	"go-htmx-jobs/internal/database/dao"
	"go-htmx-jobs/internal/router"
	"go-htmx-jobs/internal/services"

	_ "github.com/lib/pq"
)

// web is the real bootstrap. Wiring order matters:
// config -> db pool -> services -> AppContext (DI container) -> router -> listen.
func web() {
	config := cfg.Load()

	// sql.Open is lazy: it only validates the URL and creates the pool
	// struct — no TCP happens yet. Ping forces a real connection now so a
	// dead container fails the boot instead of the first request.
	db, err := sql.Open("postgres", config.DbUrl())
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	// and the happy path blocks in ListenAndServe forever. To make it real
	// (and to stop cutting in-flight requests off mid-response on deploy),
	// add graceful shutdown later:
	//   1. create an explicit http.Server in router.Start instead of using
	//      the package-level http.ListenAndServe
	//   2. in web(), watch for SIGINT/SIGTERM with signal.NotifyContext
	//   3. on signal, call server.Shutdown(ctx) to drain in-flight requests
	//   4. only then does this Close (and the process exit) run cleanly
	defer db.Close()
	if err := db.PingContext(gctx.Background()); err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	ctx := gctx.Background()
	ctx = context.WithDB(ctx, db)

	baseURL := os.Getenv("JSONPLACEHOLDER_BASE_URL")
	if baseURL == "" {
		baseURL = services.DefaultBaseURL
	}
	client := &http.Client{Timeout: 5 * time.Second}

	sm := services.NewServiceManager(baseURL, client)
	dm := dao.NewDAOManager()
	appCtx := context.NewAppContext(sm, dm)

	appRouter := router.NewAppRouter(appCtx)
	appRouter.BuildRouter(ctx)

	if err := appRouter.Start(config); err != nil {
		slog.Error("failed to start router", "error", err)
		os.Exit(1)
	}
}
