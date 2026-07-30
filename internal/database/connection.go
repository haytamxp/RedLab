package database

import (
	"context"
	"fmt"
	"time"

	"github.com/haytamxp/redlab/backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is the global connection pool.
// All repositories will use this pool.
var DB *pgxpool.Pool

// Connect establishes a connection pool to PostgreSQL.
func Connect(cfg *config.Config) error {

	// Build the PostgreSQL connection string (DSN).
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	// Prevent waiting forever if PostgreSQL is unavailable.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create the connection pool.
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return err
	}

	// Verify the connection.
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return err
	}

	DB = pool

	fmt.Println("✅ Connected to PostgreSQL")

	return nil
}

// Close gracefully closes the connection pool.
func Close() {

	if DB != nil {

		DB.Close()

		fmt.Println("🛑 PostgreSQL connection closed")
	}
}