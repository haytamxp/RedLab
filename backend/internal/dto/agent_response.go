package dto

type AgentResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Hostname        string `json:"hostname"`
	IPAddress       string `json:"ip_address"`
	OperatingSystem string `json:"operating_system"`
	Version         string `json:"version"`
	Status          string `json:"status"`
}