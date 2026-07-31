package database

import (
	"context"
	"time"

	appErrors "github.com/haytamxp/redlab/backend/internal/errors"
)

func Health() error {

	if DB == nil {
		return appErrors.ErrDatabaseNotConnected
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return DB.Ping(ctx)
}