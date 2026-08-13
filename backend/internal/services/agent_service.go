package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/haytamxp/redlab/backend/internal/models"
	"github.com/haytamxp/redlab/backend/internal/repository"
)

type AgentService struct {
	repository *repository.AgentRepository
}

func NewAgentService(
	repository *repository.AgentRepository,
) *AgentService {

	return &AgentService{
		repository: repository,
	}
}

func (s *AgentService) Create(
	ctx context.Context,
	agent *models.Agent,
) error {

	now := time.Now()

	agent.ID = uuid.New()

	agent.CreatedAt = now
	agent.UpdatedAt = now

	agent.Status = models.AgentOnline

	return s.repository.Create(ctx, agent)
}
