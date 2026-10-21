package context

import (
	gctx "context"

	"go-htmx-jobs/internal/constant"
	"go-htmx-jobs/pkg/model"
)

// DAOManager hands out one DAO per table. These interfaces live in this
// package (like the starter's) so AppContext can expose them without
// importing internal/database/dao — which imports us for GetDB. That
// would be an import cycle.

type DAOManager interface {
	SessionsDAO() SessionsDAO
	UserRolesDAO() UserRolesDAO
	PrivilegesDAO() PrivilegesDAO
	RolesPrivilegesDAO() RolesPrivilegesDAO
}

type SessionsDAO interface {
	Insert(ctx gctx.Context, m *model.Sessions) error
	GetOne(ctx gctx.Context, pk model.SessionsPrimaryKey) (*model.Sessions, error)
	Delete(ctx gctx.Context, pk model.SessionsPrimaryKey) error
}

type UserRolesDAO interface {
	GetOne(ctx gctx.Context, oktaSub string) (*model.UserRole, error)
	Upsert(ctx gctx.Context, m *model.UserRole) error
}

type PrivilegesDAO interface {
	GetManyByName(ctx gctx.Context, names []constant.PrivilegeName) ([]model.Privilege, error)
	UpsertMany(ctx gctx.Context, privileges []model.Privilege) ([]model.Privilege, error)
}

type RolesPrivilegesDAO interface {
	GetMany(ctx gctx.Context, pks []model.RolesPrivilegesPrimaryKey) ([]model.RolesPrivilege, error)
}
