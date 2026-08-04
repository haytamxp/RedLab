package models

type Asset struct {
	Base

	Name string `db:"name" json:"name"`

	IPAddress string `db:"ip_address" json:"ip_address"`

	Hostname string `db:"hostname" json:"hostname"`

	OperatingSystem string `db:"operating_system" json:"operating_system"`
}