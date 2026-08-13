package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/haytamxp/redlab/backend/internal/dto"
)

var findingValidator = validator.New()

func ValidateFindingRequest(req dto.FindingRequest) error {
	if err := findingValidator.Struct(req); err != nil {
		return fmt.Errorf("finding validation failed: %w", err)
	}

	return nil
}
