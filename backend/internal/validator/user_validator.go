package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/haytamxp/redlab/backend/internal/dto"
)

var userValidator = validator.New()

func ValidateCreateUserRequest(req dto.CreateUserRequest) error {
	if err := userValidator.Struct(req); err != nil {
		return fmt.Errorf("user validation failed: %w", err)
	}

	return nil
}
