package database

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

// CreateTenantSchema creates a new schema and users table for a tenant (gym)
func CreateTenantSchema(db *sql.DB, schemaName string) error {
	schema := pq.QuoteIdentifier(schemaName)
	_, err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema))
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// Example: create a users table in the new schema
	_, err = db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.users (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL,
			email TEXT NOT NULL,
			is_verified BOOLEAN NOT NULL DEFAULT false,
			is_active BOOLEAN NOT NULL DEFAULT true,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`, schema))
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Add more table creation as needed

	return nil
}
