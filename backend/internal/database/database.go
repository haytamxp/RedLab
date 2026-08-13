package database

import "github.com/jackc/pgx/v5/pgxpool"

// DB is the global PostgreSQL connection pool.
// All repositories will use this variable.
var DB *pgxpool.Pool
