package cmd

import (
	gctx "context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"go-htmx-jobs/internal/cfg"
	"go-htmx-jobs/internal/context"
	"go-htmx-jobs/internal/router"
	"go-htmx-jobs/internal/services"
)

// web is the real bootstrap. Wiring order matters and matches the starter:
//   services -> AppContext (DI container) -> router -> listen.
// (The starter additionally opens a DB and puts it on the context here.)
func web() {
	ctx := gctx.Background()

	baseURL := os.Getenv("JSONPLACEHOLDER_BASE_URL")
	if baseURL == "" {
		baseURL = services.DefaultBaseURL
	}
	client := &http.Client{Timeout: 5 * time.Second}

	sm := services.NewServiceManager(baseURL, client)
	appCtx := context.NewAppContext(sm)

	appRouter := router.NewAppRouter(appCtx)
	appRouter.BuildRouter(ctx)

	if err := appRouter.Start(cfg.Load()); err != nil {
		slog.Error("failed to start router", "error", err)
		os.Exit(1)
	}
}
