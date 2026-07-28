package logger

import (
	"go.uber.org/zap"
)

func New(development bool) (*zap.Logger, error) {

	if development {
		return zap.NewDevelopment()
	}

	return zap.NewProduction()
}