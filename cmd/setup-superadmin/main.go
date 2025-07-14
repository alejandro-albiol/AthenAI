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
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load environment variables from .env file
	config.LoadEnv()

	fmt.Println("=== AthenAI Super Admin Setup ===")
	fmt.Println("This tool creates platform super administrators.")
	fmt.Println("Use with extreme caution - only run in secure environments.")
	fmt.Println()

	// Connect to database
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Check if this is the first super admin
	count, err := countExistingSuperAdmins(db)
	if err != nil {
		log.Fatal("Failed to check existing admins:", err)
	}

	if count > 0 {
		fmt.Printf("Warning: %d super admin(s) already exist.\n", count)
		if !confirmAction("Continue creating another super admin?") {
			fmt.Println("Operation cancelled.")
			return
		}
	}

	// Get admin details
	username := getInput("Username: ")
	email := getInput("Email: ")
	password := getPasswordInput("Password: ")
	confirmPassword := getPasswordInput("Confirm password: ")

	if password != confirmPassword {
		log.Fatal("Passwords do not match!")
	}

	// Validate inputs
	if len(username) < 3 {
		log.Fatal("Username must be at least 3 characters")
	}
	if len(password) < 8 {
		log.Fatal("Password must be at least 8 characters")
	}
	if !strings.Contains(email, "@") {
		log.Fatal("Invalid email format")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// Final confirmation
	fmt.Printf("\nCreating super admin:\n")
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Email: %s\n", email)
	fmt.Println()

	if !confirmAction("Create this super admin?") {
		fmt.Println("Operation cancelled.")
		return
	}

	// Create super admin
	err = createSuperAdmin(db, username, email, string(hashedPassword))
	if err != nil {
		log.Fatal("Failed to create super admin:", err)
	}

	fmt.Println()
	fmt.Println("✅ Super admin created successfully!")
	fmt.Println("⚠️  Remember to:")
	fmt.Println("   - Store credentials securely")
	fmt.Println("   - Remove this setup tool from production")
	fmt.Println("   - Monitor admin access logs")
}

func countExistingSuperAdmins(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM public.admin WHERE is_active = true").Scan(&count)
	return count, err
}

func createSuperAdmin(db *sql.DB, username, email, passwordHash string) error {
	query := `
		INSERT INTO public.admin (id, username, password_hash, email, is_active, created_at, updated_at) 
		VALUES (gen_random_uuid(), $1, $2, $3, true, NOW(), NOW())
	`
	_, err := db.Exec(query, username, passwordHash, email)
	return err
}

func getInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func getPasswordInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Failed to read password:", err)
	}
	return strings.TrimSpace(password)
}

func confirmAction(prompt string) bool {
	fmt.Printf("%s (y/N): ", prompt)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}
