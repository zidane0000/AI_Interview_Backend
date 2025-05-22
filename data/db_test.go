package data_test

import (
	"os"
	"testing"

	"github.com/zidane0000/AI_Interview_Backend/config"
	"github.com/zidane0000/AI_Interview_Backend/data"
)

func TestInitAndCloseDB(t *testing.T) {
	// Set a fake database URL for testing (should fail to connect)
	os.Setenv("DATABASE_URL", "postgres://invalid:invalid@localhost:5432/invalid_db")
	cfg := config.LoadConfig()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic due to invalid database URL, but did not panic")
		}
	}()

	// This should log.Fatal (panic) due to invalid connection
	data.InitDB(cfg.DatabaseURL)
}
