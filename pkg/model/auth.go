package model

import "time"

// Sessions is one server-side login session for an Okta subject.
type Sessions struct {
	OktaSub   string
	Token     string
	ExpiresAt *time.Time // nullable in the table
	CreatedAt time.Time
}

// SessionsPrimaryKey mirrors the table's composite PK (okta_sub, token).
type SessionsPrimaryKey struct {
	OktaSub   string
	AuthToken string
}

// UserRole maps an Okta subject to an app role.
type UserRole struct {
	OktaSub   string
	RoleID    int16
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Privilege is one named capability checked by ApplyPrivileges.
type Privilege struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// RolesPrivilege is one row of the role↔privilege join table.
type RolesPrivilege struct {
	RoleID      int16
	PrivilegeID int64
	CreatedAt   time.Time
}

// RolesPrivilegesPrimaryKey mirrors the join table's composite PK.
type RolesPrivilegesPrimaryKey struct {
	RoleID      int16
	PrivilegeID int64
}
