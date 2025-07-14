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

	// Example: create a user table in the new schema
	_, err = db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.user (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL,
			email TEXT NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`, schema))
	if err != nil {
		return fmt.Errorf("failed to create user table: %w", err)
	}

	// Add more table creation as needed

	return nil
}

// DeleteTenantSchema completely removes a tenant's schema and all its data
// WARNING: This operation is irreversible and will delete ALL tenant data
func DeleteTenantSchema(db *sql.DB, schemaName string) error {
	schema := pq.QuoteIdentifier(schemaName)

	// Drop the entire schema with CASCADE to remove all objects within it
	_, err := db.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schema))
	if err != nil {
		return fmt.Errorf("failed to drop schema %s: %w", schemaName, err)
	}

	return nil
}

// ListTenantSchemas returns all tenant schemas in the database
func ListTenantSchemas(db *sql.DB) ([]string, error) {
	query := `
		SELECT schema_name 
		FROM information_schema.schemata 
		WHERE schema_name NOT IN ('information_schema', 'pg_catalog', 'pg_toast', 'public')
		AND schema_name NOT LIKE 'pg_%'
		ORDER BY schema_name`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list schemas: %w", err)
	}
	defer rows.Close()

	var schemas []string
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			return nil, err
		}
		schemas = append(schemas, schema)
	}

	return schemas, nil
}
