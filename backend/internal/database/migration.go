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
		role TEXT NOT NULL DEFAULT 'student',
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		ldap_user BOOLEAN NOT NULL DEFAULT FALSE,
		last_login TIMESTAMP NULL,
		manager_id UUID NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
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
		token_hash TEXT UNIQUE,
		token_issued_at TIMESTAMP NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);
	`

	if _, err := db.Exec(ctx, agentsQuery); err != nil {
		return fmt.Errorf("create agents table: %w", err)
	}

	alterAgentsQuery := `
	ALTER TABLE agents
		ADD COLUMN IF NOT EXISTS token_hash TEXT UNIQUE,
		ADD COLUMN IF NOT EXISTS token_issued_at TIMESTAMP NULL;
	`

	if _, err := db.Exec(ctx, alterAgentsQuery); err != nil {
		return fmt.Errorf("update agents authentication columns: %w", err)
	}

	tasksQuery := `
	CREATE TABLE IF NOT EXISTS tasks (
		id UUID PRIMARY KEY,
		agent_id UUID NOT NULL,
		type TEXT NOT NULL,
		payload JSONB NOT NULL DEFAULT '{}'::jsonb,
		status TEXT NOT NULL DEFAULT 'PENDING',
		priority INTEGER NOT NULL DEFAULT 0,
		result JSONB NULL,
		error_message TEXT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		claimed_at TIMESTAMP NULL,
		completed_at TIMESTAMP NULL,
		CONSTRAINT fk_tasks_agent
			FOREIGN KEY (agent_id)
			REFERENCES agents(id)
			ON DELETE CASCADE
	);
	`

	if _, err := db.Exec(ctx, tasksQuery); err != nil {
		return fmt.Errorf("create tasks table: %w", err)
	}

	taskIndexes := `
	CREATE INDEX IF NOT EXISTS idx_tasks_agent_status_priority
	ON tasks(agent_id, status, priority DESC, created_at ASC);

	CREATE INDEX IF NOT EXISTS idx_tasks_created_at
	ON tasks(created_at);
	`

	if _, err := db.Exec(ctx, taskIndexes); err != nil {
		return fmt.Errorf("create task indexes: %w", err)
	}

	return nil
}
