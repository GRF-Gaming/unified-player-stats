package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

func NewPgConn(
	addr string,
	port int,
	password string,
	dbName string,
	maxActiveConns int,
) (*sqlx.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s port=%d user=postgres password=%s dbname=%s sslmode=disable",
		addr,
		port,
		password,
		dbName,
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxActiveConns)

	// verify connection
	if err := db.Ping(); err != nil {
		slog.Error("Failed to ping pg db")
		return nil, err
	}

	return db, nil
}
