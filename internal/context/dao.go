package context

import (
	gctx "context"
	"time"

	"go-htmx-jobs/internal/constant"
	dbmodel "go-htmx-jobs/pkg/db/public/model"
	"go-htmx-jobs/pkg/model"
	"github.com/go-jet/jet/v2/postgres"
)

type PostgresTable interface {
	postgres.WritableTable
	postgres.ReadableTable
}

type PrimaryKey any

type GetBaseCols interface {
	InsertCols() postgres.ColumnList
	UpdateCols() postgres.ColumnList
	AllCols() postgres.ColumnList
	OnConflictCols() postgres.ColumnList
	UpdateOnConflictCols() []postgres.ColumnAssigment
}

type PKMatcher[PK any] interface {
	PKMatch(pk PK) postgres.BoolExpression
}

type GetUpdatedAt[R any] interface {
	GetUpdatedAt(row *R) *time.Time
}

type GetTable interface {
	Table() PostgresTable
}

type DAOBaseQueries[PK PrimaryKey, R any] interface {
	GetAll(ctx gctx.Context) ([]R, error)
	GetOne(ctx gctx.Context, pk PK) (R, error)
	GetMany(ctx gctx.Context, pks []PK) ([]R, error)
	Insert(ctx gctx.Context, model *R) error
	InsertMany(ctx gctx.Context, models []R) ([]R, error)
	Upsert(ctx gctx.Context, model *R) error
	UpsertMany(ctx gctx.Context, models []R) ([]R, error)
	Update(ctx gctx.Context, model *R, pk PK) error
	Delete(ctx gctx.Context, pk PK) error
}

type DAO[PK any, R any] interface {
	GetTable
	GetBaseCols
	PKMatcher[PK]
	GetUpdatedAt[R]
	DAOBaseQueries[PK, R]
}

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
	DAO[model.SessionsPrimaryKey, dbmodel.Sessions]
}

type UserRolesDAO interface {
	DAO[string, dbmodel.UserRoles]
}

type PrivilegesDAO interface {
	DAO[int64, dbmodel.Privileges]
	GetManyByName(ctx gctx.Context, names []constant.PrivilegeName) ([]dbmodel.Privileges, error)
}

type RolesPrivilegesDAO interface {
	DAO[model.RolesPrivilegesPrimaryKey, dbmodel.RolesPrivileges]
}
