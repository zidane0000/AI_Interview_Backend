package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

// DTO types for testing (mirrored from api package)
type InterviewResponseDTO struct {
	ID            string    `json:"id"`
	CandidateName string    `json:"candidate_name"`
	Questions     []string  `json:"questions"`
	CreatedAt     time.Time `json:"created_at"`
}

type ChatMessageDTO struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type ChatInterviewSessionDTO struct {
	ID          string           `json:"id"`
	InterviewID string           `json:"interview_id"`
	Messages    []ChatMessageDTO `json:"messages"`
	Status      string           `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
}

type SendMessageRequestDTO struct {
	Message string `json:"message"`
}

type SendMessageResponseDTO struct {
	Message       ChatMessageDTO  `json:"message"`
	AIResponse    *ChatMessageDTO `json:"ai_response,omitempty"`
	SessionStatus string          `json:"session_status"` // "active" or "completed"
}

type EvaluationResponseDTO struct {
	ID          string            `json:"id"`
	InterviewID string            `json:"interview_id"`
	Answers     map[string]string `json:"answers"`
	Score       float64           `json:"score"`
	Feedback    string            `json:"feedback"`
	CreatedAt   time.Time         `json:"created_at"`
}

// Test helper functions for E2E tests

func GetAPIBaseURL() string {
	if v := os.Getenv("API_BASE_URL"); v != "" {
		return v
	}
	return "http://localhost:8080"
}

// CreateTestInterview creates a test interview and returns the response
func CreateTestInterview(t *testing.T, candidateName string, questions []string) InterviewResponseDTO {
	t.Helper()
	baseURL := GetAPIBaseURL()

	createReq := map[string]interface{}{
		"candidate_name": candidateName,
		"questions":      questions,
	}

	reqBody, _ := json.Marshal(createReq)
	resp, err := http.Post(baseURL+"/interviews", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create interview: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d", resp.StatusCode)
	}

	var interview InterviewResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&interview); err != nil {
		t.Fatalf("Failed to decode interview response: %v", err)
	}
	if interview.ID == "" {
		t.Fatalf("Interview ID is empty")
	}

	return interview
}

// StartChatSession starts a chat session for the given interview
func StartChatSession(t *testing.T, interviewID string) ChatInterviewSessionDTO {
	t.Helper()
	baseURL := GetAPIBaseURL()

	resp, err := http.Post(fmt.Sprintf("%s/interviews/%s/chat/start", baseURL, interviewID), "application/json", nil)
	if err != nil {
		t.Fatalf("Failed to start chat session: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d", resp.StatusCode)
	}

	var chatSession ChatInterviewSessionDTO
	if err := json.NewDecoder(resp.Body).Decode(&chatSession); err != nil {
		t.Fatalf("Failed to decode chat session: %v", err)
	}

	return chatSession
}

// SendMessage sends a message in a chat session
func SendMessage(t *testing.T, sessionID, message string) SendMessageResponseDTO {
	t.Helper()
	baseURL := GetAPIBaseURL()

	sendMsgReq := SendMessageRequestDTO{
		Message: message,
	}
	reqBody, _ := json.Marshal(sendMsgReq)
	resp, err := http.Post(fmt.Sprintf("%s/chat/%s/message", baseURL, sessionID), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var msgResponse SendMessageResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&msgResponse); err != nil {
		t.Fatalf("Failed to decode send message response: %v", err)
	}

	return msgResponse
}

// GetChatSession retrieves chat session state
func GetChatSession(t *testing.T, sessionID string) ChatInterviewSessionDTO {
	t.Helper()
	baseURL := GetAPIBaseURL()

	resp, err := http.Get(fmt.Sprintf("%s/chat/%s", baseURL, sessionID))
	if err != nil {
		t.Fatalf("Failed to get chat session: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var session ChatInterviewSessionDTO
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		t.Fatalf("Failed to decode chat session: %v", err)
	}

	return session
}

// EndChatSession ends a chat session and returns evaluation
func EndChatSession(t *testing.T, sessionID string) EvaluationResponseDTO {
	t.Helper()
	baseURL := GetAPIBaseURL()

	resp, err := http.Post(fmt.Sprintf("%s/chat/%s/end", baseURL, sessionID), "application/json", nil)
	if err != nil {
		t.Fatalf("Failed to end chat session: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var evaluation EvaluationResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&evaluation); err != nil {
		t.Fatalf("Failed to decode evaluation: %v", err)
	}

	return evaluation
}

// AssertErrorResponse checks if response contains expected error
func AssertErrorResponse(t *testing.T, resp *http.Response, expectedStatus int, expectedMessage string) {
	t.Helper()
	if resp.StatusCode != expectedStatus {
		t.Errorf("Expected status %d, got %d", expectedStatus, resp.StatusCode)
	}

	var errorResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if errorMsg, ok := errorResp["error"].(string); ok {
		if errorMsg != expectedMessage {
			t.Errorf("Expected error message '%s', got '%s'", expectedMessage, errorMsg)
		}
	} else {
		t.Errorf("Error response missing 'error' field")
	}
}

// Sample test data generators
func GetSampleQuestions() []string {
	return []string{
		"Tell me about yourself",
		"What are your strengths?",
		"Describe a challenging project you worked on",
		"Where do you see yourself in 5 years?",
	}
}

func GetLongMessage() string {
	return "This is a very long message that contains a lot of text to test how the system handles longer inputs. " +
		"It includes multiple sentences and should test the limits of message processing. " +
		"The message continues with more content to ensure we test various edge cases related to message length and content processing. " +
		"This helps us verify that the chat system can handle realistic user inputs of varying lengths."
}

func GetSpecialCharacterMessage() string {
	return "Message with special chars: 你好 🚀 @#$%^&*() \"quotes\" 'apostrophes' and\nnewlines"
}
