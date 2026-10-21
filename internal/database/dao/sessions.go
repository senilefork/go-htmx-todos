package dao

import (
	gctx "context"
	"fmt"

	appctx "go-htmx-jobs/internal/context"
	"go-htmx-jobs/pkg/model"
)

type sessionsDAO struct{}

func (d *sessionsDAO) Insert(ctx gctx.Context, m *model.Sessions) error {
	_, err := appctx.GetDB(ctx).ExecContext(ctx,
		`INSERT INTO sessions (okta_sub, token, expires_at)
		 VALUES ($1, $2, $3)`,
		m.OktaSub, m.Token, m.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("insert session: %w", err)
	}
	return nil
}

// GetOne returns the session matching both halves of the composite PK.
// sql.ErrNoRows passes through unwrapped so callers can
// errors.Is(err, sql.ErrNoRows) — an absent session means "not logged in".
func (d *sessionsDAO) GetOne(ctx gctx.Context, pk model.SessionsPrimaryKey) (*model.Sessions, error) {
	row := appctx.GetDB(ctx).QueryRowContext(ctx,
		`SELECT okta_sub, token, expires_at, created_at
		 FROM sessions
		 WHERE okta_sub = $1 AND token = $2`,
		pk.OktaSub, pk.AuthToken,
	)
	m := &model.Sessions{}
	if err := row.Scan(&m.OktaSub, &m.Token, &m.ExpiresAt, &m.CreatedAt); err != nil {
		return nil, err
	}
	return m, nil
}

func (d *sessionsDAO) Delete(ctx gctx.Context, pk model.SessionsPrimaryKey) error {
	_, err := appctx.GetDB(ctx).ExecContext(ctx,
		`DELETE FROM sessions WHERE okta_sub = $1 AND token = $2`,
		pk.OktaSub, pk.AuthToken,
	)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}
