package services

import (
	"context"
	"errors"
	"strings"

	"github.com/haytamxp/redlab/backend/internal/auth"
	"github.com/haytamxp/redlab/backend/internal/models"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrUserNotFound       = errors.New("user not registered in RedLab")
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
	user.LDAPUser = false
	user.IsActive = true

	return s.users.Create(ctx, user)
}

func (s *AuthService) Login(
	ctx context.Context,
	username string,
	password string,
) (string, error) {
	username = strings.TrimSpace(username)

	if username == "" || password == "" {
		return "", ErrInvalidCredentials
	}

	user, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		return "", ErrUserNotFound
	}

	if !user.IsActive {
		return "", ErrUserInactive
	}

	// Local RedLab account.
	if !user.LDAPUser {
		if !auth.CheckPassword(user.PasswordHash, password) {
			return "", ErrInvalidCredentials
		}

		if err := s.users.UpdateLastLogin(
			ctx,
			user.ID,
		); err != nil {
			return "", err
		}

		return auth.GenerateJWT(
			user.ID.String(),
			string(user.Role),
			s.secret,
			s.expiration,
		)
	}

	// Active Directory account.
	ldapUser, err := s.ldap.Authenticate(
		ctx,
		username,
		password,
	)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if !strings.EqualFold(
		ldapUser.SAMAccountName,
		user.Username,
	) {
		return "", ErrInvalidCredentials
	}

	if err := s.users.UpdateLastLogin(
		ctx,
		user.ID,
	); err != nil {
		return "", err
	}

	return auth.GenerateJWT(
		user.ID.String(),
		string(user.Role),
		s.secret,
		s.expiration,
	)
}
