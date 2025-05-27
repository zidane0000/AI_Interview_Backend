// HTTP handler functions for each endpoint
package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// Helper: write JSON response
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// Helper: write JSON error response
func writeJSONError(w http.ResponseWriter, status int, msg string, details ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	errResp := ErrorResponseDTO{Error: msg}
	if len(details) > 0 {
		errResp.Details = details[0]
	}
	json.NewEncoder(w).Encode(errResp)
}

// CreateInterviewHandler handles POST /interviews
func CreateInterviewHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateInterviewRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}
	if req.CandidateName == "" || len(req.Questions) == 0 {
		writeJSONError(w, http.StatusBadRequest, "Missing candidate_name or questions")
		return
	}
	resp := InterviewResponseDTO{
		ID:            "sample-id", // TODO: generate real ID
		CandidateName: req.CandidateName,
		Questions:     req.Questions,
		CreatedAt:     time.Now(),
	}
	writeJSON(w, http.StatusCreated, resp)
}

// ListInterviewsHandler handles GET /interviews
func ListInterviewsHandler(w http.ResponseWriter, r *http.Request) {
	resp := ListInterviewsResponseDTO{
		Interviews: []InterviewResponseDTO{}, // TODO: fetch from DB
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetInterviewHandler handles GET /interviews/{id}
func GetInterviewHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSONError(w, ErrCodeBadRequest, ErrMsgMissingInterviewID)
		return
	}
	// TODO: fetch interview by id
	resp := InterviewResponseDTO{
		ID:            id,
		CandidateName: "Sample Name",
		Questions:     []string{"Q1", "Q2"},
		CreatedAt:     time.Now(),
	}
	writeJSON(w, http.StatusOK, resp)
}

// SubmitEvaluationHandler handles POST /evaluation
func SubmitEvaluationHandler(w http.ResponseWriter, r *http.Request) {
	var req SubmitEvaluationRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}
	if req.InterviewID == "" || len(req.Answers) == 0 {
		writeJSONError(w, http.StatusBadRequest, "Missing interview_id or answers")
		return
	}
	resp := EvaluationResponseDTO{
		ID:          "sample-eval-id", // TODO: generate real ID
		InterviewID: req.InterviewID,
		Score:       0.95,
		Feedback:    "Sample feedback",
		CreatedAt:   time.Now(),
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetEvaluationHandler handles GET /evaluation/{id}
func GetEvaluationHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSONError(w, ErrCodeBadRequest, ErrMsgMissingEvaluationID)
		return
	}
	// TODO: fetch evaluation by id
	resp := EvaluationResponseDTO{
		ID:          id,
		InterviewID: "sample-interview-id",
		Score:       0.95,
		Feedback:    "Sample feedback",
		CreatedAt:   time.Now(),
	}
	writeJSON(w, http.StatusOK, resp)
}
