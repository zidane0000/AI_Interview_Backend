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

// Global AI client instance (automatically loads .env configuration)
// TODO: ARCHITECTURAL CONCERN - Global client design issues:
// - Testing difficulties: Hard to mock or inject test doubles
// - Tight coupling: All code tightly coupled to single global instance
// - Configuration inflexibility: Can't have different configs for different use cases
// - Concurrency concerns: Global state can lead to race conditions
// - Initialization order dependencies: Must be initialized before use
//
// FUTURE REFACTORING PLAN:
// 1. Introduce dependency injection in handlers
// 2. Make the global client optional/deprecated
// 3. Improve testability with proper mocking
// 4. Consider service container or context-based injection
//
// For now keeping global client to focus on fixing failing tests first.
var Client = NewAutoAIClient()

// NewAutoAIClient initializes the AI client using the best available API key (OpenAI > Gemini > none)
// This method automatically loads .env files and reads environment variables
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

	return &AIClient{
		enhancedClient: NewEnhancedAIClient(config),
	}
}

// GenerateChatResponse generates AI response for conversational interviews
func (c *AIClient) GenerateChatResponse(sessionID string, conversationHistory []map[string]string, userMessage string) (string, error) {
	return c.GenerateChatResponseWithLanguage(sessionID, conversationHistory, userMessage, "en")
}

// GenerateChatResponseWithLanguage generates AI response with language support
func (c *AIClient) GenerateChatResponseWithLanguage(sessionID string, conversationHistory []map[string]string, userMessage string, language string) (string, error) {
	// Build context for the AI including conversation history and language
	contextMap := map[string]interface{}{
		"interview_type":       "general",
		"job_title":            "Software Engineer",
		"context":              "Interview in progress",
		"conversation_history": conversationHistory,
		"language":             language,
	}

	return c.enhancedClient.GenerateInterviewResponse(sessionID, userMessage, contextMap)
}

// GenerateClosingMessage generates a closing AI response for ending interviews
func (c *AIClient) GenerateClosingMessage(sessionID string, conversationHistory []map[string]string, userMessage string) (string, error) {
	return c.GenerateClosingMessageWithLanguage(sessionID, conversationHistory, userMessage, "en")
}

// GenerateClosingMessageWithLanguage generates a closing AI response with language support
func (c *AIClient) GenerateClosingMessageWithLanguage(sessionID string, conversationHistory []map[string]string, userMessage string, language string) (string, error) {
	// Build context for the AI to indicate this is the final message
	contextMap := map[string]interface{}{
		"interview_type":       "general",
		"job_title":            "Software Engineer",
		"context":              "This is the final message - wrap up the interview professionally and thank the candidate",
		"conversation_history": conversationHistory,
		"closing_interview":    true,
		"language":             language,
	}

	return c.enhancedClient.GenerateInterviewResponse(sessionID, userMessage, contextMap)
}

// ShouldEndInterview determines if the interview should end
func (c *AIClient) ShouldEndInterview(messageCount int) bool {
	return messageCount >= 8 // End after 8 user messages
}

// EvaluateAnswers evaluates chat conversation and generates score and feedback
func (c *AIClient) EvaluateAnswers(questions []string, answers []string, language string) (float64, string, error) {
	// Use the context version with default job info
	return c.EvaluateAnswersWithContext(questions, answers, "Software Engineer", "General interview evaluation", language)
}

// EvaluateAnswersWithContext evaluates chat conversation with interview context
func (c *AIClient) EvaluateAnswersWithContext(questions []string, answers []string, jobTitle, jobDesc, language string) (float64, string, error) {
	if len(answers) == 0 {
		return 0.0, "No answers provided.", nil
	}

	// Use the enhanced AI client for real evaluation with context
	ctx := context.Background()

	// Create evaluation request with proper context including language
	req := &EvaluationRequest{
		Questions:   questions,
		Answers:     answers,
		JobTitle:    jobTitle,
		JobDesc:     jobDesc,
		Criteria:    []string{"communication", "technical_knowledge", "problem_solving", "clarity", "cultural_fit"},
		DetailLevel: "detailed",
		Language:    language, // Pass language for evaluation
		Context: map[string]interface{}{
			"interview_type":  "conversational",
			"evaluation_type": "chat_based",
			"language":        language, // Also include in context map
		},
	}

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

// GetCurrentProvider returns the currently configured AI provider
func (c *AIClient) GetCurrentProvider() string {
	return c.enhancedClient.config.DefaultProvider
}

// GetCurrentModel returns the currently configured AI model
func (c *AIClient) GetCurrentModel() string {
	return c.enhancedClient.config.DefaultModel
}
