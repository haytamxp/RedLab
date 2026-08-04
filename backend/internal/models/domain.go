package models

type Domain struct {
	Base

	Name string `db:"name" json:"name"`

	FQDN string `db:"fqdn" json:"fqdn"`

	Forest string `db:"forest" json:"forest"`
}