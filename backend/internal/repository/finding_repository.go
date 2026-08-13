package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/haytamxp/redlab/backend/internal/models"
)

var ErrFindingNotFound = errors.New("finding not found")

type FindingRepository struct {
	db *pgxpool.Pool
}

func NewFindingRepository(
	db *pgxpool.Pool,
) *FindingRepository {
	return &FindingRepository{
		db: db,
	}
}

func (r *FindingRepository) Create(
	ctx context.Context,
	finding *models.Finding,
) error {
	query := `
	INSERT INTO findings (
		id,
		assessment_id,
		task_id,
		agent_id,
		title,
		description,
		severity,
		technique_id,
		technique_name,
		evidence,
		recommendation,
		created_at,
		updated_at
	)
	VALUES (
		$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13
	)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		finding.ID,
		finding.AssessmentID,
		finding.TaskID,
		finding.AgentID,
		finding.Title,
		finding.Description,
		finding.Severity,
		finding.TechniqueID,
		finding.TechniqueName,
		finding.Evidence,
		finding.Recommendation,
		finding.CreatedAt,
		finding.UpdatedAt,
	)

	return err
}

func (r *FindingRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Finding, error) {
	query := `
	SELECT
		id,
		assessment_id,
		task_id,
		agent_id,
		title,
		description,
		severity,
		technique_id,
		technique_name,
		evidence,
		recommendation,
		created_at,
		updated_at
	FROM findings
	WHERE id = $1
	`

	var finding models.Finding

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&finding.ID,
		&finding.AssessmentID,
		&finding.TaskID,
		&finding.AgentID,
		&finding.Title,
		&finding.Description,
		&finding.Severity,
		&finding.TechniqueID,
		&finding.TechniqueName,
		&finding.Evidence,
		&finding.Recommendation,
		&finding.CreatedAt,
		&finding.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrFindingNotFound
	}

	if err != nil {
		return nil, err
	}

	return &finding, nil
}

func (r *FindingRepository) FindAll(
	ctx context.Context,
) ([]models.Finding, error) {
	query := `
	SELECT
		id,
		assessment_id,
		task_id,
		agent_id,
		title,
		description,
		severity,
		technique_id,
		technique_name,
		evidence,
		recommendation,
		created_at,
		updated_at
	FROM findings
	ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	findings := make([]models.Finding, 0)

	for rows.Next() {
		var finding models.Finding

		if err := rows.Scan(
			&finding.ID,
			&finding.AssessmentID,
			&finding.TaskID,
			&finding.AgentID,
			&finding.Title,
			&finding.Description,
			&finding.Severity,
			&finding.TechniqueID,
			&finding.TechniqueName,
			&finding.Evidence,
			&finding.Recommendation,
			&finding.CreatedAt,
			&finding.UpdatedAt,
		); err != nil {
			return nil, err
		}

		findings = append(
			findings,
			finding,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return findings, nil
}
