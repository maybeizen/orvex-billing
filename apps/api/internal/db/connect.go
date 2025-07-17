package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite" // Default to SQLite for development
	}

	var dsn string
	var dialector gorm.Dialector

	switch dbType {
	case "mysql":
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "3306"
		}
		database := os.Getenv("DB_NAME")
		if database == "" {
			database = "billing"
		}
		username := os.Getenv("DB_USER")
		if username == "" {
			username = "root"
		}
		password := os.Getenv("DB_PASSWORD")
		if password == "" {
			password = ""
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			username, password, host, port, database)
		dialector = mysql.Open(dsn)

	case "sqlite":
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = "./database/database.db"
		}
		
		if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
			return fmt.Errorf("failed to create database directory: %v", err)
		}

		dialector = sqlite.Open(dbPath)

	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Printf("Successfully connected to %s database", dbType)
	return nil
}