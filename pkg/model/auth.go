package model

// Primary-key structs for tables whose PK is composite. Single-column
// PKs (privileges.id, user_roles.okta_sub) need no wrapper — the raw
// value is the key. Row structs live in pkg/db/public/model (generated).

// SessionsPrimaryKey mirrors the table's composite PK (okta_sub, token).
type SessionsPrimaryKey struct {
	OktaSub   string
	AuthToken string
}

// RolesPrivilegesPrimaryKey mirrors the join table's composite PK.
type RolesPrivilegesPrimaryKey struct {
	RoleID      int16
	PrivilegeID int64
}
