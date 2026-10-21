package dao

import (
	"time"

	"go-htmx-jobs/internal/context"
	"go-htmx-jobs/internal/database/gen/go_jobs_db/public/table"
	dbmodel "go-htmx-jobs/pkg/db/public/model"
	"go-htmx-jobs/pkg/model"
	"github.com/go-jet/jet/v2/postgres"
)

type sessionsDAO struct {
	context.DAOBaseQueries[model.SessionsPrimaryKey, dbmodel.Sessions]
}

func NewSessionsDAO() *sessionsDAO {
	dao := &sessionsDAO{
		DAOBaseQueries: nil,
	}
	queries := newDAOQueryable[model.SessionsPrimaryKey, dbmodel.Sessions](dao)
	dao.DAOBaseQueries = &queries
	return dao
}

func (dao *sessionsDAO) Table() context.PostgresTable {
	return table.Sessions
}

func (dao *sessionsDAO) InsertCols() postgres.ColumnList {
	return table.Sessions.AllColumns.Except(
		table.Sessions.CreatedAt,
	)
}

func (dao *sessionsDAO) UpdateCols() postgres.ColumnList {
	return table.Sessions.AllColumns.Except(
		table.Sessions.CreatedAt,
		table.Sessions.OktaSub,
		table.Sessions.Token,
	)
}

func (dao *sessionsDAO) AllCols() postgres.ColumnList {
	return table.Sessions.AllColumns
}

func (dao *sessionsDAO) OnConflictCols() postgres.ColumnList {
	return []postgres.Column{}
}

func (dao *sessionsDAO) UpdateOnConflictCols() []postgres.ColumnAssigment {
	return []postgres.ColumnAssigment{}
}

func (dao *sessionsDAO) PKMatch(pk model.SessionsPrimaryKey) postgres.BoolExpression {
	return table.Sessions.
		OktaSub.EQ(postgres.String(pk.OktaSub)).
		AND(table.Sessions.Token.EQ(postgres.String(pk.AuthToken)))
}

func (dao *sessionsDAO) GetUpdatedAt(row *dbmodel.Sessions) *time.Time {
	return nil
}
