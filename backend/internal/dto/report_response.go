package dto

type ReportResponse struct {
	ID           string `json:"id"`
	AssessmentID string `json:"assessment_id"`
	FileName     string `json:"file_name"`
	FilePath     string `json:"file_path"`
}