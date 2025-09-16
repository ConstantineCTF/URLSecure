package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver imported for side effects
)

// ConnectMySQL opens a pooled MySQL DB connection based on provided credentials.
func ConnectMySQL(user, pass, host, port, dbname string) (*sql.DB, error) {
	// Format MySQL DSN (data source name) with parseTime=true to handle time fields
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, dbname)

	// Open connection to MySQL database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Configure connection pool settings for performance and resource management
	db.SetMaxOpenConns(25)            // Max open connections
	db.SetMaxIdleConns(25)            // Max idle connections for reuse
	db.SetConnMaxLifetime(time.Hour) // Max lifetime before recycling connection
	db.SetConnMaxIdleTime(10 * time.Minute) // Max idle time before recycling

	// Ping DB to verify connectivity
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
