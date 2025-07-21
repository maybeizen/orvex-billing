package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/google/uuid"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database management commands",
	Long:  `Commands for managing the database.`,
}

var wipeCmd = &cobra.Command{
	Use:   "wipe",
	Short: "Wipe the entire database",
	Long:  `Delete all data from the database. This action cannot be undone!`,
	Run: func(cmd *cobra.Command, args []string) {
		// Confirm dangerous operation
		fmt.Println("⚠️  WARNING: This will delete ALL data from the database!")
		fmt.Print("Type 'WIPE' to confirm: ")
		
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		if input != "WIPE" {
			fmt.Println("Database wipe cancelled")
			return
		}

		if err := wipeDatabase(); err != nil {
			fmt.Printf("Error wiping database: %v\n", err)
			return
		}
		
		fmt.Println("Database wiped successfully")
	},
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with initial data",
	Long:  `Add initial data to the database (admin user, sample data, etc.).`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := seedDatabase(); err != nil {
			fmt.Printf("Error seeding database: %v\n", err)
			return
		}
		fmt.Println("Database seeded successfully")
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the database (wipe + migrate + seed)",
	Long:  `Completely reset the database by wiping all data, running migrations, and seeding initial data.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Confirm dangerous operation
		fmt.Println("⚠️  WARNING: This will delete ALL data and reset the database!")
		fmt.Print("Type 'RESET' to confirm: ")
		
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		if input != "RESET" {
			fmt.Println("Database reset cancelled")
			return
		}

		fmt.Println("Wiping database...")
		if err := wipeDatabase(); err != nil {
			fmt.Printf("Error wiping database: %v\n", err)
			return
		}

		fmt.Println("Running migrations...")
		if err := runMigrations(); err != nil {
			fmt.Printf("Error running migrations: %v\n", err)
			return
		}

		fmt.Println("Seeding database...")
		if err := seedDatabase(); err != nil {
			fmt.Printf("Error seeding database: %v\n", err)
			return
		}
		
		fmt.Println("Database reset completed successfully")
	},
}

func wipeDatabase() error {
	// Wipe sessions database first
	if err := wipeSessionsDatabase(); err != nil {
		return fmt.Errorf("failed to wipe sessions database: %v", err)
	}

	// Get all table names from main database
	var tableNames []string
	
	// For SQLite, query sqlite_master table
	rows, err := db.DB.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").Rows()
	if err != nil {
		return fmt.Errorf("failed to get table names: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("failed to scan table name: %v", err)
		}
		tableNames = append(tableNames, tableName)
	}

	// Disable foreign key constraints temporarily
	if err := db.DB.Exec("PRAGMA foreign_keys = OFF").Error; err != nil {
		return fmt.Errorf("failed to disable foreign keys: %v", err)
	}

	// Drop all tables
	for _, tableName := range tableNames {
		if err := db.DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)).Error; err != nil {
			return fmt.Errorf("failed to drop table %s: %v", tableName, err)
		}
	}

	// Re-enable foreign key constraints
	if err := db.DB.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return fmt.Errorf("failed to re-enable foreign keys: %v", err)
	}

	// Recreate tables using AutoMigrate
	err = db.DB.AutoMigrate(
		&types.User{}, 
		&types.Session{}, 
		&types.BackupCode{},
		&types.Service{}, 
		&types.Invoice{}, 
		&types.InvoiceItem{},
		&types.Category{},
		&types.Product{},
		&types.Cart{},
		&types.CartItem{},
		&types.Coupon{},
		&types.CouponUsage{},
	)
	if err != nil {
		return fmt.Errorf("failed to recreate tables: %v", err)
	}

	return nil
}

func wipeSessionsDatabase() error {
	sessionsDbPath := "./database/sessions.db"
	
	// Check if sessions database exists
	if _, err := os.Stat(sessionsDbPath); os.IsNotExist(err) {
		// Sessions database doesn't exist, nothing to wipe
		return nil
	}

	// Remove the sessions database file
	if err := os.Remove(sessionsDbPath); err != nil {
		return fmt.Errorf("failed to remove sessions database: %v", err)
	}

	fmt.Println("Wiped sessions database")
	return nil
}

func seedDatabase() error {
	// Check if we should seed an admin user
	var userCount int64
	db.DB.Model(&types.User{}).Count(&userCount)
	
	if userCount == 0 {
		fmt.Print("Create admin user? (Y/n): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		
		if input == "" || input == "y" || input == "yes" {
			if err := createAdminUser(); err != nil {
				return fmt.Errorf("failed to create admin user: %v", err)
			}
		}
	}

	// Add other seed data here if needed
	
	return nil
}

func createAdminUser() error {
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Print("Admin email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)
	if email == "" {
		return fmt.Errorf("email is required")
	}

	fmt.Print("Admin username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	if username == "" {
		return fmt.Errorf("username is required")
	}

	fmt.Print("Admin password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("error reading password: %v", err)
	}
	fmt.Println() // New line after password input
	password := strings.TrimSpace(string(passwordBytes))
	if password == "" {
		return fmt.Errorf("password is required")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	// Create admin user
	adminUser := types.User{
		UUID:         uuid.New(),
		Email:        email,
		Username:     username,
		PasswordHash: string(hashedPassword),
		FirstName:    "Admin",
		LastName:     "User",
		Role:         types.UserRoleAdmin,
	}

	if err := db.DB.Create(&adminUser).Error; err != nil {
		return fmt.Errorf("error creating admin user: %v", err)
	}

	fmt.Printf("Created admin user: %s (%s)\n", username, email)
	return nil
}

func init() {
	dbCmd.AddCommand(wipeCmd)
	dbCmd.AddCommand(seedCmd)
	dbCmd.AddCommand(resetCmd)
}