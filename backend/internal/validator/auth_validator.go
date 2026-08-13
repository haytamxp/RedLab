package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/haytamxp/redlab/backend/internal/dto"
)

var authValidator = validator.New()

func ValidateRegisterRequest(req dto.RegisterRequest) error {
	if err := authValidator.Struct(req); err != nil {
		return fmt.Errorf("register validation failed: %w", err)
	}

	return nil
}

func ValidateLoginRequest(req dto.LoginRequest) error {
	if err := authValidator.Struct(req); err != nil {
		return fmt.Errorf("login validation failed: %w", err)
	}

	return nil
}
