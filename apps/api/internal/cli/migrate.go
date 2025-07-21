package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/orvexcc/billing/api/internal/db"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
	Long:  `Commands for managing database migrations.`,
}

var runMigrationsCmd = &cobra.Command{
	Use:   "run",
	Short: "Run pending database migrations",
	Long:  `Execute all pending database migrations in chronological order.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := runMigrations()
		if err != nil {
			fmt.Printf("Error running migrations: %v\n", err)
			return
		}
		fmt.Println("Migrations completed successfully")
	},
}

var createMigrationCmd = &cobra.Command{
	Use:   "create [migration_name]",
	Short: "Create a new migration file",
	Long:  `Create a new migration file with a timestamp prefix.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		migrationName := args[0]
		err := createMigrationFile(migrationName)
		if err != nil {
			fmt.Printf("Error creating migration file: %v\n", err)
			return
		}
		fmt.Printf("Migration file created successfully\n")
	},
}

var migrationStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Long:  `Display the status of all migrations (pending/completed).`,
	Run: func(cmd *cobra.Command, args []string) {
		err := showMigrationStatus()
		if err != nil {
			fmt.Printf("Error showing migration status: %v\n", err)
			return
		}
	},
}

// Migration represents a database migration record
type Migration struct {
	ID        uint      `gorm:"primaryKey"`
	Filename  string    `gorm:"uniqueIndex;not null"`
	AppliedAt time.Time `gorm:"not null"`
}

// Ensure migrations table exists
func ensureMigrationsTable() error {
	return db.DB.AutoMigrate(&Migration{})
}

func runMigrations() error {
	// Ensure migrations table exists
	if err := ensureMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Get migration directory path
	migrationsDir := "./database/migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(migrationsDir, 0755); err != nil {
			return fmt.Errorf("failed to create migrations directory: %v", err)
		}
		fmt.Println("No migrations found")
		return nil
	}

	// Read migration files
	migrationFiles, err := getMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migration files: %v", err)
	}

	if len(migrationFiles) == 0 {
		fmt.Println("No migration files found")
		return nil
	}

	// Get applied migrations
	var appliedMigrations []Migration
	db.DB.Find(&appliedMigrations)

	appliedMap := make(map[string]bool)
	for _, migration := range appliedMigrations {
		appliedMap[migration.Filename] = true
	}

	// Run pending migrations
	pendingCount := 0
	for _, filename := range migrationFiles {
		if !appliedMap[filename] {
			if err := runSingleMigration(migrationsDir, filename); err != nil {
				return fmt.Errorf("failed to run migration %s: %v", filename, err)
			}
			pendingCount++
		}
	}

	if pendingCount == 0 {
		fmt.Println("No pending migrations")
	} else {
		fmt.Printf("Applied %d migrations\n", pendingCount)
	}

	return nil
}

func runSingleMigration(migrationsDir, filename string) error {
	filePath := filepath.Join(migrationsDir, filename)
	
	// Read migration file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %v", err)
	}

	// Execute SQL in transaction
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Execute the SQL
		if err := tx.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("failed to execute SQL: %v", err)
		}

		// Record the migration as applied
		migration := Migration{
			Filename:  filename,
			AppliedAt: time.Now(),
		}
		if err := tx.Create(&migration).Error; err != nil {
			return fmt.Errorf("failed to record migration: %v", err)
		}

		fmt.Printf("Applied migration: %s\n", filename)
		return nil
	})
}

func getMigrationFiles(migrationsDir string) ([]string, error) {
	var files []string
	
	err := filepath.WalkDir(migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".sql") {
			files = append(files, d.Name())
		}
		return nil
	})
	
	if err != nil {
		return nil, err
	}

	// Sort files chronologically
	sort.Strings(files)
	return files, nil
}

func createMigrationFile(migrationName string) error {
	// Create migrations directory if it doesn't exist
	migrationsDir := "./database/migrations"
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %v", err)
	}

	// Generate timestamp-based filename
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s.sql", timestamp, strings.ReplaceAll(migrationName, " ", "_"))
	filePath := filepath.Join(migrationsDir, filename)

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("migration file already exists: %s", filename)
	}

	// Create migration file with template
	template := fmt.Sprintf(`-- Migration: %s
-- Created: %s

-- Write your migration SQL here
-- Example:
-- CREATE TABLE example (
--   id INTEGER PRIMARY KEY AUTOINCREMENT,
--   name TEXT NOT NULL,
--   created_at DATETIME DEFAULT CURRENT_TIMESTAMP
-- );

-- For rollback, create a corresponding rollback file: %s_%s_rollback.sql
`, migrationName, time.Now().Format("2006-01-02 15:04:05"), timestamp, strings.ReplaceAll(migrationName, " ", "_"))

	if err := os.WriteFile(filePath, []byte(template), 0644); err != nil {
		return fmt.Errorf("failed to create migration file: %v", err)
	}

	fmt.Printf("Created migration file: %s\n", filename)
	return nil
}

func showMigrationStatus() error {
	// Ensure migrations table exists
	if err := ensureMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Get migration directory path
	migrationsDir := "./database/migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		fmt.Println("No migrations directory found")
		return nil
	}

	// Read migration files
	migrationFiles, err := getMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migration files: %v", err)
	}

	if len(migrationFiles) == 0 {
		fmt.Println("No migration files found")
		return nil
	}

	// Get applied migrations
	var appliedMigrations []Migration
	db.DB.Find(&appliedMigrations)

	appliedMap := make(map[string]Migration)
	for _, migration := range appliedMigrations {
		appliedMap[migration.Filename] = migration
	}

	// Display status
	fmt.Printf("%-50s %-12s %s\n", "Migration", "Status", "Applied At")
	fmt.Println(strings.Repeat("-", 80))

	for _, filename := range migrationFiles {
		if migration, exists := appliedMap[filename]; exists {
			fmt.Printf("%-50s %-12s %s\n", 
				filename, 
				"Applied", 
				migration.AppliedAt.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Printf("%-50s %-12s %s\n", filename, "Pending", "-")
		}
	}

	return nil
}

func init() {
	migrateCmd.AddCommand(runMigrationsCmd)
	migrateCmd.AddCommand(createMigrationCmd)
	migrateCmd.AddCommand(migrationStatusCmd)
}