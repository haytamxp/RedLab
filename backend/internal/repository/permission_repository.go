package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/haytamxp/redlab/backend/internal/models"
)

type PermissionRepository struct {
	db *pgxpool.Pool
}

func NewPermissionRepository(db *pgxpool.Pool) *PermissionRepository {
	return &PermissionRepository{
		db: db,
	}
}

func (r *PermissionRepository) FindAll(ctx context.Context) ([]models.Permission, error) {
	return []models.Permission{}, nil
}