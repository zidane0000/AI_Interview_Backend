package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateInterviewHandler_Success(t *testing.T) {
	body := CreateInterviewRequestDTO{
		CandidateName: "Alice",
		Questions:     []string{"Q1", "Q2"},
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/interviews", bytes.NewReader(b))
	w := httptest.NewRecorder()
	CreateInterviewHandler(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", w.Code)
	}
}

func TestCreateInterviewHandler_BadRequest(t *testing.T) {
	// Invalid JSON
	req := httptest.NewRequest("POST", "/interviews", bytes.NewReader([]byte("{")))
	w := httptest.NewRecorder()
	CreateInterviewHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", w.Code)
	}

	// Missing fields
	body := CreateInterviewRequestDTO{}
	b, _ := json.Marshal(body)
	req = httptest.NewRequest("POST", "/interviews", bytes.NewReader(b))
	w = httptest.NewRecorder()
	CreateInterviewHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for missing fields, got %d", w.Code)
	}
}

func TestListInterviewsHandler_Empty(t *testing.T) {
	router := SetupRouter()
	req := httptest.NewRequest("GET", "/interviews", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

func TestGetInterviewHandler_BadRequest(t *testing.T) {
	router := SetupRouter()
	req := httptest.NewRequest("GET", "/interviews/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

func TestGetInterviewHandler_Success(t *testing.T) {
	router := SetupRouter()

	// Step 1: Create an interview
	createBody := CreateInterviewRequestDTO{
		CandidateName: "Test User",
		Questions:     []string{"Q1", "Q2"},
	}
	b, _ := json.Marshal(createBody)
	createReq := httptest.NewRequest("POST", "/interviews", bytes.NewReader(b))
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)
	if createW.Code != http.StatusCreated {
		t.Fatalf("failed to create interview, got %d", createW.Code)
	}
	var createdResp InterviewResponseDTO
	if err := json.Unmarshal(createW.Body.Bytes(), &createdResp); err != nil {
		t.Fatalf("failed to decode create response: %v", err)
	}

	// Step 2: Use the real ID for GET
	req := httptest.NewRequest("GET", "/interviews/"+createdResp.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

func TestSubmitEvaluationHandler_Success(t *testing.T) {
	body := SubmitEvaluationRequestDTO{
		InterviewID: "abc",
		Answers:     map[string]string{"Q1": "A1"},
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/evaluation", bytes.NewReader(b))
	w := httptest.NewRecorder()
	SubmitEvaluationHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

func TestSubmitEvaluationHandler_BadRequest(t *testing.T) {
	// Invalid JSON
	req := httptest.NewRequest("POST", "/evaluation", bytes.NewReader([]byte("{")))
	w := httptest.NewRecorder()
	SubmitEvaluationHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", w.Code)
	}

	// Missing fields
	body := SubmitEvaluationRequestDTO{}
	b, _ := json.Marshal(body)
	req = httptest.NewRequest("POST", "/evaluation", bytes.NewReader(b))
	w = httptest.NewRecorder()
	SubmitEvaluationHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for missing fields, got %d", w.Code)
	}
}

func TestGetEvaluationHandler_BadRequest(t *testing.T) {
	router := SetupRouter()
	req := httptest.NewRequest("GET", "/evaluation/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 Method Not Allowed, got %d", w.Code)
	}
}

func TestGetEvaluationHandler_Success(t *testing.T) {
	router := SetupRouter()
	req := httptest.NewRequest("GET", "/evaluation/456", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}
