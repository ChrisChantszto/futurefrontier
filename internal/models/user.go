package models

import "time"

type Role string

const (
	RoleSuperAdmin Role = "superadmin"
	RoleAdmin      Role = "admin"
	RoleEditor     Role = "editor"
	RoleViewer     Role = "viewer"
)

type User struct {
	ID            string   `bson:"_id,omitempty" json:"id"`
	Email         string   `bson:"email" json:"email"`
	PasswordHash  string   `bson:"passwordHash" json:"-"`
	Roles         []Role   `bson:"roles" json:"roles"`
	EmailVerified bool     `bson:"emailVerified" json:"emailVerified"`
	CreatedAt     time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time `bson:"updatedAt" json:"updatedAt"`
}