package reporting

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(
	service *Service,
) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Get(
	c *gin.Context,
) {
	assessmentID, err := uuid.Parse(
		c.Param("id"),
	)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"error":   "invalid assessment ID",
			},
		)
		return
	}

	report, err := h.service.Build(
		c.Request.Context(),
		assessmentID,
	)

	if errors.Is(
		err,
		ErrReportAssessmentNotFound,
	) {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"success": false,
				"error":   "assessment not found",
			},
		)
		return
	}

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"error":   "failed to build report",
			},
		)
		return
	}

	c.Header(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	if err := RenderReport(
		c.Writer,
		report,
	); err != nil {
		if !c.IsAborted() {
			c.AbortWithStatus(
				http.StatusInternalServerError,
			)
		}
		return
	}
}
