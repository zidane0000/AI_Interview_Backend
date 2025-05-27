package api

import "time"

// Data Transfer Objects (DTOs) for API request and response payloads:
// - CreateInterviewRequestDTO
// - InterviewResponseDTO
// - ListInterviewsResponseDTO
// - SubmitEvaluationRequestDTO
// - EvaluationResponseDTO
//
// These DTOs define the JSON structure for all RESTful API endpoints.
// Use these types for marshaling/unmarshaling and handler signatures.

// --- Interview DTOs ---
type CreateInterviewRequestDTO struct {
	CandidateName string   `json:"candidate_name"`
	Questions     []string `json:"questions"`
}

type InterviewResponseDTO struct {
	ID            string    `json:"id"`
	CandidateName string    `json:"candidate_name"`
	Questions     []string  `json:"questions"`
	CreatedAt     time.Time `json:"created_at"`
}

type ListInterviewsResponseDTO struct {
	Interviews []InterviewResponseDTO `json:"interviews"`
}

// --- Evaluation DTOs ---
type SubmitEvaluationRequestDTO struct {
	InterviewID string            `json:"interview_id"`
	Answers     map[string]string `json:"answers"`
}

type EvaluationResponseDTO struct {
	ID          string    `json:"id"`
	InterviewID string    `json:"interview_id"`
	Score       float64   `json:"score"`
	Feedback    string    `json:"feedback"`
	CreatedAt   time.Time `json:"created_at"`
}

// --- Error DTO ---
type ErrorResponseDTO struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
