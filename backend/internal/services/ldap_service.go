package services

import (
	"context"
	"fmt"

	"github.com/haytamxp/redlab/backend/internal/config"
	redlabldap "github.com/haytamxp/redlab/backend/internal/ldap"
)

type LDAPService struct {
	client *redlabldap.Client
}

func NewLDAPService(cfg config.LDAPConfig) *LDAPService {
	return &LDAPService{
		client: redlabldap.NewClient(cfg),
	}
}

func (s *LDAPService) Authenticate(
	ctx context.Context,
	username string,
	password string,
) (*redlabldap.User, error) {

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	user, err := s.client.Authenticate(username, password)
	if err != nil {
		return nil, fmt.Errorf("LDAP authentication failed: %w", err)
	}

	if err := redlabldap.ValidateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *LDAPService) FindUser(
	ctx context.Context,
	username string,
) (*redlabldap.User, error) {

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	user, err := s.client.FindUser(username)
	if err != nil {
		return nil, err
	}

	if err := redlabldap.ValidateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *LDAPService) ListUsers(
	ctx context.Context,
	search string,
) ([]redlabldap.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return s.client.ListUsers(search)
}

func (s *LDAPService) Close() {
	s.client.Close()
}
