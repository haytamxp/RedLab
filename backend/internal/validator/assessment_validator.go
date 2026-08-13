package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/haytamxp/redlab/backend/internal/dto"
)

var assessmentValidator = validator.New()

func ValidateAssessmentRequest(req dto.AssessmentRequest) error {
	if err := assessmentValidator.Struct(req); err != nil {
		return fmt.Errorf("assessment validation failed: %w", err)
	}

	return nil
}
