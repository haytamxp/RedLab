package dto

import (
	"encoding/json"
	"time"
)

type CreateTaskRequest struct {
	AgentID  string          `json:"agent_id" binding:"required"`
	Type     string          `json:"type" binding:"required"`
	Payload  json.RawMessage `json:"payload"`
	Priority int             `json:"priority"`
}

type TaskResultRequest struct {
	Status string          `json:"status" binding:"required"`
	Result json.RawMessage `json:"result"`
	Error  string          `json:"error"`
}

type TaskResponse struct {
	ID           string          `json:"id"`
	AgentID      string          `json:"agent_id"`
	Type         string          `json:"type"`
	Payload      json.RawMessage `json:"payload"`
	Status       string          `json:"status"`
	Priority     int             `json:"priority"`
	Result       json.RawMessage `json:"result,omitempty"`
	ErrorMessage *string         `json:"error_message,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	ClaimedAt    *time.Time      `json:"claimed_at,omitempty"`
	CompletedAt  *time.Time      `json:"completed_at,omitempty"`
	ReviewStatus string          `json:"review_status"`
	ReviewedAt   *time.Time      `json:"reviewed_at,omitempty"`
}

type ReviewTaskRequest struct {
	Status string `json:"status" binding:"required"`
}
