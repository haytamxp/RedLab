package services

import (
	"context"
	"errors"

	"github.com/haytamxp/redlab/backend/internal/auth"
	"github.com/haytamxp/redlab/backend/internal/models"
)

type AuthService struct {
	users *UserService

	jwtSecret string

	expiration int
}

func NewAuthService(
	users *UserService,
	secret string,
	expiration int,
) *AuthService {

	return &AuthService{
		users: users,

		jwtSecret: secret,

		expiration: expiration,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	user *models.User,
	password string,
) error {

	hash, err := auth.HashPassword(password)

	if err != nil {
		return err
	}

	user.PasswordHash = hash

	return s.users.CreateUser(
		ctx,
		user,
	)
}

func (s *AuthService) Login(
	ctx context.Context,
	username,
	password string,
) (string, error) {

	user, err := s.users.FindByUsername(
		ctx,
		username,
	)

	if err != nil {
		return "", err
	}

	if !auth.CheckPassword(
		user.PasswordHash,
		password,
	) {

		return "", errors.New("invalid username or password")
	}

	token, err := auth.GenerateJWT(
		user.ID.String(),
		string(user.Role),
		s.jwtSecret,
		s.expiration,
	)

	if err != nil {
		return "", err
	}

	return token, nil
}