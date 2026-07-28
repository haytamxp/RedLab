package middleware

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/haytamxp/redlab/backend/internal/logger"
)

func Logger() gin.HandlerFunc {

	return func(c *gin.Context) {

		start := time.Now()

		c.Next()

		logger.Log.Info(
			"HTTP Request",
		)
		_ = start
	}
}