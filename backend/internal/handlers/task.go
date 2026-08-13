package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/haytamxp/redlab/backend/internal/dto"
	"github.com/haytamxp/redlab/backend/internal/models"
	"github.com/haytamxp/redlab/backend/internal/repository"
	"github.com/haytamxp/redlab/backend/internal/services"
)

type TaskHandler struct {
	service      *services.TaskService
	agentService *services.AgentService
}

func NewTaskHandler(
	service *services.TaskService,
	agentService *services.AgentService,
) *TaskHandler {
	return &TaskHandler{
		service:      service,
		agentService: agentService,
	}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req dto.CreateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	agentID, err := uuid.Parse(req.AgentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid agent_id",
		})
		return
	}

	task, err := h.service.Create(
		c.Request.Context(),
		agentID,
		req.Type,
		req.Payload,
		req.Priority,
	)

	if err != nil {
		if errors.Is(err, repository.ErrAgentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "agent not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to create task",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Task created successfully",
		"data":    toTaskResponse(task),
	})
}

func (h *TaskHandler) Next(c *gin.Context) {
	agentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid agent ID",
		})
		return
	}

	agent, ok := h.authenticateAgent(c)
	if !ok {
		return
	}

	if agent.ID != agentID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "agent token does not match requested agent",
		})
		return
	}

	task, err := h.service.Next(
		c.Request.Context(),
		agentID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to retrieve next task",
		})
		return
	}

	if task == nil {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    toTaskResponse(task),
	})
}

func (h *TaskHandler) Complete(c *gin.Context) {
	agentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid agent ID",
		})
		return
	}

	taskID, err := uuid.Parse(c.Param("taskId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid task ID",
		})
		return
	}

	agent, ok := h.authenticateAgent(c)
	if !ok {
		return
	}

	if agent.ID != agentID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "agent token does not match requested agent",
		})
		return
	}

	var req dto.TaskResultRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	task, err := h.service.Complete(
		c.Request.Context(),
		taskID,
		agentID,
		req.Status,
		req.Result,
		req.Error,
	)

	if errors.Is(err, services.ErrInvalidTaskResult) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if errors.Is(err, repository.ErrTaskNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "task not found or task is not claimed by this agent",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to complete task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Task result recorded",
		"data":    toTaskResponse(task),
	})
}

func (h *TaskHandler) ListForAgent(c *gin.Context) {
	agentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid agent ID",
		})
		return
	}

	agent, ok := h.authenticateAgent(c)
	if !ok {
		return
	}

	if agent.ID != agentID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "agent token does not match requested agent",
		})
		return
	}

	tasks, err := h.service.ListForAgent(
		c.Request.Context(),
		agentID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to retrieve tasks",
		})
		return
	}

	response := make([]dto.TaskResponse, 0, len(tasks))

	for i := range tasks {
		response = append(
			response,
			toTaskResponse(&tasks[i]),
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

func (h *TaskHandler) authenticateAgent(
	c *gin.Context,
) (*models.Agent, bool) {
	token := extractBearerToken(
		c.GetHeader("Authorization"),
	)

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "missing agent token",
		})
		return nil, false
	}

	agent, err := h.agentService.AuthenticateToken(
		c.Request.Context(),
		token,
	)

	if errors.Is(err, services.ErrInvalidAgentToken) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "invalid agent token",
		})
		return nil, false
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "agent authentication failed",
		})
		return nil, false
	}

	return agent, true
}

func toTaskResponse(task *models.Task) dto.TaskResponse {
	response := dto.TaskResponse{
		ID:           task.ID.String(),
		AgentID:      task.AgentID.String(),
		Type:         task.Type,
		Payload:      task.Payload,
		Status:       string(task.Status),
		Priority:     task.Priority,
		ErrorMessage: task.ErrorMessage,
		CreatedAt:    task.CreatedAt,
		UpdatedAt:    task.UpdatedAt,
		ClaimedAt:    task.ClaimedAt,
		CompletedAt:  task.CompletedAt,
	}

	if task.Result != nil {
		response.Result = *task.Result
	}

	return response
}
