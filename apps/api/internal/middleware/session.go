package middleware

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
)

var SessionStore *session.Store

func InitSessions() {
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "default-secret-change-in-production"
	}

	storage := sqlite3.New(sqlite3.Config{
		Database:   "./database/sessions.db",
		Table:      "sessions",
		Reset:      false,
		GCInterval: 10 * time.Second,
	})

	SessionStore = session.New(session.Config{
		Storage:        storage,
		Expiration:     24 * time.Hour,
		CookieSecure:   os.Getenv("APP_ENV") == "production",
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
		KeyLookup:      "cookie:session_id",
	})
}

func RequireAuth(c *fiber.Ctx) error {
	sess, err := SessionStore.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Session error",
		})
	}

	userID := sess.Get("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	// Validate that the user actually exists in the database
	var user types.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		// User doesn't exist, destroy session
		sess.Destroy()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	c.Locals("user_id", userID)
	c.Locals("user", user) // Also provide the user object
	return c.Next()
}

func OptionalAuth(c *fiber.Ctx) error {
	sess, err := SessionStore.Get(c)
	if err == nil {
		userID := sess.Get("user_id")
		if userID != nil {
			// Validate that the user actually exists in the database
			var user types.User
			if err := db.DB.First(&user, userID).Error; err == nil {
				c.Locals("user_id", userID)
				c.Locals("user", user)
			} else {
				// User doesn't exist, destroy session
				sess.Destroy()
			}
		}
	}
	return c.Next()
}