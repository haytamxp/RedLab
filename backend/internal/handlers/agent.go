package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/haytamxp/redlab/backend/internal/dto"
	"github.com/haytamxp/redlab/backend/internal/models"
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
			"error": err.Error(),
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

	if err := h.service.Create(c.Request.Context(), agent); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Agent created successfully",
	})
}