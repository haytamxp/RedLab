package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/haytamxp/redlab/backend/internal/models"
)

type AssessmentRepository struct {
	db *pgxpool.Pool
}

func NewAssessmentRepository(db *pgxpool.Pool) *AssessmentRepository {
	return &AssessmentRepository{
		db: db,
	}
}

func (r *AssessmentRepository) Create(ctx context.Context, assessment *models.Assessment) error {
	return nil
}

func (r *AssessmentRepository) FindAll(ctx context.Context) ([]models.Assessment, error) {
	return []models.Assessment{}, nil
}