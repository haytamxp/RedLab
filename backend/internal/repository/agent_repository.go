package repository

import (
	"context"

	"github.com/haytamxp/redlab/backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
) error {

	query := `
	INSERT INTO agents
	(
		id,
		name,
		hostname,
		ip_address,
		operating_system,
		version,
		status,
		created_at,
		updated_at
	)
	VALUES
	($1,$2,$3,$4,$5,$6,$7,$8,$9)
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
		agent.CreatedAt,
		agent.UpdatedAt,
	)

	return err
}
