package reporting

import "time"

type Report struct {
	Assessment     AssessmentReport
	Findings       []FindingReport
	SeverityCounts map[string]int
	GeneratedAt    time.Time
}

type AssessmentReport struct {
	ID          string
	Name        string
	Description string
	AgentID     string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FindingReport struct {
	ID             string
	Title          string
	Description    string
	Severity       string
	TechniqueID    string
	TechniqueName  string
	Recommendation string
	CreatedAt      time.Time
}
