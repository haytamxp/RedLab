package models

import "time"

type AgentStatus string

const (
	AgentOnline  AgentStatus = "ONLINE"
	AgentOffline AgentStatus = "OFFLINE"
)

type Agent struct {
	Base

	Name            string      `db:"name" json:"name"`
	Hostname        string      `db:"hostname" json:"hostname"`
	IPAddress       string      `db:"ip_address" json:"ip_address"`
	OperatingSystem string      `db:"operating_system" json:"operating_system"`
	Version         string      `db:"version" json:"version"`
	Status          AgentStatus `db:"status" json:"status"`

	LastSeen *time.Time `db:"last_seen" json:"last_seen,omitempty"`
}