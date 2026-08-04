package models

import "github.com/google/uuid"

type Report struct {
	Base

	AssessmentID uuid.UUID `db:"assessment_id" json:"assessment_id"`

	FileName string `db:"file_name" json:"file_name"`

	FilePath string `db:"file_path" json:"file_path"`
}