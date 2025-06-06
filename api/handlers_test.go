package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/zidane0000/AI_Interview_Backend/data"
)

// clearMemoryStore clears all data from the memory store for test isolation
func clearMemoryStore() {
	var err error
	data.GlobalStore, err = data.NewHybridStore(data.BackendMemory, "")
	if err != nil {
		panic("Failed to initialize test store: " + err.Error())
	}
}

func TestCreateInterviewHandler_Success(t *testing.T) {
	clearMemoryStore() // Clear store for test isolation
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
	clearMemoryStore() // Clear store for test isolation
	router := SetupRouter()
	req := httptest.NewRequest("GET", "/interviews", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	var resp ListInterviewsResponseDTO
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Interviews) != 0 {
		t.Errorf("expected empty interviews list, got %d interviews", len(resp.Interviews))
	}
	if resp.Total != 0 {
		t.Errorf("expected total count of 0, got %d", resp.Total)
	}
}

func TestListInterviewsHandler_WithData(t *testing.T) {
	clearMemoryStore() // Clear store for test isolation
	router := SetupRouter()

	// Create multiple test interviews
	interviews := []CreateInterviewRequestDTO{
		{CandidateName: "Alice Johnson", Questions: []string{"Q1", "Q2"}},
		{CandidateName: "Bob Smith", Questions: []string{"Q3", "Q4"}},
		{CandidateName: "Charlie Brown", Questions: []string{"Q5", "Q6"}},
	}

	// Create interviews
	for _, interview := range interviews {
		b, _ := json.Marshal(interview)
		req := httptest.NewRequest("POST", "/interviews", bytes.NewReader(b))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("failed to create interview for %s, got %d", interview.CandidateName, w.Code)
		}
	}

	// Test listing all interviews
	req := httptest.NewRequest("GET", "/interviews", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	var resp ListInterviewsResponseDTO
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Interviews) != 3 {
		t.Errorf("expected 3 interviews, got %d", len(resp.Interviews))
	}
	if resp.Total != 3 {
		t.Errorf("expected total count of 3, got %d", resp.Total)
	}
}

func TestListInterviewsHandler_Pagination(t *testing.T) {
	clearMemoryStore() // Clear store for test isolation
	router := SetupRouter()

	// Create 5 test interviews
	for i := 1; i <= 5; i++ {
		interview := CreateInterviewRequestDTO{
			CandidateName: fmt.Sprintf("Candidate %d", i),
			Questions:     []string{"Q1", "Q2"},
		}
		b, _ := json.Marshal(interview)
		req := httptest.NewRequest("POST", "/interviews", bytes.NewReader(b))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("failed to create interview %d, got %d", i, w.Code)
		}
	}

	// Test pagination with limit=2
	req := httptest.NewRequest("GET", "/interviews?limit=2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	var resp ListInterviewsResponseDTO
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Interviews) != 2 {
		t.Errorf("expected 2 interviews with limit=2, got %d", len(resp.Interviews))
	}
	if resp.Total != 5 {
		t.Errorf("expected total count of 5, got %d", resp.Total)
	}

	// Test pagination with offset
	req = httptest.NewRequest("GET", "/interviews?limit=2&offset=2", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Interviews) != 2 {
		t.Errorf("expected 2 interviews with offset=2, got %d", len(resp.Interviews))
	}
	if resp.Total != 5 {
		t.Errorf("expected total count of 5, got %d", resp.Total)
	}
}

func TestListInterviewsHandler_Filtering(t *testing.T) {
	clearMemoryStore() // Clear store for test isolation
	router := SetupRouter()

	// Create test interviews with different names
	interviews := []CreateInterviewRequestDTO{
		{CandidateName: "Alice Johnson", Questions: []string{"Q1", "Q2"}},
		{CandidateName: "Bob Alice", Questions: []string{"Q3", "Q4"}},
		{CandidateName: "Charlie Brown", Questions: []string{"Q5", "Q6"}},
	}

	for _, interview := range interviews {
		b, _ := json.Marshal(interview)
		req := httptest.NewRequest("POST", "/interviews", bytes.NewReader(b))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("failed to create interview for %s, got %d", interview.CandidateName, w.Code)
		}
	}

	// Test filtering by candidate name
	req := httptest.NewRequest("GET", "/interviews?candidate_name=Alice", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	var resp ListInterviewsResponseDTO
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Interviews) != 2 {
		t.Errorf("expected 2 interviews containing 'Alice', got %d", len(resp.Interviews))
	}
	if resp.Total != 2 {
		t.Errorf("expected total count of 2, got %d", resp.Total)
	}

	// Verify the filtered results contain "Alice"
	for _, interview := range resp.Interviews {
		if !strings.Contains(strings.ToLower(interview.CandidateName), "alice") {
			t.Errorf("expected interview name to contain 'alice', got %s", interview.CandidateName)
		}
	}
}

