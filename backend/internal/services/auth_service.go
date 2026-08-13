package services

import (
	"context"
	"errors"

	"github.com/haytamxp/redlab/backend/internal/auth"
	"github.com/haytamxp/redlab/backend/internal/models"
)

type AuthService struct {
	users *UserService
	ldap  *LDAPService

	secret     string
	expiration int
}

func NewAuthService(
	users *UserService,
	ldap *LDAPService,
	secret string,
	expiration int,
) *AuthService {

	return &AuthService{
		users:      users,
		ldap:       ldap,
		secret:     secret,
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

	return s.users.Create(ctx, user)
}

func (s *AuthService) Login(
	ctx context.Context,
	username string,
	password string,
) (string, error) {

	if username == "" || password == "" {
		return "", errors.New("username and password are required")
	}

	// Authenticate the credentials against Active Directory.
	ldapUser, err := s.ldap.Authenticate(ctx, username, password)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// Find the corresponding RedLab application user.
	user, err := s.users.FindByUsername(ctx, ldapUser.SAMAccountName)
	if err != nil {
		return "", errors.New("user is not registered in RedLab")
	}

	if !user.IsActive {
		return "", errors.New("user account is inactive")
	}

	// Generate the RedLab JWT using the application's role.
	return auth.GenerateJWT(
		user.ID.String(),
		string(user.Role),
		s.secret,
		s.expiration,
	)
}
