package dto

import "time"

type AgentResponse struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Hostname        string     `json:"hostname"`
	IPAddress       string     `json:"ip_address"`
	OperatingSystem string     `json:"operating_system"`
	Version         string     `json:"version"`
	Status          string     `json:"status"`
	LastSeen        *time.Time `json:"last_seen,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
