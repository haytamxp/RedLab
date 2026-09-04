package reporting

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrReportAssessmentNotFound = errors.New(
	"assessment not found",
)

type Service struct {
	db *pgxpool.Pool
}

func NewService(
	db *pgxpool.Pool,
) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Build(
	ctx context.Context,
	assessmentID uuid.UUID,
) (*Report, error) {
	assessment, err := s.loadAssessment(
		ctx,
		assessmentID,
	)
	if err != nil {
		return nil, err
	}

	findings, err := s.loadFindings(
		ctx,
		assessmentID,
	)
	if err != nil {
		return nil, err
	}

	severityCounts := make(map[string]int)
	for _, finding := range findings {
		severityCounts[strings.ToUpper(
			strings.TrimSpace(
				finding.Severity,
			),
		)]++
	}

	return &Report{
		Assessment:     *assessment,
		Findings:       findings,
		SeverityCounts: severityCounts,
		GeneratedAt:    time.Now(),
	}, nil
}

func (s *Service) loadAssessment(
	ctx context.Context,
	id uuid.UUID,
) (*AssessmentReport, error) {
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

	var report AssessmentReport

	err := s.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&report.ID,
		&report.Name,
		&report.Description,
		&report.AgentID,
		&report.Status,
		&report.CreatedAt,
		&report.UpdatedAt,
	)

	if errors.Is(
		err,
		pgx.ErrNoRows,
	) {
		return nil, ErrReportAssessmentNotFound
	}

	if err != nil {
		return nil, fmt.Errorf(
			"load assessment: %w",
			err,
		)
	}

	return &report, nil
}

func (s *Service) loadFindings(
	ctx context.Context,
	assessmentID uuid.UUID,
) ([]FindingReport, error) {
	query := `
	SELECT
		id,
		title,
		description,
		severity,
		technique_id,
		technique_name,
		recommendation,
		created_at
	FROM findings
	WHERE assessment_id = $1
	ORDER BY
		CASE severity
			WHEN 'CRITICAL' THEN 1
			WHEN 'HIGH' THEN 2
			WHEN 'MEDIUM' THEN 3
			WHEN 'LOW' THEN 4
			WHEN 'INFO' THEN 5
			ELSE 6
		END,
		created_at DESC
	`

	rows, err := s.db.Query(
		ctx,
		query,
		assessmentID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"load findings: %w",
			err,
		)
	}
	defer rows.Close()

	findings := make(
		[]FindingReport,
		0,
	)

	for rows.Next() {
		var finding FindingReport

		if err := rows.Scan(
			&finding.ID,
			&finding.Title,
			&finding.Description,
			&finding.Severity,
			&finding.TechniqueID,
			&finding.TechniqueName,
			&finding.Recommendation,
			&finding.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf(
				"scan finding: %w",
				err,
			)
		}

		findings = append(
			findings,
			finding,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate findings: %w",
			err,
		)
	}

	return findings, nil
}
