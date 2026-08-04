package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/haytamxp/redlab/backend/internal/models"
)

type AssetRepository struct {
	db *pgxpool.Pool
}

func NewAssetRepository(db *pgxpool.Pool) *AssetRepository {
	return &AssetRepository{
		db: db,
	}
}

func (r *AssetRepository) Create(ctx context.Context, asset *models.Asset) error {
	return nil
}

func (r *AssetRepository) FindAll(ctx context.Context) ([]models.Asset, error) {
	return []models.Asset{}, nil
}