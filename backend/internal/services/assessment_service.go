package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/haytamxp/redlab/backend/internal/models"
	"github.com/haytamxp/redlab/backend/internal/repository"
)

var (
	ErrAssessmentNameRequired = errors.New(
		"assessment name is required",
	)
	ErrInvalidAssessmentStatus = errors.New(
		"invalid assessment status",
	)
)

type AssessmentService struct {
	repository   *repository.AssessmentRepository
	agentService *AgentService
}

func NewAssessmentService(
	repository *repository.AssessmentRepository,
	agentService *AgentService,
) *AssessmentService {
	return &AssessmentService{
		repository:   repository,
		agentService: agentService,
	}
}

func (s *AssessmentService) Create(
	ctx context.Context,
	name string,
	description string,
	agentID uuid.UUID,
) (*models.Assessment, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, ErrAssessmentNameRequired
	}

	if _, err := s.agentService.Get(
		ctx,
		agentID,
	); err != nil {
		return nil, err
	}

	now := time.Now()

	assessment := &models.Assessment{
		Base: models.Base{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:        name,
		Description: strings.TrimSpace(description),
		AgentID:     agentID,
		Status:      models.AssessmentPending,
	}

	if err := s.repository.Create(
		ctx,
		assessment,
	); err != nil {
		return nil, err
	}

	return assessment, nil
}

func (s *AssessmentService) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Assessment, error) {
	return s.repository.FindByID(
		ctx,
		id,
	)
}

func (s *AssessmentService) FindAll(
	ctx context.Context,
) ([]models.Assessment, error) {
	return s.repository.FindAll(ctx)
}

func (s *AssessmentService) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	status string,
) error {
	switch models.AssessmentStatus(
		strings.ToUpper(
			strings.TrimSpace(status),
		),
	) {
	case models.AssessmentPending,
		models.AssessmentRunning,
		models.AssessmentCompleted,
		models.AssessmentFailed:

		return s.repository.UpdateStatus(
			ctx,
			id,
			models.AssessmentStatus(
				strings.ToUpper(
					strings.TrimSpace(status),
				),
			),
		)

	default:
		return ErrInvalidAssessmentStatus
	}
}
