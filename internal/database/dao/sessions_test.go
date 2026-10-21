package dao

import (
	gctx "context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"

	appctx "go-htmx-jobs/internal/context"
	"go-htmx-jobs/pkg/model"
)

// Requires the Docker Postgres container. Skips instead of failing so
// `make test` still passes when it's down.
func TestSessionsRoundTrip(t *testing.T) {
	// Tests run with cwd = this package's dir, so .env at the project
	// root isn't autoloaded — load it relative to here.
	_ = godotenv.Load("../../../.env")

	db, err := sql.Open("postgres", cfgURL(t))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.PingContext(gctx.Background()); err != nil {
		t.Skipf("postgres unavailable: %v", err)
	}

	ctx := appctx.WithDB(gctx.Background(), db)
	dm := NewDAOManager()
	pk := model.SessionsPrimaryKey{OktaSub: "00utest", AuthToken: "tok-integration"}
	t.Cleanup(func() {
		_ = dm.SessionsDAO().Delete(ctx, pk) // don't leak rows on mid-test failure
	})

	exp := time.Now().Add(time.Hour)
	if err := dm.SessionsDAO().Insert(ctx, &model.Sessions{
		OktaSub: pk.OktaSub, Token: pk.AuthToken, ExpiresAt: &exp,
	}); err != nil {
		t.Fatalf("insert: %v", err)
	}

	got, err := dm.SessionsDAO().GetOne(ctx, pk)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.OktaSub != pk.OktaSub || got.Token != pk.AuthToken || got.ExpiresAt == nil {
		t.Fatalf("wrong row: %+v", got)
	}

	if err := dm.SessionsDAO().Delete(ctx, pk); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := dm.SessionsDAO().GetOne(ctx, pk); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("want sql.ErrNoRows after delete, got %v", err)
	}
}

func cfgURL(t *testing.T) string {
	t.Helper()
	if v := os.Getenv("DB_USER"); v == "" {
		t.Fatal("env not loaded — godotenv.Load failed?")
	}
	return fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_EXTERNAL_PORT"), os.Getenv("DB_NAME"))
}
