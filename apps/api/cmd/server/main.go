package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/api"
	"github.com/orvexcc/billing/api/internal/cli"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/middleware"
	"github.com/orvexcc/billing/api/internal/types"
	"github.com/orvexcc/billing/api/internal/utils"
)

func main() {
	utils.LoadEnv()

	// Check if CLI arguments are provided
	if len(os.Args) > 1 {
		// Connect to database for CLI operations
		log.Println("Connecting to database...")
		err := db.Connect()
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		// Execute CLI commands
		if err := cli.Execute(); err != nil {
			os.Exit(1)
		}
		return
	}

	// Server mode
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Connecting to database...")
	err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Running database migrations...")
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
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Initializing sessions...")
	middleware.InitSessions()

	app := fiber.New(fiber.Config{
		ProxyHeader: "X-Forwarded-For",
	})

	api.RegisterRoutes(app)

	log.Printf("Server is running on port %s", port)
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

