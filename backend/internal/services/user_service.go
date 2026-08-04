package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/haytamxp/redlab/backend/internal/models"
	"github.com/haytamxp/redlab/backend/internal/repository"
)

type UserService struct {
	repository *repository.UserRepository
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) Create(ctx context.Context, user *models.User) error {

	user.ID = uuid.New()

	user.IsActive = true

	return s.repository.Create(ctx, user)
}

func (s *UserService) FindByUsername(ctx context.Context, username string) (*models.User, error) {

	return s.repository.FindByUsername(ctx, username)
}
func (s *UserService) GetAll(
	ctx context.Context,
) ([]models.User, error) {

	return s.repository.FindAll(ctx)
}