package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run(address string, engine *gin.Engine) error {

	server := &http.Server{
		Addr:    address,
		Handler: engine,
	}

	return server.ListenAndServe()
}