package dto

import (
	"encoding/json"
	"time"
)

type FindingResponse struct {
	ID             string          `json:"id"`
	AssessmentID   *string         `json:"assessment_id,omitempty"`
	TaskID         *string         `json:"task_id,omitempty"`
	AgentID        string          `json:"agent_id"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	Severity       string          `json:"severity"`
	TechniqueID    string          `json:"technique_id"`
	TechniqueName  string          `json:"technique_name"`
	Evidence       json.RawMessage `json:"evidence"`
	Recommendation string          `json:"recommendation"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}
