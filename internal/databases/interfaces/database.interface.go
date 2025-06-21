package interfaces

import "database/sql"

type DBService interface {
    Connect() error
    Close() error
    Exec(query string, args ...interface{}) (sql.Result, error)
    Query(query string, args ...interface{}) (*sql.Rows, error)
}