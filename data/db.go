// Database connection and setup using GORM
package data

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the PostgreSQL connection using GORM
func InitDB(databaseURL string) {
	var err error
	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
}

// CloseDB closes the database connection (if needed)
func CloseDB() {
	db, err := DB.DB()
	if err == nil {
		db.Close()
	}
}
