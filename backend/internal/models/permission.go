package models

type Permission struct {
	Base

	Name string `db:"name" json:"name"`

	Description string `db:"description" json:"description"`
}