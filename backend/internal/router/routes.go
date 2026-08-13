package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/haytamxp/redlab/backend/internal/auth"
	"github.com/haytamxp/redlab/backend/internal/database"
	"github.com/haytamxp/redlab/backend/internal/handlers"
	"github.com/haytamxp/redlab/backend/internal/permissions"
	"github.com/haytamxp/redlab/backend/internal/repository"
	"github.com/haytamxp/redlab/backend/internal/services"
)

func (r *Router) RegisterRoutes(
	authHandler *handlers.AuthHandler,
	agentHandler *handlers.AgentHandler,
	jwtSecret string,
) {
	r.Engine.GET("/health", handlers.Health)

	api := r.Engine.Group("/api")
	v1 := api.Group("/v1")

	// -------------------------
	// Public routes
	// -------------------------

	v1.GET("/health", handlers.Health)

	authGroup := v1.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)

	// -------------------------
	// Task infrastructure
	// -------------------------

	agentRepository := repository.NewAgentRepository(database.DB)
	agentService := services.NewAgentService(agentRepository)

	taskRepository := repository.NewTaskRepository(database.DB)
	taskService := services.NewTaskService(
		taskRepository,
		agentService,
	)

	taskHandler := handlers.NewTaskHandler(
		taskService,
		agentService,
	)

	// Agent-authenticated endpoints.
	v1.POST(
		"/agents/:id/tasks/next",
		taskHandler.Next,
	)

	v1.GET(
		"/agents/:id/tasks",
		taskHandler.ListForAgent,
	)

	v1.POST(
		"/agents/:id/tasks/:taskId/result",
		taskHandler.Complete,
	)

	// -------------------------
	// Existing agent heartbeat
	// -------------------------

	v1.POST(
		"/agents/:id/heartbeat",
		agentHandler.Heartbeat,
	)

	// -------------------------
	// Human authenticated routes
	// -------------------------

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

	// -------------------------
	// Agent administration
	// -------------------------

	agents := protected.Group("/agents")
	agents.Use(auth.RequirePermission(permissions.ManageAgents))

	agents.POST("", agentHandler.Create)
	agents.GET("", agentHandler.List)
	agents.GET("/:id", agentHandler.Get)

	// -------------------------
	// Task administration
	// -------------------------

	tasks := protected.Group("/tasks")
	tasks.Use(auth.RequirePermission(permissions.ManageAgents))

	tasks.POST("", taskHandler.Create)
}
