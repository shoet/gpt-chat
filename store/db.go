package store

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func NewRDB(dsn string) (*sqlx.DB, func() error, error) {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, conn.Close, fmt.Errorf("failed to open db: %w", err)
	}
	if err := conn.Ping(); err != nil {
		return nil, conn.Close, fmt.Errorf("failed to ping db: %w", err)
	}
	xdb := sqlx.NewDb(conn, "mysql")
	return xdb, conn.Close, nil
}
