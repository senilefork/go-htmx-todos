package dao

import (
	"time"

	"go-htmx-jobs/internal/context"
	"go-htmx-jobs/internal/database/gen/go_jobs_db/public/table"
	dbmodel "go-htmx-jobs/pkg/db/public/model"
	"github.com/go-jet/jet/v2/postgres"
)

type userRolesDAO struct {
	context.DAOBaseQueries[string, dbmodel.UserRoles]
}

func NewUserRolesDAO() *userRolesDAO {
	dao := &userRolesDAO{
		DAOBaseQueries: nil,
	}
	queries := newDAOQueryable[string, dbmodel.UserRoles](dao)
	dao.DAOBaseQueries = &queries
	return dao
}

func (dao *userRolesDAO) Table() context.PostgresTable {
	return table.UserRoles
}

func (dao *userRolesDAO) InsertCols() postgres.ColumnList {
	return table.UserRoles.AllColumns.Except(
		table.UserRoles.CreatedAt,
		table.UserRoles.UpdatedAt,
	)
}

func (dao *userRolesDAO) UpdateCols() postgres.ColumnList {
	return table.UserRoles.AllColumns.Except(
		table.UserRoles.CreatedAt,
		table.UserRoles.UpdatedAt,
		table.UserRoles.OktaSub,
	)
}

func (dao *userRolesDAO) AllCols() postgres.ColumnList {
	return table.UserRoles.AllColumns
}

func (dao *userRolesDAO) OnConflictCols() postgres.ColumnList {
	return []postgres.Column{table.UserRoles.OktaSub}
}

func (dao *userRolesDAO) UpdateOnConflictCols() []postgres.ColumnAssigment {
	return []postgres.ColumnAssigment{
		table.UserRoles.RoleID.SET(table.UserRoles.EXCLUDED.RoleID),
		table.UserRoles.UpdatedAt.SET(table.UserRoles.EXCLUDED.UpdatedAt),
	}
}

func (dao *userRolesDAO) PKMatch(pk string) postgres.BoolExpression {
	return table.UserRoles.OktaSub.EQ(postgres.String(pk))
}

func (dao *userRolesDAO) GetUpdatedAt(row *dbmodel.UserRoles) *time.Time {
	return row.UpdatedAt
}
