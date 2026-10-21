package dao

import (
	gctx "context"
	"fmt"

	"github.com/lib/pq"

	appctx "go-htmx-jobs/internal/context"
	"go-htmx-jobs/internal/constant"
	"go-htmx-jobs/pkg/model"
)

type privilegesDAO struct{}

// GetManyByName fetches privileges whose names are in the slice. ANY($1)
// with a pq.Array beats building $1..$n placeholders by hand — the slice
// binds as a Postgres text[] in one round trip.
func (d *privilegesDAO) GetManyByName(ctx gctx.Context, names []constant.PrivilegeName) ([]model.Privilege, error) {
	if len(names) == 0 {
		return nil, nil
	}
	strs := make([]string, len(names))
	for i, n := range names {
		strs[i] = string(n)
	}

	rows, err := appctx.GetDB(ctx).QueryContext(ctx,
		`SELECT id, name, created_at, updated_at
		 FROM privileges
		 WHERE name = ANY($1)`,
		pq.Array(strs),
	)
	if err != nil {
		return nil, fmt.Errorf("get privileges by name: %w", err)
	}
	defer rows.Close()

	var out []model.Privilege
	for rows.Next() {
		var p model.Privilege
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan privilege: %w", err)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// UpsertMany inserts any names not yet present and returns the full set
// with ids — existing rows come back via the ON CONFLICT no-op update,
// which forces RETURNING to fire for them too.
func (d *privilegesDAO) UpsertMany(ctx gctx.Context, privileges []model.Privilege) ([]model.Privilege, error) {
	if len(privileges) == 0 {
		return nil, nil
	}
	names := make([]string, len(privileges))
	for i, p := range privileges {
		names[i] = p.Name
	}

	rows, err := appctx.GetDB(ctx).QueryContext(ctx,
		`INSERT INTO privileges (name)
		 SELECT * FROM UNNEST($1::text[])
		 ON CONFLICT (name)
		 DO UPDATE SET updated_at = CURRENT_TIMESTAMP
		 RETURNING id, name, created_at, updated_at`,
		pq.Array(names),
	)
	if err != nil {
		return nil, fmt.Errorf("upsert privileges: %w", err)
	}
	defer rows.Close()

	var out []model.Privilege
	for rows.Next() {
		var p model.Privilege
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan privilege: %w", err)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
