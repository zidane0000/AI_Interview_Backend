// Database connection and setup using GORM
package data

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes and returns a PostgreSQL connection using GORM
func InitDB(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CloseDB closes the provided database connection (if needed)
func CloseDB(db *gorm.DB) {
	if db == nil {
		return
	}
	dbConn, err := db.DB()
	if err == nil {
		dbConn.Close()
	}
}
