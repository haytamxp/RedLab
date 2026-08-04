package models

import "github.com/google/uuid"

type Severity string

const (
	SeverityLow      Severity = "LOW"
	SeverityMedium   Severity = "MEDIUM"
	SeverityHigh     Severity = "HIGH"
	SeverityCritical Severity = "CRITICAL"
)

type Finding struct {
	Base

	AssessmentID uuid.UUID `db:"assessment_id" json:"assessment_id"`

	Title string `db:"title" json:"title"`

	Description string `db:"description" json:"description"`

	Severity Severity `db:"severity" json:"severity"`
}