package dto

type FindingRequest struct {
	AssessmentID string `json:"assessment_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Severity     string `json:"severity"`
}
