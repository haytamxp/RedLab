package router

import (
	"github.com/haytamxp/redlab/backend/internal/handlers"
)

func (r *Router) registerRoutes() {

	r.Engine.GET("/health", handlers.Health)

	api := r.Engine.Group("/api")

	v1 := api.Group("/v1")

	{
		v1.GET("/health", handlers.Health)

		auth := v1.Group("/auth")

		{
			auth.POST("/login", handlers.Login)
			auth.POST("/register", handlers.Register)
		}
	}
}