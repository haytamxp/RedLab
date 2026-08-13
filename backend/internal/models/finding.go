package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Severity string

const (
	SeverityInfo     Severity = "INFO"
	SeverityLow      Severity = "LOW"
	SeverityMedium   Severity = "MEDIUM"
	SeverityHigh     Severity = "HIGH"
	SeverityCritical Severity = "CRITICAL"
)

type Finding struct {
	Base

	AssessmentID *uuid.UUID `db:"assessment_id" json:"assessment_id,omitempty"`
	TaskID       *uuid.UUID `db:"task_id" json:"task_id,omitempty"`
	AgentID      uuid.UUID  `db:"agent_id" json:"agent_id"`

	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`

	Severity Severity `db:"severity" json:"severity"`

	TechniqueID   string `db:"technique_id" json:"technique_id"`
	TechniqueName string `db:"technique_name" json:"technique_name"`

	Evidence json.RawMessage `db:"evidence" json:"evidence"`

	Recommendation string `db:"recommendation" json:"recommendation"`
}
