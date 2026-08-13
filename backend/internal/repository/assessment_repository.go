package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/haytamxp/redlab/backend/internal/models"
)

var ErrAssessmentNotFound = errors.New(
	"assessment not found",
)

type AssessmentRepository struct {
	db *pgxpool.Pool
}

func NewAssessmentRepository(
	db *pgxpool.Pool,
) *AssessmentRepository {
	return &AssessmentRepository{
		db: db,
	}
}

func (r *AssessmentRepository) Create(
	ctx context.Context,
	assessment *models.Assessment,
) error {
	query := `
	INSERT INTO assessments (
		id,
		name,
		description,
		agent_id,
		status,
		created_at,
		updated_at
	)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7
	)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		assessment.ID,
		assessment.Name,
		assessment.Description,
		assessment.AgentID,
		assessment.Status,
		assessment.CreatedAt,
		assessment.UpdatedAt,
	)

	return err
}

func (r *AssessmentRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Assessment, error) {
	query := `
	SELECT
		id,
		name,
		description,
		agent_id,
		status,
		created_at,
		updated_at
	FROM assessments
	WHERE id = $1
	`

	var assessment models.Assessment

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&assessment.ID,
		&assessment.Name,
		&assessment.Description,
		&assessment.AgentID,
		&assessment.Status,
		&assessment.CreatedAt,
		&assessment.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrAssessmentNotFound
	}

	if err != nil {
		return nil, err
	}

	return &assessment, nil
}

func (r *AssessmentRepository) FindAll(
	ctx context.Context,
) ([]models.Assessment, error) {
	query := `
	SELECT
		id,
		name,
		description,
		agent_id,
		status,
		created_at,
		updated_at
	FROM assessments
	ORDER BY created_at DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assessments := make(
		[]models.Assessment,
		0,
	)

	for rows.Next() {
		var assessment models.Assessment

		if err := rows.Scan(
			&assessment.ID,
			&assessment.Name,
			&assessment.Description,
			&assessment.AgentID,
			&assessment.Status,
			&assessment.CreatedAt,
			&assessment.UpdatedAt,
		); err != nil {
			return nil, err
		}

		assessments = append(
			assessments,
			assessment,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assessments, nil
}

func (r *AssessmentRepository) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	status models.AssessmentStatus,
) error {
	query := `
	UPDATE assessments
	SET
		status = $1,
		updated_at = NOW()
	WHERE id = $2
	`

	result, err := r.db.Exec(
		ctx,
		query,
		status,
		id,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrAssessmentNotFound
	}

	return nil
}
