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
	// Generate more substantial mock responses for testing
	mockResponses := []string{
		"[MOCK] Thank you for sharing that information. Can you tell me more about your experience with software development and the technologies you've worked with?",
		"[MOCK] That's interesting. How do you approach problem-solving when faced with complex technical challenges?",
		"[MOCK] I appreciate your detailed response. Could you walk me through a specific project where you demonstrated leadership skills?",
		"[MOCK] Thank you for the explanation. What motivates you in your professional work, and how do you stay updated with industry trends?",
		"[MOCK] That's a great example. How do you handle working under pressure and tight deadlines?",
		"[MOCK] I see. Can you describe a situation where you had to collaborate with cross-functional teams?",
		"[MOCK] Thank you for sharing your experience. What are your career goals for the next few years?",
		"[MOCK] That concludes our interview. Thank you for your time and thoughtful responses throughout our conversation.",
	}

	// Use a simple hash of the request to get consistent but varied responses
	responseIndex := len(req.Messages) % len(mockResponses)

	return &ChatResponse{
		Content:      mockResponses[responseIndex],
		FinishReason: "stop",
		TokensUsed:   TokenUsage{PromptTokens: 10, CompletionTokens: 20, TotalTokens: 30},
		Model:        "mock-model",
		Provider:     "mock",
		ResponseTime: 10 * time.Millisecond,
		Timestamp:    time.Now(),
	}, nil
}

func (m *MockProvider) GenerateStreamResponse(ctx context.Context, req *ChatRequest) (<-chan *ChatResponse, error) {
	ch := make(chan *ChatResponse, 1)
	ch <- &ChatResponse{
		Content:      "[MOCK] Thank you for your response. This is a streaming mock response that provides more detailed feedback and continues the conversation naturally.",
		FinishReason: "stop",
		TokensUsed:   TokenUsage{PromptTokens: 10, CompletionTokens: 20, TotalTokens: 30},
		Model:        "mock-model",
		Provider:     "mock",
		ResponseTime: 5 * time.Millisecond,
		Timestamp:    time.Now(),
	}
	close(ch)
	return ch, nil
}

func (m *MockProvider) GenerateInterviewQuestions(ctx context.Context, req *QuestionGenerationRequest) (*QuestionGenerationResponse, error) {
	// Generate more realistic mock questions
	questions := []InterviewQuestion{
		{
			Question:   "[MOCK] Can you walk me through your experience with software development and the technologies you've worked with?",
			Category:   "technical",
			Difficulty: "medium",
		},
		{
			Question:   "[MOCK] Describe a challenging project you worked on and how you overcame the obstacles.",
			Category:   "behavioral",
			Difficulty: "medium",
		},
		{
			Question:   "[MOCK] How do you approach debugging and troubleshooting complex issues?",
			Category:   "technical",
			Difficulty: "medium",
		},
	}

	return &QuestionGenerationResponse{
		Questions:  questions,
		Rationale:  "[MOCK] These questions are designed to assess both technical competency and problem-solving abilities, providing a comprehensive evaluation of the candidate's skills and experience.",
		TokensUsed: TokenUsage{PromptTokens: 20, CompletionTokens: 40, TotalTokens: 60},
		Provider:   "mock",
		Model:      "mock-model",
		Timestamp:  time.Now(),
	}, nil
}

func (m *MockProvider) EvaluateAnswers(ctx context.Context, req *EvaluationRequest) (*EvaluationResponse, error) {
	// Generate substantial mock evaluation feedback for testing
	feedback := `[MOCK] Overall Performance Assessment:

Technical Competency: The candidate demonstrated solid understanding of core concepts and provided well-structured responses. Their approach to problem-solving shows logical thinking and attention to detail.

Communication Skills: Responses were clear and articulate, showing good ability to explain complex technical concepts. The candidate maintained professional communication throughout.

Areas of Strength:
- Strong analytical thinking
- Good communication skills
- Relevant experience and knowledge
- Professional demeanor and attitude

Areas for Improvement:
- Could provide more specific examples
- Consider elaborating on implementation details
- Opportunity to discuss alternative approaches

Overall, this candidate shows promise and would benefit from continued development in the identified areas. The responses indicate readiness for the next level of technical challenges.`

	return &EvaluationResponse{
		OverallScore:    0.8,
		CategoryScores:  map[string]float64{"technical": 0.8, "communication": 0.85, "problem_solving": 0.75},
		Feedback:        feedback,
		Strengths:       []string{"[MOCK] Strong analytical thinking", "[MOCK] Clear communication", "[MOCK] Relevant experience"},
		Weaknesses:      []string{"[MOCK] Could provide more specific examples", "[MOCK] Opportunity for more detailed explanations"},
		Recommendations: []string{"[MOCK] Continue developing technical skills", "[MOCK] Practice explaining complex concepts", "[MOCK] Seek opportunities for hands-on experience"},
		TokensUsed:      TokenUsage{PromptTokens: 50, CompletionTokens: 150, TotalTokens: 200},
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
