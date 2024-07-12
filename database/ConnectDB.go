package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type Database struct {
	User     string
	Password string
	Host     string
	Name     string
}

func (db *Database) Connect() (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", db.User, db.Password, db.Host, db.Name)
	conn, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error opening MySQL database connection: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging MySQL database: %w", err)
	}

	log.Println("Connected to MySQL database Successfully!")

	// Setting connection pool
	conn.SetMaxOpenConns(10)                 // number connection max open
	conn.SetMaxIdleConns(8)                  // number connection free
	conn.SetConnMaxLifetime(1 * time.Second) // time connection max can use
	return conn, nil
}
