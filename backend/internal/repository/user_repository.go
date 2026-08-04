package repository

import (
	"context"

	"github.com/haytamxp/redlab/backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create inserts a new user into the database.
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {

	query := `
	INSERT INTO users
	(id, username, email, password_hash, first_name, last_name, role, is_active)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Role,
		user.IsActive,
	)

	return err
}

// FindByUsername returns one user.
func (r *UserRepository) FindByUsername(
	ctx context.Context,
	username string,
) (*models.User, error) {

	user := &models.User{}

	query := `
	SELECT
	id,
	username,
	email,
	password_hash,
	first_name,
	last_name,
	role,
	is_active
	FROM users
	WHERE username=$1
	`

	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}