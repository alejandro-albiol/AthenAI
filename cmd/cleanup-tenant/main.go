package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alejandro-albiol/athenai/config"
	"github.com/alejandro-albiol/athenai/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("🗄️  AthenAI Tenant Schema Cleanup Tool")
	fmt.Println("=====================================")
	fmt.Println("⚠️  WARNING: This will permanently delete tenant data!")
	fmt.Println()

	// List existing tenant schemas
	schemas, err := database.ListTenantSchemas(db)
	if err != nil {
		log.Fatal("Failed to list schemas:", err)
	}

	if len(schemas) == 0 {
		fmt.Println("✅ No tenant schemas found.")
		return
	}

	fmt.Println("📋 Found tenant schemas:")
	for i, schema := range schemas {
		fmt.Printf("  %d. %s\n", i+1, schema)
	}
	fmt.Println()

	// Get user input
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the schema name to delete (or 'quit' to exit): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Failed to read input:", err)
	}

	schemaName := strings.TrimSpace(input)

	if schemaName == "quit" || schemaName == "" {
		fmt.Println("👋 Cancelled.")
		return
	}

	// Verify schema exists
	found := false
	for _, schema := range schemas {
		if schema == schemaName {
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("❌ Schema '%s' not found.\n", schemaName)
		return
	}

	// Final confirmation
	fmt.Printf("⚠️  You are about to PERMANENTLY DELETE schema '%s' and ALL its data.\n", schemaName)
	fmt.Print("Type 'DELETE' to confirm: ")

	confirmation, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Failed to read confirmation:", err)
	}

	if strings.TrimSpace(confirmation) != "DELETE" {
		fmt.Println("👋 Cancelled.")
		return
	}

	// Delete the schema
	fmt.Printf("🗑️  Deleting schema '%s'...\n", schemaName)

	err = database.DeleteTenantSchema(db, schemaName)
	if err != nil {
		log.Fatal("Failed to delete schema:", err)
	}

	fmt.Printf("✅ Schema '%s' has been successfully deleted.\n", schemaName)

	// Show remaining schemas
	remainingSchemas, err := database.ListTenantSchemas(db)
	if err != nil {
		log.Printf("Warning: Failed to list remaining schemas: %v", err)
		return
	}

	if len(remainingSchemas) > 0 {
		fmt.Println("\n📋 Remaining tenant schemas:")
		for _, schema := range remainingSchemas {
			fmt.Printf("  - %s\n", schema)
		}
	} else {
		fmt.Println("\n✅ No tenant schemas remaining.")
	}
}
