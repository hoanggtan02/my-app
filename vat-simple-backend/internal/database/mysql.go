package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/config"
)

// InitMySQL initializes and returns a connection to the MySQL database.
func InitMySQL(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// CloseMySQL closes the database connection.
func CloseMySQL(db *sql.DB) {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Error closing the database: %v", err)
		} else {
			log.Println("Database connection closed.")
		}
	}
}
