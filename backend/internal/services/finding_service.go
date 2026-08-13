package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/haytamxp/redlab/backend/internal/models"
	"github.com/haytamxp/redlab/backend/internal/repository"
)

var ErrInvalidFindingResult = errors.New(
	"task result does not contain MITRE technique metadata",
)

type FindingService struct {
	repository *repository.FindingRepository
}

func NewFindingService(
	repository *repository.FindingRepository,
) *FindingService {
	return &FindingService{
		repository: repository,
	}
}

type taskTechnique struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type taskResultEnvelope struct {
	Module         string          `json:"module"`
	TaskType       string          `json:"task_type"`
	MitreTechnique taskTechnique   `json:"mitre_technique"`
	Data           json.RawMessage `json:"data"`
}

func (s *FindingService) CreateFromTaskResult(
	ctx context.Context,
	task *models.Task,
	result json.RawMessage,
) (*models.Finding, error) {
	var envelope taskResultEnvelope

	if err := json.Unmarshal(
		result,
		&envelope,
	); err != nil {
		return nil, fmt.Errorf(
			"decode task result: %w",
			err,
		)
	}

	if envelope.MitreTechnique.ID == "" {
		return nil, ErrInvalidFindingResult
	}

	now := time.Now()

	taskID := task.ID

	finding := &models.Finding{
		Base: models.Base{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		TaskID:        &taskID,
		AgentID:       task.AgentID,
		Title:         buildFindingTitle(envelope),
		Description:   buildFindingDescription(envelope),
		Severity:      models.SeverityInfo,
		TechniqueID:   envelope.MitreTechnique.ID,
		TechniqueName: envelope.MitreTechnique.Name,
		Evidence:      result,
		Recommendation: buildRecommendation(
			envelope.TaskType,
		),
	}

	if err := s.repository.Create(
		ctx,
		finding,
	); err != nil {
		return nil, err
	}

	return finding, nil
}

func (s *FindingService) FindAll(
	ctx context.Context,
) ([]models.Finding, error) {
	return s.repository.FindAll(ctx)
}

func (s *FindingService) Get(
	ctx context.Context,
	id uuid.UUID,
) (*models.Finding, error) {
	return s.repository.FindByID(ctx, id)
}

func buildFindingTitle(
	result taskResultEnvelope,
) string {
	if result.Module != "" {
		return fmt.Sprintf(
			"RedLab discovery: %s",
			result.Module,
		)
	}

	return fmt.Sprintf(
		"RedLab discovery: %s",
		strings.ToLower(result.TaskType),
	)
}

func buildFindingDescription(
	result taskResultEnvelope,
) string {
	return fmt.Sprintf(
		"RedLab executed %s against the authorized assessment environment. The activity maps to MITRE ATT&CK %s (%s).",
		result.TaskType,
		result.MitreTechnique.ID,
		result.MitreTechnique.Name,
	)
}

func buildRecommendation(
	taskType string,
) string {
	switch strings.ToUpper(taskType) {
	case "AD_USER_ENUMERATION":
		return "Review discovered domain accounts and validate that privileged and service accounts follow least-privilege and lifecycle controls."

	case "AD_GROUP_ENUMERATION":
		return "Review privileged and sensitive groups for unnecessary membership and excessive permissions."

	case "AD_COMPUTER_ENUMERATION":
		return "Review discovered systems and verify that asset inventory and administrative exposure are appropriately controlled."

	case "DOMAIN_INFO":
		return "Review discovered domain information and trust relationships against the intended Active Directory security architecture."

	default:
		return "Review the collected assessment evidence and validate the discovered configuration against the organization's security baseline."
	}
}
