package dao

import (
	gctx "context"
	"fmt"

	appctx "go-htmx-jobs/internal/context"
	"go-htmx-jobs/pkg/model"
)

type userRolesDAO struct{}

func (d *userRolesDAO) GetOne(ctx gctx.Context, oktaSub string) (*model.UserRole, error) {
	row := appctx.GetDB(ctx).QueryRowContext(ctx,
		`SELECT okta_sub, role_id, created_at, updated_at
		 FROM user_roles
		 WHERE okta_sub = $1`,
		oktaSub,
	)
	m := &model.UserRole{}
	if err := row.Scan(&m.OktaSub, &m.RoleID, &m.CreatedAt, &m.UpdatedAt); err != nil {
		return nil, err
	}
	return m, nil
}

// Upsert assigns (or reassigns) the role for a subject. The no-op-write
// ON CONFLICT form keeps it one round trip; updated_at is touched so a
// reassignment is visible in the data.
func (d *userRolesDAO) Upsert(ctx gctx.Context, m *model.UserRole) error {
	_, err := appctx.GetDB(ctx).ExecContext(ctx,
		`INSERT INTO user_roles (okta_sub, role_id)
		 VALUES ($1, $2)
		 ON CONFLICT (okta_sub)
		 DO UPDATE SET role_id = EXCLUDED.role_id,
		               updated_at = CURRENT_TIMESTAMP`,
		m.OktaSub, m.RoleID,
	)
	if err != nil {
		return fmt.Errorf("upsert user role: %w", err)
	}
	return nil
}
