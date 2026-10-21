package dao

import (
	appctx "go-htmx-jobs/internal/context"
)

type daoManager struct {
	sessionsDAO        appctx.SessionsDAO
	userRolesDAO       appctx.UserRolesDAO
	privilegesDAO      appctx.PrivilegesDAO
	rolesPrivilegesDAO appctx.RolesPrivilegesDAO
}

func NewDAOManager() appctx.DAOManager {
	return &daoManager{}
}

func (dm *daoManager) SessionsDAO() appctx.SessionsDAO {
	if dm.sessionsDAO == nil {
		dm.sessionsDAO = NewSessionsDAO()
	}
	return dm.sessionsDAO
}

func (dm *daoManager) UserRolesDAO() appctx.UserRolesDAO {
	if dm.userRolesDAO == nil {
		dm.userRolesDAO = NewUserRolesDAO()
	}
	return dm.userRolesDAO
}

func (dm *daoManager) PrivilegesDAO() appctx.PrivilegesDAO {
	if dm.privilegesDAO == nil {
		dm.privilegesDAO = NewPrivilegesDAO()
	}
	return dm.privilegesDAO
}

func (dm *daoManager) RolesPrivilegesDAO() appctx.RolesPrivilegesDAO {
	if dm.rolesPrivilegesDAO == nil {
		dm.rolesPrivilegesDAO = NewRolesPrivilegesDAO()
	}
	return dm.rolesPrivilegesDAO
}
