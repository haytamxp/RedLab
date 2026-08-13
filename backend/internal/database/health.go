package database

import (
	"context"
	"errors"
	"time"
)

var ErrDatabaseNotConnected = errors.New("database not connected")

func Health() error {
	if DB == nil {
		return ErrDatabaseNotConnected
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	return DB.Ping(ctx)
}
