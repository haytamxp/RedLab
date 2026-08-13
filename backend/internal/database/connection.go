package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/haytamxp/redlab/backend/internal/config"
)

// Connect creates and validates the PostgreSQL connection pool
// and stores it in the global DB variable.
func Connect(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("invalid PostgreSQL configuration: %w", err)
	}

	if cfg.Database.MaxOpenConns > 0 {
		poolConfig.MaxConns = int32(cfg.Database.MaxOpenConns)
	}

	if cfg.Database.MaxIdleConns > 0 {
		poolConfig.MinConns = int32(cfg.Database.MaxIdleConns)
	}

	if cfg.Database.ConnMaxLifetime > 0 {
		poolConfig.MaxConnLifetime =
			time.Duration(cfg.Database.ConnMaxLifetime) * time.Second
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("PostgreSQL connection failed: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("PostgreSQL ping failed: %w", err)
	}

	DB = pool

	fmt.Println("Connected to PostgreSQL")

	return pool, nil
}

// Close closes the global PostgreSQL connection pool.
func Close() {
	if DB == nil {
		return
	}

	DB.Close()
	DB = nil

	fmt.Println("PostgreSQL connection closed")
}

// Pool returns the active PostgreSQL connection pool.
func Pool() (*pgxpool.Pool, error) {
	if DB == nil {
		return nil, ErrDatabaseNotConnected
	}

	return DB, nil
}
