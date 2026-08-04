package services

import (
	"context"

	"github.com/haytamxp/redlab/backend/internal/models"
	"github.com/haytamxp/redlab/backend/internal/repository"
)

type AssetService struct {
	repository *repository.AssetRepository
}

func NewAssetService(
	repository *repository.AssetRepository,
) *AssetService {

	return &AssetService{
		repository: repository,
	}
}

func (s *AssetService) Create(
	ctx context.Context,
	asset *models.Asset,
) error {

	return s.repository.Create(ctx, asset)
}

func (s *AssetService) FindAll(
	ctx context.Context,
) ([]models.Asset, error) {

	return s.repository.FindAll(ctx)
}