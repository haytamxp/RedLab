package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/haytamxp/redlab/backend/internal/models"
)

type FindingRepository struct {
	db *pgxpool.Pool
}

func NewFindingRepository(db *pgxpool.Pool) *FindingRepository {
	return &FindingRepository{
		db: db,
	}
}

func (r *FindingRepository) Create(ctx context.Context, finding *models.Finding) error {
	return nil
}

func (r *FindingRepository) FindAll(ctx context.Context) ([]models.Finding, error) {
	return []models.Finding{}, nil
}