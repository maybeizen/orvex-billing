package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/middleware"
	"github.com/orvexcc/billing/api/internal/types"
	"github.com/orvexcc/billing/api/internal/utils"
)

func Login(c *fiber.Ctx) error {
	var req types.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.AuthResponse{
			Success: false,
			Message: "Invalid request format",
		})
	}

	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.AuthResponse{
			Success: false,
			Message: "Validation failed: " + err.Error(),
		})
	}

	var user types.User
	result := db.DB.Where("email = ? OR username = ?", req.Email, req.Email).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(types.AuthResponse{
			Success: false,
			Message: "Invalid credentials",
		})
	}

	if user.AccountLocked {
		return c.Status(fiber.StatusUnauthorized).JSON(types.AuthResponse{
			Success: false,
			Message: "Account is locked",
		})
	}

	if !middleware.CheckPasswordHash(req.Password, user.PasswordHash) {
		user.FailedAttempts++
		if user.FailedAttempts >= 5 {
			user.AccountLocked = true
		}
		db.DB.Save(&user)
		return c.Status(fiber.StatusUnauthorized).JSON(types.AuthResponse{
			Success: false,
			Message: "Invalid credentials",
		})
	}

	// Don't require 2FA during initial login
	// This will be handled by middleware after login

	user.FailedAttempts = 0
	user.AccountLocked = false
	now := time.Now()
	user.LastLoginAt = &now
	db.DB.Save(&user)

	sess, err := middleware.SessionStore.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.AuthResponse{
			Success: false,
			Message: "Session error",
		})
	}

	sess.Set("user_id", user.ID)
	sess.Set("authenticated", true)
	// Don't set 2fa_verified during login - this will be set by the verify endpoint
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.AuthResponse{
			Success: false,
			Message: "Session save error",
		})
	}

	_, err = middleware.CreateSession(user.ID, c.IP(), c.Get("User-Agent"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.AuthResponse{
			Success: false,
			Message: "Session creation error",
		})
	}

	user.PasswordHash = ""
	return c.JSON(types.AuthResponse{
		Success: true,
		Message: "Login successful",
		User:    &user,
	})
}
