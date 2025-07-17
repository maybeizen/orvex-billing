package auth

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
)

func GenerateTOTPSecret() string {
	secret := make([]byte, 20)
	rand.Read(secret)
	return base32.StdEncoding.EncodeToString(secret)
}

func Enable2FA(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	var user types.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if user.TwoFactorEnabled {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Two-factor authentication is already enabled",
		})
	}

	secret := GenerateTOTPSecret()
	
	qrURL := fmt.Sprintf("otpauth://totp/YourApp:%s?secret=%s&issuer=YourApp", 
		user.Email, secret)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Scan this QR code with your authenticator app",
		"qr_url":  qrURL,
		"secret":  secret,
		"note":    "Save this secret safely. You'll need to verify with a TOTP code to complete setup.",
	})
}

func Confirm2FA(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	var req struct {
		Secret   string `json:"secret" validate:"required"`
		TOTPCode string `json:"totp_code" validate:"required,len=6"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	var user types.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// In production, verify the TOTP code here
	// For now, we'll just accept any 6-digit code
	if len(req.TOTPCode) != 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid TOTP code",
		})
	}

	// Enable 2FA
	user.TwoFactorSecret = req.Secret
	user.TwoFactorEnabled = true

	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to enable two-factor authentication",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Two-factor authentication enabled successfully",
	})
}

func Disable2FA(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	var req struct {
		TOTPCode string `json:"totp_code" validate:"required,len=6"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	var user types.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if !user.TwoFactorEnabled {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Two-factor authentication is not enabled",
		})
	}

	// In production, verify the TOTP code here
	if len(req.TOTPCode) != 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid TOTP code",
		})
	}

	// Disable 2FA
	user.TwoFactorSecret = ""
	user.TwoFactorEnabled = false

	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to disable two-factor authentication",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Two-factor authentication disabled successfully",
	})
}