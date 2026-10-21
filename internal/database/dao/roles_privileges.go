package dao

import (
	gctx "context"
	"errors"
	"fmt"

	appctx "go-htmx-jobs/internal/context"
	"go-htmx-jobs/pkg/model"
)

type rolesPrivilegesDAO struct{}

// GetMany returns the join rows matching any of the composite keys —
// "does this role hold at least one of these privileges" is then just
// len(rows) > 0, same as the starter's HasPermissionsByIDS.
func (d *rolesPrivilegesDAO) GetMany(ctx gctx.Context, pks []model.RolesPrivilegesPrimaryKey) ([]model.RolesPrivilege, error) {
	if len(pks) == 0 {
		return nil, nil
	}

	where := ""
	args := make([]any, 0, len(pks)*2)
	for i, pk := range pks {
		if i > 0 {
			where += " OR "
		}
		where += fmt.Sprintf("(role_id = $%d AND privilege_id = $%d)", i*2+1, i*2+2)
		args = append(args, pk.RoleID, pk.PrivilegeID)
	}

	rows, err := appctx.GetDB(ctx).QueryContext(ctx,
		`SELECT role_id, privilege_id, created_at
		 FROM roles_privileges
		 WHERE `+where,
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("get roles_privileges: %w", err)
	}
	defer rows.Close()

	var out []model.RolesPrivilege
	for rows.Next() {
		var rp model.RolesPrivilege
		if err := rows.Scan(&rp.RoleID, &rp.PrivilegeID, &rp.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan roles_privilege: %w", err)
		}
		out = append(out, rp)
	}
	return out, errors.Join(rows.Err())
}
