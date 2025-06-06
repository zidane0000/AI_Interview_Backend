// Configuration loading (env, files)
package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string

	// TODO: Add AI service configuration
	// AIProvider   string // "openai", "anthropic", "google", "local"
	// AIAPIKey     string
	// AIBaseURL    string
	// AIModel      string
	// AITimeout    time.Duration

	// TODO: Add file upload configuration
	// UploadPath      string
	// MaxFileSize     int64
	// AllowedFileTypes []string

	// TODO: Add security configuration
	// JWTSecret       string
	// CORSOrigins     []string
	// RateLimitRPS    int
	// SessionTimeout  time.Duration

	// TODO: Add logging configuration
	// LogLevel        string
	// LogFormat       string // "json", "text"
	// LogOutput       string // "stdout", "file"

	// TODO: Add internationalization configuration
	// DefaultLanguage string
	// SupportedLangs  []string
	// TranslationPath string

	// TODO: Add email/notification configuration
	// SMTPHost     string
	// SMTPPort     int
	// SMTPUser     string
	// SMTPPassword string
	// EmailFrom    string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
	}

	// DATABASE_URL is optional with hybrid store architecture
	// - If present: Uses PostgreSQL database backend
	// - If absent: Uses in-memory store for development/testing

	if cfg.Port == "" {
		cfg.Port = "8080" // default port
	}

	// TODO: Load AI service configuration
	// cfg.AIProvider = getEnvWithDefault("AI_PROVIDER", "openai")
	// cfg.AIAPIKey = os.Getenv("AI_API_KEY")
	// cfg.AIBaseURL = getEnvWithDefault("AI_BASE_URL", "")
	// cfg.AIModel = getEnvWithDefault("AI_MODEL", "gpt-4")

	// TODO: Load file upload configuration
	// cfg.UploadPath = getEnvWithDefault("UPLOAD_PATH", "./uploads")
	// cfg.MaxFileSize = parseInt64WithDefault("MAX_FILE_SIZE", 10*1024*1024) // 10MB

	// TODO: Load security configuration
	// cfg.JWTSecret = os.Getenv("JWT_SECRET")
	// cfg.CORSOrigins = strings.Split(getEnvWithDefault("CORS_ORIGINS", "*"), ",")

	// TODO: Validate required AI configuration
	// if cfg.AIAPIKey == "" && cfg.AIProvider != "local" {
	//     return nil, errors.New("AI_API_KEY is required for external AI providers")
	// }

	// TODO: Validate file paths and create directories if needed
	// TODO: Validate email configuration if notifications are enabled
	// TODO: Load configuration from config files (YAML, JSON, TOML)
	// TODO: Add configuration hot-reloading capability
	// TODO: Add configuration validation with detailed error messages

	return cfg, nil
}

// TODO: Add helper functions for configuration parsing
// func getEnvWithDefault(key, defaultValue string) string {
//     if value := os.Getenv(key); value != "" {
//         return value
//     }
//     return defaultValue
// }

// TODO: Add configuration for different environments (dev, staging, prod)
// TODO: Add configuration documentation and examples
// TODO: Add configuration schema validation
// TODO: Add sensitive data masking in logs
