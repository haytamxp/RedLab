package models

import "github.com/google/uuid"

type AssessmentStatus string

const (
	AssessmentPending   AssessmentStatus = "PENDING"
	AssessmentRunning   AssessmentStatus = "RUNNING"
	AssessmentCompleted AssessmentStatus = "COMPLETED"
	AssessmentFailed    AssessmentStatus = "FAILED"
)

type Assessment struct {
	Base

	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`

	AgentID uuid.UUID `db:"agent_id" json:"agent_id"`

	Status AssessmentStatus `db:"status" json:"status"`
}
