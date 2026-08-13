package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TaskPending   TaskStatus = "PENDING"
	TaskClaimed   TaskStatus = "CLAIMED"
	TaskCompleted TaskStatus = "COMPLETED"
	TaskFailed    TaskStatus = "FAILED"
)

type Task struct {
	Base

	AgentID uuid.UUID `json:"agent_id"`

	Type string `json:"type"`

	Payload json.RawMessage `json:"payload"`

	Status TaskStatus `json:"status"`

	Priority int `json:"priority"`

	Result *json.RawMessage `json:"result,omitempty"`

	ErrorMessage *string `json:"error_message,omitempty"`

	ClaimedAt *time.Time `json:"claimed_at,omitempty"`

	CompletedAt *time.Time `json:"completed_at,omitempty"`
}
