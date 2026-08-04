package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Base

	Username     string    `db:"username" json:"username"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	FirstName    string    `db:"first_name" json:"first_name"`
	LastName     string    `db:"last_name" json:"last_name"`

	Role Role `db:"role" json:"role"`

	IsActive bool `db:"is_active" json:"is_active"`

	LDAPUser bool `db:"ldap_user" json:"ldap_user"`

	LastLogin *time.Time `db:"last_login" json:"last_login,omitempty"`

	ManagerID *uuid.UUID `db:"manager_id" json:"manager_id,omitempty"`
}