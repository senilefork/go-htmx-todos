package dao

import (
	"time"

	"go-htmx-jobs/internal/context"
	"go-htmx-jobs/internal/database/gen/go_jobs_db/public/table"
	dbmodel "go-htmx-jobs/pkg/db/public/model"
	"go-htmx-jobs/pkg/model"
	"github.com/go-jet/jet/v2/postgres"
)

type rolesPrivilegesDAO struct {
	context.DAOBaseQueries[model.RolesPrivilegesPrimaryKey, dbmodel.RolesPrivileges]
}

func NewRolesPrivilegesDAO() *rolesPrivilegesDAO {
	dao := &rolesPrivilegesDAO{
		DAOBaseQueries: nil,
	}
	queries := newDAOQueryable[model.RolesPrivilegesPrimaryKey, dbmodel.RolesPrivileges](dao)
	dao.DAOBaseQueries = &queries
	return dao
}

func (dao *rolesPrivilegesDAO) Table() context.PostgresTable {
	return table.RolesPrivileges
}

func (dao *rolesPrivilegesDAO) InsertCols() postgres.ColumnList {
	return table.RolesPrivileges.AllColumns.Except(
		table.RolesPrivileges.CreatedAt,
	)
}

func (dao *rolesPrivilegesDAO) UpdateCols() postgres.ColumnList {
	return table.RolesPrivileges.AllColumns.Except(
		table.RolesPrivileges.CreatedAt,
	)
}

func (dao *rolesPrivilegesDAO) AllCols() postgres.ColumnList {
	return table.RolesPrivileges.AllColumns
}

func (dao *rolesPrivilegesDAO) OnConflictCols() postgres.ColumnList {
	return []postgres.Column{
		table.RolesPrivileges.RoleID,
		table.RolesPrivileges.PrivilegeID,
	}
}

func (dao *rolesPrivilegesDAO) UpdateOnConflictCols() []postgres.ColumnAssigment {
	return []postgres.ColumnAssigment{
		table.RolesPrivileges.RoleID.SET(table.RolesPrivileges.EXCLUDED.RoleID),
		table.RolesPrivileges.PrivilegeID.SET(table.RolesPrivileges.EXCLUDED.PrivilegeID),
	}
}

func (dao *rolesPrivilegesDAO) PKMatch(pk model.RolesPrivilegesPrimaryKey) postgres.BoolExpression {
	return table.RolesPrivileges.
		RoleID.EQ(postgres.Int16(pk.RoleID)).
		AND(table.RolesPrivileges.PrivilegeID.EQ(postgres.Int(pk.PrivilegeID)))
}

func (dao *rolesPrivilegesDAO) GetUpdatedAt(row *dbmodel.RolesPrivileges) *time.Time {
	return nil
}
