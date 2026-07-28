package router

import (
	"github.com/gin-gonic/gin"

	"github.com/haytamxp/redlab/backend/internal/middleware"
)

type Router struct {
	Engine *gin.Engine
}

func New() *Router {

	engine := gin.New()

	engine.Use(middleware.Logger())
	engine.Use(middleware.Recovery())
	engine.Use(middleware.CORS())

	router := &Router{
		Engine: engine,
	}

	router.registerRoutes()

	return router
}