func TestListInterviewsHandler_Sorting(t *testing.T) {
	clearMemoryStore() // Clear store for test isolation
	router := SetupRouter()

	// Create test interviews in a specific order
	interviews := []CreateInterviewRequestDTO{
		{CandidateName: "Charlie Brown", Questions: []string{"Q1", "Q2"}},
		{CandidateName: "Alice Johnson", Questions: []string{"Q3", "Q4"}},
		{CandidateName: "Bob Smith", Questions: []string{"Q5", "Q6"}},
	}

	for _, interview := range interviews {
		b, _ := json.Marshal(interview)
		req := httptest.NewRequest("POST", "/interviews", bytes.NewReader(b))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("failed to create interview for %s, got %d", interview.CandidateName, w.Code)
		}
		// Add small delay to ensure different creation times
		time.Sleep(1 * time.Millisecond)
	}

	// Test sorting by name ascending
	req := httptest.NewRequest("GET", "/interviews?sort_by=name&sort_order=asc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	var resp ListInterviewsResponseDTO
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Interviews) != 3 {
		t.Errorf("expected 3 interviews, got %d", len(resp.Interviews))
	}

	// Verify the sorting order
	expectedOrder := []string{"Alice Johnson", "Bob Smith", "Charlie Brown"}
	for i, interview := range resp.Interviews {
		if interview.CandidateName != expectedOrder[i] {
			t.Errorf("expected interview %d to be %s, got %s", i, expectedOrder[i], interview.CandidateName)
		}
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
	clearMemoryStore() // Clear store for test isolation
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
	clearMemoryStore() // Clear store for test isolation
	// First create a valid interview
	interview := &data.Interview{
		ID:            "test-interview-123",
		CandidateName: "Test Candidate",
		Questions:     []string{"What is your experience?", "Tell me about yourself"}, CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := data.GlobalStore.CreateInterview(interview); err != nil {
		t.Fatalf("failed to create interview: %v", err)
	}

	body := SubmitEvaluationRequestDTO{
		InterviewID: "test-interview-123",
		Answers:     map[string]string{"question_0": "5 years of experience", "question_1": "I am a developer"},
	}
	b, _ := json.Marshal(body)

	router := SetupRouter()
	req := httptest.NewRequest("POST", "/evaluation", bytes.NewReader(b))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

func TestSubmitEvaluationHandler_BadRequest(t *testing.T) {
	router := SetupRouter()

	// Invalid JSON
	req := httptest.NewRequest("POST", "/evaluation", bytes.NewReader([]byte("{")))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", w.Code)
	}

	// Missing fields
	body := SubmitEvaluationRequestDTO{}
	b, _ := json.Marshal(body)
	req = httptest.NewRequest("POST", "/evaluation", bytes.NewReader(b))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
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
	clearMemoryStore() // Clear store for test isolation
	// First create a valid evaluation
	evaluation := &data.Evaluation{
		ID:          "test-evaluation-456",
		InterviewID: "test-interview-456",
		Answers:     map[string]string{"question_0": "Test answer"},
		Score:       0.8,
		Feedback:    "Good performance",
		CreatedAt:   time.Now(), UpdatedAt: time.Now(),
	}
	if err := data.GlobalStore.CreateEvaluation(evaluation); err != nil {
		t.Fatalf("failed to create evaluation: %v", err)
	}

	router := SetupRouter()
	req := httptest.NewRequest("GET", "/evaluation/test-evaluation-456", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}
