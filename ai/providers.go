// Provider registry and initialization utilities
package ai

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsIntOrDefault returns environment variable as int or default
func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getEnvAsFloatOrDefault returns environment variable as float64 or default
func getEnvAsFloatOrDefault(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return defaultValue
}

// getEnvAsBoolOrDefault returns environment variable as bool or default
func getEnvAsBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

// getEnvAsDurationOrDefault returns environment variable as duration or default
func getEnvAsDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// NewDefaultAIConfig creates a default AI configuration from environment variables
// This function automatically loads .env files to ensure configuration is available
func NewDefaultAIConfig() *AIConfig {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	return &AIConfig{
		OpenAIAPIKey:     getEnvOrDefault("OPENAI_API_KEY", ""),
		GeminiAPIKey:     getEnvOrDefault("GEMINI_API_KEY", ""),
		DefaultProvider:  getEnvOrDefault("AI_DEFAULT_PROVIDER", ProviderMock),
		DefaultModel:     getEnvOrDefault("AI_DEFAULT_MODEL", "mock-model"),
		MaxRetries:       getEnvAsIntOrDefault("AI_MAX_RETRIES", 3),
		RequestTimeout:   getEnvAsDurationOrDefault("AI_REQUEST_TIMEOUT", 60*time.Second),
		DefaultMaxTokens: getEnvAsIntOrDefault("AI_DEFAULT_MAX_TOKENS", 1000),
		DefaultTemp:      getEnvAsFloatOrDefault("AI_DEFAULT_TEMPERATURE", 0.7),
		EnableCaching:    getEnvAsBoolOrDefault("AI_ENABLE_CACHING", true),
		EnableMetrics:    getEnvAsBoolOrDefault("AI_ENABLE_METRICS", true),
		EnableStreaming:  getEnvAsBoolOrDefault("AI_ENABLE_STREAMING", false),
		RateLimitRPM:     getEnvAsIntOrDefault("AI_RATE_LIMIT_RPM", 60),
		RateLimitTPM:     getEnvAsIntOrDefault("AI_RATE_LIMIT_TPM", 60000),
		DailyTokenLimit:  getEnvAsIntOrDefault("AI_DAILY_TOKEN_LIMIT", 100000),
		CostPerToken:     getEnvAsFloatOrDefault("AI_COST_PER_TOKEN", 0.000002),
		MaxCostPerDay:    getEnvAsFloatOrDefault("AI_MAX_COST_PER_DAY", 10.0),
	}
}

// ValidateConfig validates the AI configuration
func ValidateConfig(config *AIConfig) error {
	if config.OpenAIAPIKey == "" && config.GeminiAPIKey == "" && config.DefaultProvider != ProviderMock {
		return fmt.Errorf("at least one AI provider API key must be configured, or use mock provider")
	}

	if config.DefaultProvider != ProviderOpenAI && config.DefaultProvider != ProviderGemini && config.DefaultProvider != ProviderMock {
		return fmt.Errorf("invalid default provider: %s", config.DefaultProvider)
	}

	if config.DefaultProvider == ProviderOpenAI && config.OpenAIAPIKey == "" {
		return fmt.Errorf("OpenAI API key required when using OpenAI as default provider")
	}

	if config.DefaultProvider == ProviderGemini && config.GeminiAPIKey == "" {
		return fmt.Errorf("Gemini API key required when using Gemini as default provider")
	}

	// Mock provider doesn't require API keys

	if config.MaxRetries < 0 {
		return fmt.Errorf("max retries cannot be negative")
	}

	if config.RequestTimeout <= 0 {
		return fmt.Errorf("request timeout must be positive")
	}

	if config.DefaultMaxTokens <= 0 {
		return fmt.Errorf("default max tokens must be positive")
	}

	if config.DefaultTemp < 0 || config.DefaultTemp > 2 {
		return fmt.Errorf("default temperature must be between 0 and 2")
	}

	return nil
}

