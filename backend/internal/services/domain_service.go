package services

import (
	"context"

	"github.com/haytamxp/redlab/backend/internal/models"
	"github.com/haytamxp/redlab/backend/internal/repository"
)

type DomainService struct {
	repository *repository.DomainRepository
}

func NewDomainService(
	repository *repository.DomainRepository,
) *DomainService {

	return &DomainService{
		repository: repository,
	}
}

func (s *DomainService) Create(
	ctx context.Context,
	domain *models.Domain,
) error {

	return s.repository.Create(ctx, domain)
}

func (s *DomainService) FindAll(
	ctx context.Context,
) ([]models.Domain, error) {

	return s.repository.FindAll(ctx)
}