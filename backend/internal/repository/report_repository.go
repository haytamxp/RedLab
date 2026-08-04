package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/haytamxp/redlab/backend/internal/models"
)

type ReportRepository struct {
	db *pgxpool.Pool
}

func NewReportRepository(db *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{
		db: db,
	}
}

func (r *ReportRepository) Create(ctx context.Context, report *models.Report) error {
	return nil
}

func (r *ReportRepository) FindAll(ctx context.Context) ([]models.Report, error) {
	return []models.Report{}, nil
}