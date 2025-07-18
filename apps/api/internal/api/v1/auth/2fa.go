package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/middleware"
	"github.com/orvexcc/billing/api/internal/types"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

func GenerateTOTPKey(email string) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      "Orvex",
		AccountName: email,
		SecretSize:  32,
	})
}

func VerifyTOTP(secret, token string) bool {
	return totp.Validate(token, secret)
}

func generateBackupCodes() []string {
	codes := make([]string, 8)
	for i := range codes {
		code := make([]byte, 5)
		rand.Read(code)
		// Generate 8-character alphanumeric codes
		codes[i] = fmt.Sprintf("%x", code)[:8]
	}
	return codes
}

func Setup2FA(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "User not authenticated",
		})
	}

	var user types.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}

	if user.TwoFactorEnabled {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Two-factor authentication is already enabled",
		})
	}

	// Generate new TOTP key
	key, err := GenerateTOTPKey(user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to generate TOTP key",
		})
	}

	// Generate QR code as base64 encoded PNG
	qrCode, err := qrcode.Encode(key.String(), qrcode.Medium, 256)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to generate QR code",
		})
	}

	qrCodeBase64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrCode)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"secret":  key.Secret(),
			"qr_code": qrCodeBase64,
		},
	})
}

func Enable2FA(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "User not authenticated",
		})
	}

	var req struct {
		Secret   string `json:"secret" validate:"required"`
		Token    string `json:"token" validate:"required,len=6"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request format",
		})
	}

	var user types.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}

	if user.TwoFactorEnabled {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Two-factor authentication is already enabled",
		})
	}

	// Verify password
	if !middleware.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid password",
		})
	}

	// Verify TOTP code
	if !VerifyTOTP(req.Secret, req.Token) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid verification code",
		})
	}

	// Generate backup codes
	backupCodes := generateBackupCodes()

	// Enable 2FA
	user.TwoFactorSecret = req.Secret
	user.TwoFactorEnabled = true

	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to enable two-factor authentication",
		})
	}

	// Save backup codes to database
	for _, code := range backupCodes {
		backupCode := types.BackupCode{
			UserID: user.ID,
			Code:   code,
			Used:   false,
		}
		if err := db.DB.Create(&backupCode).Error; err != nil {
			// If backup code creation fails, disable 2FA and return error
			user.TwoFactorSecret = ""
			user.TwoFactorEnabled = false
			db.DB.Save(&user)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to create backup codes",
			})
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"backup_codes": backupCodes,
		},
	})
}

func Disable2FA(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "User not authenticated",
		})
	}

	var req struct {
		Token string `json:"token" validate:"required,len=6"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request format",
		})
	}

	var user types.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}

	if !user.TwoFactorEnabled {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Two-factor authentication is not enabled",
		})
	}

	// Verify TOTP code before disabling
	if !VerifyTOTP(user.TwoFactorSecret, req.Token) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid TOTP code",
		})
	}

	// Delete all backup codes for this user
	db.DB.Where("user_id = ?", userID).Delete(&types.BackupCode{})

	// Disable 2FA
	user.TwoFactorSecret = ""
	user.TwoFactorEnabled = false

	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to disable two-factor authentication",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func Verify2FA(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "User not authenticated",
		})
	}

	var req struct {
		Token string `json:"token" validate:"required,len=6"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request format",
		})
	}

	var user types.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}

	if !user.TwoFactorEnabled {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Two-factor authentication is not enabled",
		})
	}

	// Verify TOTP code
	if !VerifyTOTP(user.TwoFactorSecret, req.Token) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid verification code",
		})
	}

	// Update session to mark 2FA as verified
	sess, err := middleware.SessionStore.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Session error",
		})
	}

	sess.Set("2fa_verified", true)
	sess.Set("2fa_verified_at", time.Now().Unix())
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to save session",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func GenerateBackupCodes(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "User not authenticated",
		})
	}

	var user types.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}

	if !user.TwoFactorEnabled {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Two-factor authentication is not enabled",
		})
	}

	// Delete all existing backup codes for this user
	db.DB.Where("user_id = ?", userID).Delete(&types.BackupCode{})

	// Generate new backup codes
	backupCodes := generateBackupCodes()

	// Save new backup codes to database
	for _, code := range backupCodes {
		backupCode := types.BackupCode{
			UserID: user.ID,
			Code:   code,
			Used:   false,
		}
		if err := db.DB.Create(&backupCode).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to create backup codes",
			})
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"backup_codes": backupCodes,
		},
	})
}

func VerifyBackupCode(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "User not authenticated",
		})
	}

	var req struct {
		Code string `json:"code" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request format",
		})
	}

	// Validate and normalize the backup code
	req.Code = strings.ToLower(strings.TrimSpace(req.Code))
	if len(req.Code) != 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Backup code must be 8 characters long",
		})
	}

	var user types.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}

	if !user.TwoFactorEnabled {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Two-factor authentication is not enabled",
		})
	}

	// Check if backup code exists and is unused
	var backupCode types.BackupCode
	result = db.DB.Where("user_id = ? AND code = ? AND used = false", userID, req.Code).First(&backupCode)
	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid or already used backup code",
		})
	}

	// Mark backup code as used
	backupCode.Used = true
	now := time.Now()
	backupCode.UsedAt = &now
	if err := db.DB.Save(&backupCode).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to mark backup code as used",
		})
	}

	// Set 2FA verified in session
	sess, err := middleware.SessionStore.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Session error",
		})
	}

	sess.Set("2fa_verified", true)
	sess.Set("2fa_verified_at", time.Now().Unix())
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to save session",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}