// Client for communicating with AI service/model
package ai

import (
	"context"
	"fmt"
)

// Legacy AIClient for backward compatibility
type AIClient struct {
	enhancedClient *EnhancedAIClient
}

// Global AI client instance
var Client = NewAutoAIClient()

// NewAutoAIClient initializes the AI client using the best available API key (OpenAI > Gemini > none)
func NewAutoAIClient() *AIClient {
	config := NewDefaultAIConfig() // loads from env

	// Priority: OpenAI > Gemini > fallback
	if config.OpenAIAPIKey != "" {
		config.DefaultProvider = ProviderOpenAI
		config.DefaultModel = "gpt-4o"

	} else if config.GeminiAPIKey != "" {
		config.DefaultProvider = ProviderGemini
		config.DefaultModel = "gemini-2.0-flash"
	}
	fmt.Printf("Using AI provider: %s and model: %s\n", config.DefaultProvider, config.DefaultModel)

	return &AIClient{
		enhancedClient: NewEnhancedAIClient(config),
	}
}

// GenerateChatResponse generates AI response for conversational interviews
func (c *AIClient) GenerateChatResponse(conversationHistory []map[string]string, userMessage string) (string, error) {
	sessionID := "default-session" // TODO: Use proper session ID from context

	// Build context for the AI including conversation history
	contextMap := map[string]interface{}{
		"interview_type":       "general",
		"job_title":            "Software Engineer",
		"context":              "Interview in progress",
		"conversation_history": conversationHistory,
	}

	return c.enhancedClient.GenerateInterviewResponse(sessionID, userMessage, contextMap)
}

// GenerateClosingMessage generates a closing AI response for ending interviews
func (c *AIClient) GenerateClosingMessage(conversationHistory []map[string]string, userMessage string) (string, error) {
	sessionID := "default-session" // TODO: Use proper session ID from context

	// Build context for the AI to indicate this is the final message
	contextMap := map[string]interface{}{
		"interview_type":       "general",
		"job_title":            "Software Engineer",
		"context":              "This is the final message - wrap up the interview professionally and thank the candidate",
		"conversation_history": conversationHistory,
		"closing_interview":    true,
	}

	return c.enhancedClient.GenerateInterviewResponse(sessionID, userMessage, contextMap)
}

// ShouldEndInterview determines if the interview should end
func (c *AIClient) ShouldEndInterview(messageCount int) bool {
	return messageCount >= 8 // End after 8 user messages
}

// EvaluateAnswers evaluates chat conversation and generates score and feedback
func (c *AIClient) EvaluateAnswers(questions []string, answers []string) (float64, string, error) {
	if len(answers) == 0 {
		return 0.0, "No answers provided.", nil
	}

	// Use the enhanced AI client for real evaluation
	ctx := context.Background()

	// Create evaluation request
	req := &EvaluationRequest{
		Questions:   questions,
		Answers:     answers,
		JobTitle:    "Software Engineer", // Default job title
		JobDesc:     "General software engineering position",
		Criteria:    []string{"communication", "technical_knowledge", "problem_solving", "clarity"},
		DetailLevel: "detailed",
		Context: map[string]interface{}{
			"interview_type":  "conversational",
			"evaluation_type": "chat_based",
		}}

	// Call enhanced client for evaluation
	resp, err := c.enhancedClient.EvaluateAnswers(ctx, req)
	if err != nil {
		return 0.0, "Evaluation failed", err
	}

	return resp.OverallScore, resp.Feedback, nil
}

// EvaluateAnswersWithContext evaluates chat conversation with interview context
func (c *AIClient) EvaluateAnswersWithContext(questions []string, answers []string, jobTitle, jobDesc string) (float64, string, error) {
	if len(answers) == 0 {
		return 0.0, "No answers provided.", nil
	}

	// Use the enhanced AI client for real evaluation with context
	ctx := context.Background()

	// Create evaluation request with proper context
	req := &EvaluationRequest{
		Questions:   questions,
		Answers:     answers,
		JobTitle:    jobTitle,
		JobDesc:     jobDesc,
		Criteria:    []string{"communication", "technical_knowledge", "problem_solving", "clarity", "cultural_fit"},
		DetailLevel: "detailed",
		Context: map[string]interface{}{
			"interview_type":  "conversational",
			"evaluation_type": "chat_based",
		}}

	// Call enhanced client for evaluation
	resp, err := c.enhancedClient.EvaluateAnswers(ctx, req)
	if err != nil {
		return 0.0, "Evaluation failed", err
	}

	return resp.OverallScore, resp.Feedback, nil
}

// GenerateQuestionsFromResume generates interview questions based on resume and job description
func (c *AIClient) GenerateQuestionsFromResume(resumeText, jobDescription, jobTitle string) ([]InterviewQuestion, error) {
	ctx := context.Background()

	req := &QuestionGenerationRequest{
		JobTitle:        jobTitle,
		JobDescription:  jobDescription,
		ResumeContent:   resumeText,
		InterviewType:   "mixed",
		NumQuestions:    8,
		ExperienceLevel: "mid",
		Difficulty:      "medium",
	}

	resp, err := c.enhancedClient.GenerateQuestions(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Questions, nil
}

// GenerateInterviewQuestions generates questions for a specific interview setup
func (c *AIClient) GenerateInterviewQuestions(jobTitle, jobDesc string, questionCount int) ([]InterviewQuestion, error) {
	ctx := context.Background()

	req := &QuestionGenerationRequest{
		JobTitle:        jobTitle,
		JobDescription:  jobDesc,
		InterviewType:   "general",
		NumQuestions:    questionCount,
		ExperienceLevel: "mid",
		Difficulty:      "medium",
	}

	resp, err := c.enhancedClient.GenerateQuestions(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Questions, nil
}

// GetProviderInfo returns information about available AI providers
func (c *AIClient) GetProviderInfo() map[string]interface{} {
	info := make(map[string]interface{})
	providers := c.enhancedClient.GetAvailableProviders()

	for _, providerName := range providers {
		info[providerName] = GetProviderInfo(providerName)
	}

	return info
}

// SwitchProvider changes the active AI provider
func (c *AIClient) SwitchProvider(providerName string) error {
	c.enhancedClient.mu.Lock()
	defer c.enhancedClient.mu.Unlock()

	if _, exists := c.enhancedClient.providers[providerName]; !exists {
		return fmt.Errorf("provider not available: %s", providerName)
	}

	c.enhancedClient.config.DefaultProvider = providerName
	return nil
}
