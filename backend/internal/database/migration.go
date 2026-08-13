package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(db *pgxpool.Pool) error {
	query := `
CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY,
	username TEXT UNIQUE NOT NULL,
	email TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	first_name TEXT,
	last_name TEXT,
	role TEXT,
	is_active BOOLEAN,
	ldap_user BOOLEAN,
	last_login TIMESTAMP,
	manager_id UUID,
	created_at TIMESTAMP,
	updated_at TIMESTAMP
);
`

	_, err := db.Exec(context.Background(), query)

	return err
}
