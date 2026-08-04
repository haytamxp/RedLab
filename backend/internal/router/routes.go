package router

import (
	"github.com/gin-gonic/gin"

	"github.com/haytamxp/redlab/backend/internal/auth"
	"github.com/haytamxp/redlab/backend/internal/handlers"
)

func (r *Router) RegisterRoutes(
	authHandler *handlers.AuthHandler,
) {

	r.Engine.GET("/health", handlers.Health)

	api := r.Engine.Group("/api")

	v1 := api.Group("/v1")

	{
		v1.GET("/health", handlers.Health)

		authGroup := v1.Group("/auth")

		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}

		protected := v1.Group("/")

		protected.Use(auth.JWTMiddleware("my_secret_key"))

		{
			protected.GET("/profile", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Authenticated",
				})
			})
		}
	}
	func (r *Router) RegisterAgentRoutes(agent *handlers.AgentHandler) {

	api := r.Engine.Group("/api/v1")
	api.GET("/users", userHandler.GetAll)
	api.POST("/agents", agent.Create)
}
}