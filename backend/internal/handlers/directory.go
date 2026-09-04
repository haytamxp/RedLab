package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/haytamxp/redlab/backend/internal/dto"
	"github.com/haytamxp/redlab/backend/internal/services"
)

type DirectoryHandler struct {
	service *services.LDAPService
}

func NewDirectoryHandler(service *services.LDAPService) *DirectoryHandler {
	return &DirectoryHandler{service: service}
}

func (h *DirectoryHandler) ListUsers(c *gin.Context) {
	users, err := h.service.ListUsers(
		c.Request.Context(),
		c.Query("search"),
	)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"success": false,
			"error":   "Active Directory search failed",
		})
		return
	}

	response := make([]dto.DirectoryUserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, dto.DirectoryUserResponse{
			DN:                user.DN,
			Username:          user.SAMAccountName,
			UserPrincipalName: user.UserPrincipalName,
			Email:             user.Email,
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			DisplayName:       user.DisplayName,
			Enabled:           user.Enabled,
			Groups:            user.Groups,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}
