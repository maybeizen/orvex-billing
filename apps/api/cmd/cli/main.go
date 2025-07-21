package main

import (
	"log"
	"os"

	"github.com/orvexcc/billing/api/internal/cli"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/utils"
)

func main() {
	// Load environment variables
	utils.LoadEnv()

	// Connect to database
	log.Println("Connecting to database...")
	err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Execute CLI commands
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}