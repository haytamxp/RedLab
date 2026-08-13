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

type FindingHandler struct {
	service *services.FindingService
}

func NewFindingHandler(
	service *services.FindingService,
) *FindingHandler {
	return &FindingHandler{
		service: service,
	}
}

func (h *FindingHandler) List(c *gin.Context) {
	findings, err := h.service.FindAll(
		c.Request.Context(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to retrieve findings",
		})
		return
	}

	response := make(
		[]dto.FindingResponse,
		0,
		len(findings),
	)

	for i := range findings {
		response = append(
			response,
			toFindingResponse(
				&findings[i],
			),
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

func (h *FindingHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(
		c.Param("id"),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid finding ID",
		})
		return
	}

	finding, err := h.service.Get(
		c.Request.Context(),
		id,
	)

	if errors.Is(
		err,
		repository.ErrFindingNotFound,
	) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "finding not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to retrieve finding",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    toFindingResponse(finding),
	})
}

func toFindingResponse(
	finding *models.Finding,
) dto.FindingResponse {
	var assessmentID *string

	if finding.AssessmentID != nil {
		value := finding.AssessmentID.String()
		assessmentID = &value
	}

	var taskID *string

	if finding.TaskID != nil {
		value := finding.TaskID.String()
		taskID = &value
	}

	return dto.FindingResponse{
		ID:             finding.ID.String(),
		AssessmentID:   assessmentID,
		TaskID:         taskID,
		AgentID:        finding.AgentID.String(),
		Title:          finding.Title,
		Description:    finding.Description,
		Severity:       string(finding.Severity),
		TechniqueID:    finding.TechniqueID,
		TechniqueName:  finding.TechniqueName,
		Evidence:       finding.Evidence,
		Recommendation: finding.Recommendation,
		CreatedAt:      finding.CreatedAt,
		UpdatedAt:      finding.UpdatedAt,
	}
}
