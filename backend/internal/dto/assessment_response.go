package dto

type AssessmentResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AgentID     string `json:"agent_id"`
	Status      string `json:"status"`
}