// GetAvailableProviders returns list of providers with valid API keys
func GetAvailableProviders(config *AIConfig) []string {
	var providers []string

	if config.OpenAIAPIKey != "" {
		providers = append(providers, ProviderOpenAI)
	}

	if config.GeminiAPIKey != "" {
		providers = append(providers, ProviderGemini)
	}

	// Mock provider is always available
	providers = append(providers, ProviderMock)

	return providers
}

// GetProviderInfo returns information about a specific provider
func GetProviderInfo(provider string) map[string]interface{} {
	switch provider {
	case ProviderOpenAI:
		return map[string]interface{}{
			"name":               "OpenAI",
			"models":             []string{"gpt-4", "gpt-4-turbo", "gpt-3.5-turbo"},
			"supports_vision":    true,
			"supports_functions": true,
			"max_tokens":         4096,
			"website":            "https://platform.openai.com/",
		}
	case ProviderGemini:
		return map[string]interface{}{
			"name":               "Google Gemini",
			"models":             []string{"gemini-1.5-pro", "gemini-1.5-flash", "gemini-pro"},
			"supports_vision":    true,
			"supports_functions": true,
			"max_tokens":         8192,
			"website":            "https://ai.google.dev/gemini-api",
		}
	case ProviderMock:
		return map[string]interface{}{
			"name":               "Mock Provider",
			"models":             []string{"mock-model"},
			"supports_vision":    false,
			"supports_functions": false,
			"max_tokens":         1000,
			"website":            "https://localhost/mock",
		}
	default:
		return map[string]interface{}{
			"error": "Unknown provider",
		}
	}
}

// CreateAIProviderFromConfig creates an AI provider instance from configuration
func CreateAIProviderFromConfig(providerName string, config *AIConfig) (AIProvider, error) {
	switch providerName {
	case ProviderOpenAI:
		if config.OpenAIAPIKey == "" {
			return nil, fmt.Errorf("OpenAI API key not configured")
		}
		return NewOpenAIProvider(config.OpenAIAPIKey, config), nil
	case ProviderGemini:
		if config.GeminiAPIKey == "" {
			return nil, fmt.Errorf("Gemini API key not configured")
		}
		return NewGeminiProvider(config.GeminiAPIKey, config), nil
	case ProviderMock:
		return NewMockProvider(), nil

	default:
		return nil, fmt.Errorf("unknown provider: %s", providerName)
	}
}

// GetRecommendedProvider returns the recommended provider based on task type
func GetRecommendedProvider(taskType string, availableProviders []string) string {
	if len(availableProviders) == 0 {
		return ""
	}

	// Default to first available provider
	defaultProvider := availableProviders[0]

	switch taskType {
	case "chat", "conversation":
		// Both are good for chat, prefer OpenAI for consistency
		for _, provider := range availableProviders {
			if provider == ProviderOpenAI {
				return provider
			}
		}
	case "evaluation", "analysis":
		// Gemini might be good for analytical tasks
		for _, provider := range availableProviders {
			if provider == ProviderGemini {
				return provider
			}
		}
	case "question_generation":
		// Both are suitable, prefer whichever is available
		return defaultProvider
	}

	return defaultProvider
}

// GetModelRecommendation returns recommended model for a provider and task
func GetModelRecommendation(provider, taskType string) string {
	switch provider {
	case ProviderOpenAI:
		switch taskType {
		case "chat", "conversation":
			return "gpt-3.5-turbo" // Fast and cost-effective for chat
		case "evaluation", "analysis":
			return "gpt-4" // More accurate for complex analysis
		case "question_generation":
			return "gpt-3.5-turbo" // Good balance for question generation
		default:
			return "gpt-3.5-turbo"
		}
	case ProviderGemini:
		switch taskType {
		case "chat", "conversation":
			return "gemini-1.5-flash" // Fast responses for chat
		case "evaluation", "analysis":
			return "gemini-1.5-pro" // Better for complex reasoning
		case "question_generation":
			return "gemini-1.5-flash" // Good for generation tasks
		default:
			return "gemini-1.5-flash"
		}
	case ProviderMock:
		return "mock-model" // Mock provider always uses mock-model
	default:
		return ""
	}
}
