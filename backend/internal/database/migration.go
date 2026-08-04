package database

import (
	"context"
)

func Migrate(db *Database) error {

	query := `
CREATE TABLE IF NOT EXISTS users(

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

	_, err := db.Pool.Exec(
		context.Background(),
		query,
	)

	return err
}