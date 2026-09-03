package models

import "github.com/google/uuid"

type Department struct {
	Base

	Name        string `db:"name" json:"name"`
	Code        string `db:"code" json:"code"`
	Description string `db:"description" json:"description"`
	Teams       []Team `json:"teams,omitempty"`
	UserCount   int    `json:"user_count"`
}

type Team struct {
	Base

	DepartmentID uuid.UUID `db:"department_id" json:"department_id"`
	Name         string    `db:"name" json:"name"`
	Code         string    `db:"code" json:"code"`
	Description  string    `db:"description" json:"description"`
	UserCount    int       `json:"user_count"`
}
