package dto

type CreateAssessmentRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	AgentID     string `json:"agent_id" binding:"required"`
}