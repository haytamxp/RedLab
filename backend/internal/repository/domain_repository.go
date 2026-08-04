package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/haytamxp/redlab/backend/internal/models"
)

type DomainRepository struct {
	db *pgxpool.Pool
}

func NewDomainRepository(db *pgxpool.Pool) *DomainRepository {
	return &DomainRepository{
		db: db,
	}
}

func (r *DomainRepository) Create(ctx context.Context, domain *models.Domain) error {
	return nil
}

func (r *DomainRepository) FindAll(ctx context.Context) ([]models.Domain, error) {
	return []models.Domain{}, nil
}