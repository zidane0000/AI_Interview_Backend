package config_test

import (
	"os"
	"testing"

	"github.com/zidane0000/AI_Interview_Backend/config"
)

func TestLoadConfig_DefaultPort(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	os.Unsetenv("PORT")
	cfg := config.LoadConfig()
	if cfg.Port != "8080" {
		t.Errorf("expected default port 8080, got %s", cfg.Port)
	}
}

func TestLoadConfig_CustomPort(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	os.Setenv("PORT", "1234")
	cfg := config.LoadConfig()
	if cfg.Port != "1234" {
		t.Errorf("expected port 1234, got %s", cfg.Port)
	}
}

func TestLoadConfig_MissingDatabaseURL(t *testing.T) {
	os.Unsetenv("DATABASE_URL")
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected log.Fatal to exit, but it did not")
		}
	}()
	_ = config.LoadConfig()
}
