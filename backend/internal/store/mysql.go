package store

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLDSN(user, pass, host, port, dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, dbname)
}

func ConnectMySQL(dsn string) (*sql.DB, error) {
	return sql.Open("mysql", dsn)
}
