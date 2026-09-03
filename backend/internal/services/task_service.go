package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/haytamxp/redlab/backend/internal/models"
	"github.com/haytamxp/redlab/backend/internal/repository"
)

var (
	ErrTaskTypeRequired = errors.New(
		"task type is required",
	)

	ErrInvalidTaskState = errors.New(
		"invalid task state",
	)

	ErrInvalidTaskResult = errors.New(
		"task result status must be COMPLETED or FAILED",
	)
)

type TaskService struct {
	repository   *repository.TaskRepository
	agentService *AgentService
}

func NewTaskService(
	repository *repository.TaskRepository,
	agentService *AgentService,
) *TaskService {
	return &TaskService{
		repository:   repository,
		agentService: agentService,
	}
}


func (s *TaskService) Create(
	ctx context.Context,
	agentID uuid.UUID,
	taskType string,
	payload json.RawMessage,
	priority int,
) (*models.Task, error) {
	taskType =
		strings.TrimSpace(
			taskType,
		)

	if taskType == "" {
		return nil,
			ErrTaskTypeRequired
	}

	if len(payload) == 0 {
		payload =
			json.RawMessage(`{}`)
	}

	if _, err :=
		s.agentService.Get(
			ctx,
			agentID,
		); err != nil {
		return nil, err
	}

	now := time.Now()

	task := &models.Task{
		Base: models.Base{
			ID: uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},

		AgentID: agentID,

		Type: taskType,

		Payload: payload,

		Status:
			models.TaskPending,

		Priority: priority,
	}

	if err :=
		s.repository.Create(
			ctx,
			task,
		); err != nil {
		return nil, err
	}

	return task, nil
}


func (s *TaskService) Next(
	ctx context.Context,
	agentID uuid.UUID,
) (*models.Task, error) {
	return s.repository.ClaimNext(
		ctx,
		agentID,
	)
}


func (s *TaskService) Complete(
	ctx context.Context,
	taskID uuid.UUID,
	agentID uuid.UUID,
	status string,
	result json.RawMessage,
	errorMessage string,
) (*models.Task, error) {
	var taskStatus models.TaskStatus

	switch strings.ToUpper(
		strings.TrimSpace(status),
	) {

	case string(models.TaskCompleted):
		taskStatus =
			models.TaskCompleted

	case string(models.TaskFailed):
		taskStatus =
			models.TaskFailed

	default:
		return nil,
			ErrInvalidTaskResult
	}

	if len(result) == 0 {
		result =
			json.RawMessage(`{}`)
	}

	var taskError *string

	if strings.TrimSpace(
		errorMessage,
	) != "" {
		message :=
			strings.TrimSpace(
				errorMessage,
			)

		taskError = &message
	}

	return s.repository.Complete(
		ctx,
		taskID,
		agentID,
		taskStatus,
		result,
		taskError,
	)
}


func (s *TaskService) ListForAgent(
	ctx context.Context,
	agentID uuid.UUID,
) ([]models.Task, error) {
	return s.repository.ListForAgent(
		ctx,
		agentID,
	)
}


/*
 * Operator task management.
 */

func (s *TaskService) ListAll(
	ctx context.Context,
) ([]models.Task, error) {
	return s.repository.ListAll(
		ctx,
	)
}


func (s *TaskService) DeletePending(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.repository.DeletePending(
		ctx,
		id,
	)
}