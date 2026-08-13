package validator

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/haytamxp/redlab/backend/internal/dto"
)

var (
	ErrAssessmentNameRequired = errors.New(
		"assessment name is required",
	)
	ErrAssessmentNameTooLong = errors.New(
		"assessment name exceeds 120 characters",
	)
	ErrInvalidAssessmentAgent = errors.New(
		"invalid assessment agent ID",
	)
)

func ValidateCreateAssessment(
	request dto.CreateAssessmentRequest,
) error {
	name := strings.TrimSpace(
		request.Name,
	)

	if name == "" {
		return ErrAssessmentNameRequired
	}

	if len(name) > 120 {
		return ErrAssessmentNameTooLong
	}

	if _, err := uuid.Parse(
		request.AgentID,
	); err != nil {
		return ErrInvalidAssessmentAgent
	}

	return nil
}
