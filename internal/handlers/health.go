package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/haytamxp/redlab/backend/internal/response"
)

func Health(c *gin.Context) {

	data := gin.H{
		"status":  "running",
		"service": "RedLab Backend",
		"version": "1.0.0",
	}

	response.Success(c, "Backend is healthy", data)
}