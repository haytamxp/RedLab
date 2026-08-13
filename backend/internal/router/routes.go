package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/haytamxp/redlab/backend/internal/auth"
	"github.com/haytamxp/redlab/backend/internal/handlers"
	"github.com/haytamxp/redlab/backend/internal/permissions"
)

func (r *Router) RegisterRoutes(
	authHandler *handlers.AuthHandler,
	agentHandler *handlers.AgentHandler,
	jwtSecret string,
) {
	r.Engine.GET("/health", handlers.Health)

	api := r.Engine.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/health", handlers.Health)

	authGroup := v1.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)

	protected := v1.Group("/")
	protected.Use(auth.JWTMiddleware(jwtSecret))

	protected.GET("/profile", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Authenticated",
			"user_id": c.GetString("userID"),
			"role":    c.GetString("role"),
		})
	})

	agents := protected.Group("/agents")
	agents.Use(auth.RequirePermission(permissions.ManageAgents))

	agents.POST("", agentHandler.Create)
	agents.GET("", agentHandler.List)
	agents.GET("/:id", agentHandler.Get)
	agents.POST("/:id/heartbeat", agentHandler.Heartbeat)
}
