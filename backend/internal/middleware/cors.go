package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {

	return cors.New(cors.Config{
		AllowOrigins: []string{"*"},

		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},

		AllowHeaders: []string{
			"Origin",
			"Authorization",
			"Content-Type",
		},

		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
}