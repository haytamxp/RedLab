package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/haytamxp/redlab/backend/internal/dto"
	"github.com/haytamxp/redlab/backend/internal/models"
	"github.com/haytamxp/redlab/backend/internal/services"
)

type AuthHandler struct {
	auth *services.AuthService
}

func NewAuthHandler(
	auth *services.AuthService,
) *AuthHandler {
	return &AuthHandler{
		auth: auth,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

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

	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      models.RoleStudent,
		LDAPUser:  false,
		IsActive:  true,
	}

	if err := h.auth.Register(
		c.Request.Context(),
		user,
		req.Password,
	); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"error":   err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusCreated,
		dto.RegisterResponse{
			Message: "User registered successfully",
		},
	)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

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

	token, err := h.auth.Login(
		c.Request.Context(),
		req.Username,
		req.Password,
	)

	if errors.Is(
		err,
		services.ErrInvalidCredentials,
	) ||
		errors.Is(
			err,
			services.ErrUserNotFound,
		) {

		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"success": false,
				"error":   "invalid username or password",
			},
		)

		return
	}

	if errors.Is(
		err,
		services.ErrUserInactive,
	) {
		c.JSON(
			http.StatusForbidden,
			gin.H{
				"success": false,
				"error":   "user account is inactive",
			},
		)

		return
	}

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"error":   "authentication failed",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		dto.LoginResponse{
			Token: token,
		},
	)
}
