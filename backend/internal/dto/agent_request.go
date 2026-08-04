package dto

type CreateAgentRequest struct {

	Name string `json:"name" binding:"required"`

	Hostname string `json:"hostname" binding:"required"`

	IPAddress string `json:"ip_address" binding:"required"`

	OperatingSystem string `json:"operating_system"`

	Version string `json:"version"`
}