package dao

import (
	gctx "context"
	"time"

	"go-htmx-jobs/internal/constant"
	"go-htmx-jobs/internal/context"
	"go-htmx-jobs/internal/database/gen/go_jobs_db/public/table"
	dbmodel "go-htmx-jobs/pkg/db/public/model"
	"github.com/go-jet/jet/v2/postgres"
)

type privilegesDAO struct {
	context.DAOBaseQueries[int64, dbmodel.Privileges]
}

func NewPrivilegesDAO() *privilegesDAO {
	dao := &privilegesDAO{
		DAOBaseQueries: nil,
	}
	queries := newDAOQueryable[int64, dbmodel.Privileges](dao)
	dao.DAOBaseQueries = &queries
	return dao
}

func (dao *privilegesDAO) Table() context.PostgresTable {
	return table.Privileges
}

func (dao *privilegesDAO) InsertCols() postgres.ColumnList {
	return table.Privileges.AllColumns.Except(
		table.Privileges.ID,
		table.Privileges.CreatedAt,
		table.Privileges.UpdatedAt,
	)
}

func (dao *privilegesDAO) UpdateCols() postgres.ColumnList {
	return table.Privileges.AllColumns.Except(
		table.Privileges.ID,
		table.Privileges.CreatedAt,
	)
}

func (dao *privilegesDAO) AllCols() postgres.ColumnList {
	return table.Privileges.AllColumns
}

func (dao *privilegesDAO) OnConflictCols() postgres.ColumnList {
	return []postgres.Column{table.Privileges.Name}
}

func (dao *privilegesDAO) UpdateOnConflictCols() []postgres.ColumnAssigment {
	return []postgres.ColumnAssigment{
		table.Privileges.Name.SET(table.Privileges.EXCLUDED.Name),
	}
}

func (dao *privilegesDAO) PKMatch(pk int64) postgres.BoolExpression {
	return table.Privileges.ID.EQ(postgres.Int(pk))
}

func (dao *privilegesDAO) GetUpdatedAt(row *dbmodel.Privileges) *time.Time {
	return row.UpdatedAt
}

func (dao *privilegesDAO) GetManyByName(ctx gctx.Context, names []constant.PrivilegeName) ([]dbmodel.Privileges, error) {
	var privileges []dbmodel.Privileges

	if len(names) == 0 {
		return privileges, nil
	}

	exprs := make([]postgres.Expression, len(names))
	for i, name := range names {
		exprs[i] = postgres.String(string(name))
	}

	err := table.Privileges.
		SELECT(table.Privileges.AllColumns).
		WHERE(table.Privileges.Name.IN(exprs...)).
		QueryContext(ctx, context.GetDB(ctx), &privileges)
	if err != nil {
		return nil, err
	}

	return privileges, nil
}
