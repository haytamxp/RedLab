package main

import (
	"log"

	"github.com/haytamxp/redlab/backend/internal/app"
)

func main() {

	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := application.Start(); err != nil {
		log.Fatal(err)
	}
}