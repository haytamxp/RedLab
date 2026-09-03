package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/haytamxp/redlab/backend/internal/models"
)

var ErrAgentNotFound = errors.New("agent not found")

type AgentRepository struct {
	db *pgxpool.Pool
}

func NewAgentRepository(db *pgxpool.Pool) *AgentRepository {
	return &AgentRepository{
		db: db,
	}
}

func (r *AgentRepository) Create(
	ctx context.Context,
	agent *models.Agent,
	tokenHash string,
	tokenIssuedAt time.Time,
) error {
	query := `
INSERT INTO agents (
id,
name,
hostname,
ip_address,
operating_system,
version,
status,
last_seen,
token_hash,
token_issued_at,
created_at,
updated_at
)
VALUES (
$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12
)
`

	_, err := r.db.Exec(
		ctx,
		query,
		agent.ID,
		agent.Name,
		agent.Hostname,
		agent.IPAddress,
		agent.OperatingSystem,
		agent.Version,
		agent.Status,
		agent.LastSeen,
		tokenHash,
		tokenIssuedAt,
		agent.CreatedAt,
		agent.UpdatedAt,
	)

	return err
}

func (r *AgentRepository) List(
	ctx context.Context,
) ([]models.Agent, error) {
	query := `
SELECT
id,
name,
hostname,
ip_address,
operating_system,
version,
status,
last_seen,
created_at,
updated_at
FROM agents
ORDER BY created_at DESC
`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	agents := make([]models.Agent, 0)

	for rows.Next() {
		var agent models.Agent

		if err := rows.Scan(
			&agent.ID,
			&agent.Name,
			&agent.Hostname,
			&agent.IPAddress,
			&agent.OperatingSystem,
			&agent.Version,
			&agent.Status,
			&agent.LastSeen,
			&agent.CreatedAt,
			&agent.UpdatedAt,
		); err != nil {
			return nil, err
		}

		agents = append(agents, agent)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return agents, nil
}

func (r *AgentRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Agent, error) {
	query := `
SELECT
id,
name,
hostname,
ip_address,
operating_system,
version,
status,
last_seen,
created_at,
updated_at
FROM agents
WHERE id = $1
`

	var agent models.Agent

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&agent.ID,
		&agent.Name,
		&agent.Hostname,
		&agent.IPAddress,
		&agent.OperatingSystem,
		&agent.Version,
		&agent.Status,
		&agent.LastSeen,
		&agent.CreatedAt,
		&agent.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrAgentNotFound
	}

	if err != nil {
		return nil, err
	}

	return &agent, nil
}

func (r *AgentRepository) FindByTokenHash(
	ctx context.Context,
	tokenHash string,
) (*models.Agent, error) {
	query := `
SELECT
id,
name,
hostname,
ip_address,
operating_system,
version,
status,
last_seen,
created_at,
updated_at
FROM agents
WHERE token_hash = $1
  AND deleted_at IS NULL
`

	var agent models.Agent

	err := r.db.QueryRow(
		ctx,
		query,
		tokenHash,
	).Scan(
		&agent.ID,
		&agent.Name,
		&agent.Hostname,
		&agent.IPAddress,
		&agent.OperatingSystem,
		&agent.Version,
		&agent.Status,
		&agent.LastSeen,
		&agent.CreatedAt,
		&agent.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrAgentNotFound
	}

	if err != nil {
		return nil, err
	}

	return &agent, nil
}

func (r *AgentRepository) Heartbeat(
	ctx context.Context,
	id uuid.UUID,
) error {
	query := `
UPDATE agents
SET
status = $1,
last_seen = NOW(),
updated_at = NOW()
WHERE id = $2
`

	result, err := r.db.Exec(
		ctx,
		query,
		models.AgentOnline,
		id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrAgentNotFound
	}

	return nil
}
