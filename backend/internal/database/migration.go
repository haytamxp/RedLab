package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(db *pgxpool.Pool) error {
	ctx := context.Background()

	usersQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		first_name TEXT,
		last_name TEXT,
		role TEXT,
		is_active BOOLEAN DEFAULT TRUE,
		ldap_user BOOLEAN DEFAULT FALSE,
		last_login TIMESTAMP,
		manager_id UUID,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);
	`

	if _, err := db.Exec(ctx, usersQuery); err != nil {
		return fmt.Errorf("create users table: %w", err)
	}

	agentsQuery := `
	CREATE TABLE IF NOT EXISTS agents (
		id UUID PRIMARY KEY,
		name TEXT NOT NULL,
		hostname TEXT NOT NULL,
		ip_address TEXT NOT NULL,
		operating_system TEXT,
		version TEXT,
		status TEXT NOT NULL,
		last_seen TIMESTAMP NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);
	`

	if _, err := db.Exec(ctx, agentsQuery); err != nil {
		return fmt.Errorf("create agents table: %w", err)
	}

	return nil
}
