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

	// Registration alone does not prove that the agent is alive.
	// The agent becomes ONLINE after sending a heartbeat.
	agent.Status = models.AgentOffline
	agent.LastSeen = nil

	return s.repository.Create(ctx, agent)
}

func (s *AgentService) List(
	ctx context.Context,
) ([]models.Agent, error) {
	return s.repository.List(ctx)
}

func (s *AgentService) Get(
	ctx context.Context,
	id uuid.UUID,
) (*models.Agent, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *AgentService) Heartbeat(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.repository.Heartbeat(ctx, id)
}
