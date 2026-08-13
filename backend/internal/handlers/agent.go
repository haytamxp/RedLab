package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/haytamxp/redlab/backend/internal/dto"
	"github.com/haytamxp/redlab/backend/internal/models"
	"github.com/haytamxp/redlab/backend/internal/services"
)

type AgentHandler struct {
	service *services.AgentService
}

func NewAgentHandler(
	service *services.AgentService,
) *AgentHandler {
	return &AgentHandler{
		service: service,
	}
}

func (h *AgentHandler) Create(c *gin.Context) {
	var req dto.CreateAgentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"error":   err.Error(),
			},
		)

		return
	}

	agent := &models.Agent{
		Name:            req.Name,
		Hostname:        req.Hostname,
		IPAddress:       req.IPAddress,
		OperatingSystem: req.OperatingSystem,
		Version:         req.Version,
	}

	registrationToken, err := h.service.Create(
		c.Request.Context(),
		agent,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"error":   "failed to create agent",
			},
		)

		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"success": true,
			"message": "Agent created successfully. Store the registration token securely; it will not be shown again.",
			"data": gin.H{
				"id":                 agent.ID,
				"name":               agent.Name,
				"hostname":           agent.Hostname,
				"ip_address":         agent.IPAddress,
				"operating_system":   agent.OperatingSystem,
				"version":            agent.Version,
				"status":             agent.Status,
				"created_at":         agent.CreatedAt,
				"updated_at":         agent.UpdatedAt,
				"registration_token": registrationToken,
			},
		},
	)
}

func (h *AgentHandler) List(c *gin.Context) {
	agents, err := h.service.List(
		c.Request.Context(),
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"error":   "failed to list agents",
			},
		)

		return
	}

	response := make(
		[]dto.AgentResponse,
		0,
		len(agents),
	)

	for i := range agents {
		response = append(
			response,
			toAgentResponse(&agents[i]),
		)
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"data":    response,
		},
	)
}

func (h *AgentHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(
		c.Param("id"),
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"error":   "invalid agent ID",
			},
		)

		return
	}

	agent, err := h.service.Get(
		c.Request.Context(),
		id,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(
				http.StatusNotFound,
				gin.H{
					"success": false,
					"error":   "agent not found",
				},
			)

			return
		}

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"error":   "failed to retrieve agent",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"data":    toAgentResponse(agent),
		},
	)
}

func (h *AgentHandler) Heartbeat(c *gin.Context) {
	id, err := uuid.Parse(
		c.Param("id"),
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"error":   "invalid agent ID",
			},
		)

		return
	}

	token := extractBearerToken(
		c.GetHeader("Authorization"),
	)

	if token == "" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"success": false,
				"error":   "missing agent token",
			},
		)

		return
	}

	authenticatedAgent, err :=
		h.service.AuthenticateToken(
			c.Request.Context(),
			token,
		)

	if errors.Is(
		err,
		services.ErrInvalidAgentToken,
	) {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"success": false,
				"error":   "invalid agent token",
			},
		)

		return
	}

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"error":   "agent authentication failed",
			},
		)

		return
	}

	if authenticatedAgent.ID != id {
		c.JSON(
			http.StatusForbidden,
			gin.H{
				"success": false,
				"error":   "agent token does not match requested agent",
			},
		)

		return
	}

	if err := h.service.Heartbeat(
		c.Request.Context(),
		id,
	); err != nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"success": false,
				"error":   "agent not found",
			},
		)

		return
	}

	agent, err := h.service.Get(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"error":   "heartbeat recorded but failed to retrieve agent",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"message": "Heartbeat received",
			"data":    toAgentResponse(agent),
		},
	)
}

func extractBearerToken(header string) string {
	parts := strings.Fields(header)

	if len(parts) != 2 {
		return ""
	}

	if !strings.EqualFold(
		parts[0],
		"Bearer",
	) {
		return ""
	}

	return strings.TrimSpace(parts[1])
}

func toAgentResponse(
	agent *models.Agent,
) dto.AgentResponse {
	return dto.AgentResponse{
		ID:              agent.ID.String(),
		Name:            agent.Name,
		Hostname:        agent.Hostname,
		IPAddress:       agent.IPAddress,
		OperatingSystem: agent.OperatingSystem,
		Version:         agent.Version,
		Status:          string(agent.Status),
		LastSeen:        agent.LastSeen,
		CreatedAt:       agent.CreatedAt,
		UpdatedAt:       agent.UpdatedAt,
	}
}
