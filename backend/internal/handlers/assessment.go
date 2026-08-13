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

type AssessmentHandler struct {
	service *services.AssessmentService
}

func NewAssessmentHandler(
	service *services.AssessmentService,
) *AssessmentHandler {
	return &AssessmentHandler{
		service: service,
	}
}

func (h *AssessmentHandler) Create(
	c *gin.Context,
) {
	var req dto.CreateAssessmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	agentID, err := uuid.Parse(
		req.AgentID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid agent_id",
		})
		return
	}

	assessment, err := h.service.Create(
		c.Request.Context(),
		req.Name,
		req.Description,
		agentID,
	)

	if errors.Is(
		err,
		repository.ErrAgentNotFound,
	) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "agent not found",
		})
		return
	}

	if errors.Is(
		err,
		services.ErrAssessmentNameRequired,
	) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to create assessment",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Assessment created successfully",
		"data":    toAssessmentResponse(assessment),
	})
}

func (h *AssessmentHandler) List(
	c *gin.Context,
) {
	assessments, err := h.service.FindAll(
		c.Request.Context(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to retrieve assessments",
		})
		return
	}

	response := make(
		[]dto.AssessmentResponse,
		0,
		len(assessments),
	)

	for i := range assessments {
		response = append(
			response,
			toAssessmentResponse(
				&assessments[i],
			),
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

func (h *AssessmentHandler) Get(
	c *gin.Context,
) {
	id, err := uuid.Parse(
		c.Param("id"),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid assessment ID",
		})
		return
	}

	assessment, err := h.service.FindByID(
		c.Request.Context(),
		id,
	)

	if errors.Is(
		err,
		repository.ErrAssessmentNotFound,
	) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "assessment not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to retrieve assessment",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    toAssessmentResponse(assessment),
	})
}

func (h *AssessmentHandler) UpdateStatus(
	c *gin.Context,
) {
	id, err := uuid.Parse(
		c.Param("id"),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid assessment ID",
		})
		return
	}

	var request struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(
		&request,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	err = h.service.UpdateStatus(
		c.Request.Context(),
		id,
		request.Status,
	)

	if errors.Is(
		err,
		services.ErrInvalidAssessmentStatus,
	) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if errors.Is(
		err,
		repository.ErrAssessmentNotFound,
	) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "assessment not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to update assessment",
		})
		return
	}

	assessment, err := h.service.FindByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "status updated but assessment retrieval failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Assessment status updated",
		"data":    toAssessmentResponse(assessment),
	})
}

func toAssessmentResponse(
	assessment *models.Assessment,
) dto.AssessmentResponse {
	return dto.AssessmentResponse{
		ID:          assessment.ID.String(),
		Name:        assessment.Name,
		Description: assessment.Description,
		AgentID:     assessment.AgentID.String(),
		Status:      string(assessment.Status),
		CreatedAt:   assessment.CreatedAt,
		UpdatedAt:   assessment.UpdatedAt,
	}
}
