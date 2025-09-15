package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectMySQL opens a pooled DB connection.
func ConnectMySQL(user, pass, host, port, dbname string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// Pool settings
	db.SetMaxOpenConns(25)           // maximum open connections
	db.SetMaxIdleConns(25)           // maximum idle connections
	db.SetConnMaxLifetime(time.Hour) // maximum connection lifetime
	db.SetConnMaxIdleTime(10 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
