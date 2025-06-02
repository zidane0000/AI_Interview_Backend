// Mock AI provider for testing and CI environments
package ai

import (
	"context"
	"time"
)

// MockProvider implements the AIProvider interface with canned responses
type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (m *MockProvider) GenerateResponse(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	return &ChatResponse{
		Content:      "[MOCK] Hello! This is a mock AI response.",
		FinishReason: "stop",
		TokensUsed:   TokenUsage{PromptTokens: 1, CompletionTokens: 1, TotalTokens: 2},
		Model:        "mock-model",
		Provider:     "mock",
		ResponseTime: 10 * time.Millisecond,
		Timestamp:    time.Now(),
	}, nil
}

func (m *MockProvider) GenerateStreamResponse(ctx context.Context, req *ChatRequest) (<-chan *ChatResponse, error) {
	ch := make(chan *ChatResponse, 1)
	ch <- &ChatResponse{
		Content:      "[MOCK] Streaming response chunk.",
		FinishReason: "stop",
		TokensUsed:   TokenUsage{PromptTokens: 1, CompletionTokens: 1, TotalTokens: 2},
		Model:        "mock-model",
		Provider:     "mock",
		ResponseTime: 5 * time.Millisecond,
		Timestamp:    time.Now(),
	}
	close(ch)
	return ch, nil
}

func (m *MockProvider) GenerateInterviewQuestions(ctx context.Context, req *QuestionGenerationRequest) (*QuestionGenerationResponse, error) {
	return &QuestionGenerationResponse{
		Questions: []InterviewQuestion{{
			Question:   "[MOCK] What is your greatest strength?",
			Category:   "general",
			Difficulty: "easy",
		}},
		Rationale:  "[MOCK] Standard question.",
		TokensUsed: TokenUsage{PromptTokens: 1, CompletionTokens: 1, TotalTokens: 2},
		Provider:   "mock",
		Model:      "mock-model",
		Timestamp:  time.Now(),
	}, nil
}

func (m *MockProvider) EvaluateAnswers(ctx context.Context, req *EvaluationRequest) (*EvaluationResponse, error) {
	return &EvaluationResponse{
		OverallScore:    0.8,
		CategoryScores:  map[string]float64{"general": 0.8},
		Feedback:        "[MOCK] Good job!",
		Strengths:       []string{"[MOCK] Communication"},
		Weaknesses:      []string{"[MOCK] None"},
		Recommendations: []string{"[MOCK] Keep practicing."},
		TokensUsed:      TokenUsage{PromptTokens: 1, CompletionTokens: 1, TotalTokens: 2},
		Provider:        "mock",
		Model:           "mock-model",
		Timestamp:       time.Now(),
	}, nil
}

func (m *MockProvider) GetProviderName() string                       { return "mock" }
func (m *MockProvider) GetSupportedModels() []string                  { return []string{"mock-model"} }
func (m *MockProvider) ValidateCredentials(ctx context.Context) error { return nil }
func (m *MockProvider) IsHealthy(ctx context.Context) bool            { return true }
func (m *MockProvider) GetUsageStats(ctx context.Context) (map[string]interface{}, error) {
	return map[string]interface{}{"mock": true}, nil
}
