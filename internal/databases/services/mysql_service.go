package services

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

type MySQLService struct {
    dsn string
    db  *sql.DB
}

func NewMySQLService(dsn string) *MySQLService {
    return &MySQLService{dsn: dsn}
}

func (m *MySQLService) Connect() error {
    db, err := sql.Open("mysql", m.dsn)
    if err != nil {
        return err
    }
    m.db = db
    return nil
}

func (m *MySQLService) Close() error {
    return m.db.Close()
}

func (m *MySQLService) Exec(query string, args ...interface{}) (sql.Result, error) {
    return m.db.Exec(query, args...)
}

func (m *MySQLService) Query(query string, args ...interface{}) (*sql.Rows, error) {
    return m.db.Query(query, args...)
}