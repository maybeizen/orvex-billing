package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
	"golang.org/x/crypto/bcrypt"
)

const (
	maxFailedAttempts = 5
	lockoutDuration   = 15 * time.Minute
	bcryptCost        = 12
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func CreateSession(userID uint, ipAddress, userAgent string) (*types.Session, error) {
	sessionID, err := GenerateSecureToken()
	if err != nil {
		return nil, err
	}

	session := &types.Session{
		ID:        sessionID,
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	result := db.DB.Create(session)
	if result.Error != nil {
		return nil, result.Error
	}

	return session, nil
}

func CleanupExpiredSessions() {
	db.DB.Where("expires_at < ?", time.Now()).Delete(&types.Session{})
}

func RateLimitByIP() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Basic IP-based rate limiting can be implemented here
		// For production, consider using Redis or external rate limiting service
		return c.Next()
	}
}

func RequireAdmin(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	var user types.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Check if user has admin role
	if user.Role != types.UserRoleAdmin && user.Role != types.UserRoleSuper {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	// Store user object in context for admin operations
	c.Locals("admin_user", user)
	return c.Next()
}

func RequireSuperAdmin(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	var user types.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Check if user has super admin role
	if user.Role != types.UserRoleSuper {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Super admin access required",
		})
	}

	c.Locals("admin_user", user)
	return c.Next()
}