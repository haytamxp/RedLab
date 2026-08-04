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

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {

	query := `
	INSERT INTO users
	(
		id,
		username,
		email,
		password_hash,
		first_name,
		last_name,
		role,
		is_active,
		ldap_user
	)
	VALUES
	($1,$2,$3,$4,$5,$6,$7,$8,$9)
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
		user.LDAPUser,
	)

	return err
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {

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
	is_active,
	ldap_user,
	last_login,
	manager_id,
	created_at,
	updated_at
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
		&user.LDAPUser,
		&user.LastLogin,
		&user.ManagerID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}