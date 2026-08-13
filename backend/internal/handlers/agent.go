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

type AgentHandler struct {
	service *services.AgentService
}

func NewAgentHandler(service *services.AgentService) *AgentHandler {
	return &AgentHandler{
		service: service,
	}
}

func (h *AgentHandler) Create(c *gin.Context) {
	var req dto.CreateAgentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	agent := &models.Agent{
		Name:            req.Name,
		Hostname:        req.Hostname,
		IPAddress:       req.IPAddress,
		OperatingSystem: req.OperatingSystem,
		Version:         req.Version,
	}

	if err := h.service.Create(
		c.Request.Context(),
		agent,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to create agent",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Agent created successfully",
		"data":    agent,
	})
}

func (h *AgentHandler) List(c *gin.Context) {
	agents, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to list agents",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    agents,
	})
}

func (h *AgentHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid agent ID",
		})
		return
	}

	agent, err := h.service.Get(
		c.Request.Context(),
		id,
	)

	if errors.Is(err, repository.ErrAgentNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "agent not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to retrieve agent",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    agent,
	})
}

func (h *AgentHandler) Heartbeat(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid agent ID",
		})
		return
	}

	if err := h.service.Heartbeat(
		c.Request.Context(),
		id,
	); errors.Is(err, repository.ErrAgentNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "agent not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to record heartbeat",
		})
		return
	}

	agent, err := h.service.Get(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "heartbeat recorded but agent retrieval failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Heartbeat received",
		"data":    agent,
	})
}
