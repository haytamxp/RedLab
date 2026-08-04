package services

import (
	"context"

	"github.com/haytamxp/redlab/backend/internal/models"
	"github.com/haytamxp/redlab/backend/internal/repository"
)

type FindingService struct {
	repository *repository.FindingRepository
}

func NewFindingService(
	repository *repository.FindingRepository,
) *FindingService {

	return &FindingService{
		repository: repository,
	}
}

func (s *FindingService) Create(
	ctx context.Context,
	finding *models.Finding,
) error {

	return s.repository.Create(ctx, finding)
}

func (s *FindingService) FindAll(
	ctx context.Context,
) ([]models.Finding, error) {

	return s.repository.FindAll(ctx)
}