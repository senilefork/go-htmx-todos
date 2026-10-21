package dao

import appctx "go-htmx-jobs/internal/context"

// Manager hands out the DAO singletons. DAOs are stateless — the pool
// comes from ctx on every call — so there's nothing to configure; the
// struct exists to own the instances and satisfy appctx.DAOManager.
type Manager struct {
	sessions        sessionsDAO
	userRoles       userRolesDAO
	privileges      privilegesDAO
	rolesPrivileges rolesPrivilegesDAO
}

func NewDAOManager() *Manager {
	return &Manager{}
}

func (m *Manager) SessionsDAO() appctx.SessionsDAO             { return &m.sessions }
func (m *Manager) UserRolesDAO() appctx.UserRolesDAO           { return &m.userRoles }
func (m *Manager) PrivilegesDAO() appctx.PrivilegesDAO         { return &m.privileges }
func (m *Manager) RolesPrivilegesDAO() appctx.RolesPrivilegesDAO { return &m.rolesPrivileges }
