package services

import (
	"context"

	"github.com/haytamxp/redlab/backend/internal/models"
	"github.com/haytamxp/redlab/backend/internal/repository"
)

type AssessmentService struct {
	repository *repository.AssessmentRepository
}

func NewAssessmentService(
	repository *repository.AssessmentRepository,
) *AssessmentService {

	return &AssessmentService{
		repository: repository,
	}
}

func (s *AssessmentService) Create(
	ctx context.Context,
	assessment *models.Assessment,
) error {

	return s.repository.Create(ctx, assessment)
}

func (s *AssessmentService) FindAll(
	ctx context.Context,
) ([]models.Assessment, error) {

	return s.repository.FindAll(ctx)
}