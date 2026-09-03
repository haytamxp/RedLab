package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/haytamxp/redlab/backend/internal/models"
	"github.com/haytamxp/redlab/backend/internal/repository"
)

var ErrInvalidAgentToken = errors.New("invalid agent token")

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
) (string, error) {
	now := time.Now()

	agent.ID = uuid.New()
	agent.CreatedAt = now
	agent.UpdatedAt = now

	agent.Status = models.AgentOffline
	agent.LastSeen = nil

	token, err := generateAgentToken()
	if err != nil {
		return "", err
	}

	tokenHash := hashAgentToken(token)

	if err := s.repository.Create(
		ctx,
		agent,
		tokenHash,
		now,
	); err != nil {
		return "", err
	}

	return token, nil
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

func (s *AgentService) AuthenticateToken(
	ctx context.Context,
	token string,
) (*models.Agent, error) {
	if token == "" {
		return nil, ErrInvalidAgentToken
	}

	tokenHash := hashAgentToken(token)

	agent, err := s.repository.FindByTokenHash(
		ctx,
		tokenHash,
	)
	if err != nil {
		if errors.Is(err, repository.ErrAgentNotFound) {
			return nil, ErrInvalidAgentToken
		}

		return nil, err
	}

	return agent, nil
}

func generateAgentToken() (string, error) {
	raw := make([]byte, 32)

	if _, err := rand.Read(raw); err != nil {
		return "", err
	}

	return hex.EncodeToString(raw), nil
}

func hashAgentToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
