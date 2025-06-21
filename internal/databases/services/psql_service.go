package services

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type PostgresService struct {
	dsn string
	db  *sql.DB
}

func NewPostgresService(dsn string) *PostgresService {
	return &PostgresService{dsn: dsn}
}

func (p *PostgresService) Connect() error {
	db, err := sql.Open("postgres", p.dsn)
	if err != nil {
		return fmt.Errorf("failed to open postgres connection: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping postgres: %w", err)
	}
	p.db = db
	return nil
}

func (p *PostgresService) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func (p *PostgresService) Exec(query string, args ...interface{}) (sql.Result, error) {
	if p.db == nil {
		return nil, fmt.Errorf("database not connected")
	}
	return p.db.Exec(query, args...)
}

func (p *PostgresService) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if p.db == nil {
		return nil, fmt.Errorf("database not connected")
	}
	return p.db.Query(query, args...)
}