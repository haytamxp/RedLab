package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/haytamxp/redlab/backend/internal/models"
)

var ErrTaskNotFound = errors.New(
	"task not found",
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(
	db *pgxpool.Pool,
) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}


func (r *TaskRepository) Create(
	ctx context.Context,
	task *models.Task,
) error {
	query := `
	INSERT INTO tasks (
		id,
		agent_id,
		type,
		payload,
		status,
		priority,
		created_at,
		updated_at
	)
	VALUES (
		$1,$2,$3,$4,$5,$6,$7,$8
	)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		task.ID,
		task.AgentID,
		task.Type,
		task.Payload,
		task.Status,
		task.Priority,
		task.CreatedAt,
		task.UpdatedAt,
	)

	return err
}


func (r *TaskRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Task, error) {
	query := `
	SELECT
		id,
		agent_id,
		type,
		payload,
		status,
		priority,
		result,
		error_message,
		created_at,
		updated_at,
		claimed_at,
		completed_at
	FROM tasks
	WHERE id = $1
	`

	return r.scanTask(
		r.db.QueryRow(
			ctx,
			query,
			id,
		),
	)
}


func (r *TaskRepository) ClaimNext(
	ctx context.Context,
	agentID uuid.UUID,
) (*models.Task, error) {
	query := `
	WITH next_task AS (
		SELECT id
		FROM tasks
		WHERE agent_id = $1
		  AND status = 'PENDING'
		ORDER BY
			priority DESC,
			created_at ASC
		LIMIT 1
		FOR UPDATE SKIP LOCKED
	)
	UPDATE tasks AS t
	SET
		status = 'CLAIMED',
		claimed_at = NOW(),
		updated_at = NOW()
	FROM next_task
	WHERE t.id = next_task.id
	RETURNING
		t.id,
		t.agent_id,
		t.type,
		t.payload,
		t.status,
		t.priority,
		t.result,
		t.error_message,
		t.created_at,
		t.updated_at,
		t.claimed_at,
		t.completed_at
	`

	task, err :=
		r.scanTask(
			r.db.QueryRow(
				ctx,
				query,
				agentID,
			),
		)

	if errors.Is(
		err,
		pgx.ErrNoRows,
	) {
		return nil, nil
	}

	return task, err
}


func (r *TaskRepository) Complete(
	ctx context.Context,
	taskID uuid.UUID,
	agentID uuid.UUID,
	status models.TaskStatus,
	result []byte,
	errorMessage *string,
) (*models.Task, error) {
	query := `
	UPDATE tasks
	SET
		status = $1,
		result = $2,
		error_message = $3,
		completed_at = NOW(),
		updated_at = NOW()
	WHERE id = $4
	  AND agent_id = $5
	  AND status = 'CLAIMED'
	RETURNING
		id,
		agent_id,
		type,
		payload,
		status,
		priority,
		result,
		error_message,
		created_at,
		updated_at,
		claimed_at,
		completed_at
	`

	task, err :=
		r.scanTask(
			r.db.QueryRow(
				ctx,
				query,
				status,
				result,
				errorMessage,
				taskID,
				agentID,
			),
		)

	if errors.Is(
		err,
		pgx.ErrNoRows,
	) {
		return nil, ErrTaskNotFound
	}

	return task, err
}


func (r *TaskRepository) ListForAgent(
	ctx context.Context,
	agentID uuid.UUID,
) ([]models.Task, error) {
	query := `
	SELECT
		id,
		agent_id,
		type,
		payload,
		status,
		priority,
		result,
		error_message,
		created_at,
		updated_at,
		claimed_at,
		completed_at
	FROM tasks
	WHERE agent_id = $1
	ORDER BY created_at DESC
	`

	rows, err :=
		r.db.Query(
			ctx,
			query,
			agentID,
		)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := make(
		[]models.Task,
		0,
	)

	for rows.Next() {
		task, err :=
			r.scanTask(rows)

		if err != nil {
			return nil, err
		}

		tasks = append(
			tasks,
			*task,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}


/*
 * Operator task management.
 */

func (r *TaskRepository) ListAll(
	ctx context.Context,
) ([]models.Task, error) {
	query := `
	SELECT
		id,
		agent_id,
		type,
		payload,
		status,
		priority,
		result,
		error_message,
		created_at,
		updated_at,
		claimed_at,
		completed_at
	FROM tasks
	ORDER BY created_at DESC
	`

	rows, err :=
		r.db.Query(
			ctx,
			query,
		)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := make(
		[]models.Task,
		0,
	)

	for rows.Next() {
		task, err :=
			r.scanTask(rows)

		if err != nil {
			return nil, err
		}

		tasks = append(
			tasks,
			*task,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}


func (r *TaskRepository) DeletePending(
	ctx context.Context,
	id uuid.UUID,
) error {
	query := `
	DELETE FROM tasks
	WHERE id = $1
	  AND status = 'PENDING'
	`

	result, err :=
		r.db.Exec(
			ctx,
			query,
			id,
		)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrTaskNotFound
	}

	return nil
}


type taskScanner interface {
	Scan(dest ...any) error
}


func (r *TaskRepository) scanTask(
	row taskScanner,
) (*models.Task, error) {
	var task models.Task

	err := row.Scan(
		&task.ID,
		&task.AgentID,
		&task.Type,
		&task.Payload,
		&task.Status,
		&task.Priority,
		&task.Result,
		&task.ErrorMessage,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.ClaimedAt,
		&task.CompletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &task, nil
}