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
	"github.com/orvexcc/billing/api/internal/utils"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

const (
	BackupCodeLength = 8
	BackupCodeCount  = 8
	TOTPWindow       = 1 // Allow 1 window before/after current time
)

// TwoFactorService handles all 2FA operations
type TwoFactorService struct {
	db *gorm.DB
}

// NewTwoFactorService creates a new 2FA service
func NewTwoFactorService(database *gorm.DB) *TwoFactorService {
	return &TwoFactorService{
		db: database,
	}
}

// GenerateTOTPKey generates a new TOTP key for a user
func (s *TwoFactorService) GenerateTOTPKey(email string) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      "Orvex",
		AccountName: email,
		SecretSize:  32,
	})
}

// ValidateTOTP validates a TOTP code against a secret with time window
func (s *TwoFactorService) ValidateTOTP(secret, token string) bool {
	if len(token) != 6 {
		return false
	}

	// Validate with time window for clock skew
	valid, _ := totp.ValidateCustom(token, secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      TOTPWindow,
		Digits:    6,
		Algorithm: otp.AlgorithmSHA1,
	})
	return valid
}

// GenerateBackupCodes generates cryptographically secure backup codes
func (s *TwoFactorService) GenerateBackupCodes() ([]string, error) {
	codes := make([]string, BackupCodeCount)
	for i := range codes {
		bytes := make([]byte, BackupCodeLength/2)
		if _, err := rand.Read(bytes); err != nil {
			return nil, fmt.Errorf("failed to generate random bytes: %w", err)
		}
		codes[i] = fmt.Sprintf("%x", bytes)
	}
	return codes, nil
}

// SaveBackupCodes saves backup codes to the database
func (s *TwoFactorService) SaveBackupCodes(userID uint, codes []string) error {
	// Delete existing backup codes
	if err := s.db.Where("user_id = ?", userID).Delete(&types.BackupCode{}).Error; err != nil {
		return fmt.Errorf("failed to delete existing backup codes: %w", err)
	}

	// Create new backup codes
	backupCodes := make([]types.BackupCode, len(codes))
	for i, code := range codes {
		backupCodes[i] = types.BackupCode{
			UserID: userID,
			Code:   code,
			Used:   false,
		}
	}

	if err := s.db.Create(&backupCodes).Error; err != nil {
		return fmt.Errorf("failed to create backup codes: %w", err)
	}

	return nil
}

// ValidateBackupCode validates and uses a backup code
func (s *TwoFactorService) ValidateBackupCode(userID uint, code string) error {
	// Normalize the code
	code = strings.ToLower(strings.TrimSpace(code))
	if len(code) != BackupCodeLength {
		return fmt.Errorf("invalid backup code length")
	}

	// Find and use the backup code
	var backupCode types.BackupCode
	result := s.db.Where("user_id = ? AND code = ? AND used = false", userID, code).First(&backupCode)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("invalid or already used backup code")
		}
		return fmt.Errorf("database error: %w", result.Error)
	}

	// Mark as used
	now := time.Now()
	backupCode.Used = true
	backupCode.UsedAt = &now

	if err := s.db.Save(&backupCode).Error; err != nil {
		return fmt.Errorf("failed to mark backup code as used: %w", err)
	}

	return nil
}

// GetUserWithValidation retrieves a user and validates they exist
func (s *TwoFactorService) GetUserWithValidation(userID interface{}) (*types.User, error) {
	if userID == nil {
		return nil, fmt.Errorf("user not authenticated")
	}

	var user types.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil
}

// SetSession2FAStatus sets the 2FA verification status in the session
func (s *TwoFactorService) SetSession2FAStatus(c *fiber.Ctx, verified bool) error {
	sess, err := middleware.SessionStore.Get(c)
	if err != nil {
		return fmt.Errorf("session error: %w", err)
	}

	sess.Set("2fa_verified", verified)
	if verified {
		sess.Set("2fa_verified_at", time.Now().Unix())
	} else {
		sess.Delete("2fa_verified_at")
	}

	if err := sess.Save(); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

// getTwoFactorService returns a 2FA service instance with current DB connection
func getTwoFactorService() *TwoFactorService {
	return NewTwoFactorService(db.DB)
}

// Setup2FA generates QR code and secret for 2FA setup
func Setup2FA(c *fiber.Ctx) error {
	service := getTwoFactorService()
	user, err := service.GetUserWithValidation(c.Locals("user_id"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	if user.TwoFactorEnabled {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Two-factor authentication is already enabled",
		})
	}

	// Generate TOTP key
	key, err := service.GenerateTOTPKey(user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to generate TOTP key",
		})
	}

	// Generate QR code
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

// Enable2FA enables 2FA for a user with password and TOTP verification
func Enable2FA(c *fiber.Ctx) error {
	service := getTwoFactorService()
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

	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Validation failed: " + err.Error(),
		})
	}

	user, err := service.GetUserWithValidation(c.Locals("user_id"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
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
	if !service.ValidateTOTP(req.Secret, req.Token) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid verification code",
		})
	}

	// Begin transaction
	tx := service.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Generate backup codes
	backupCodes, err := service.GenerateBackupCodes()
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to generate backup codes",
		})
	}

	// Save backup codes
	if err := service.SaveBackupCodes(user.ID, backupCodes); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to save backup codes",
		})
	}

	// Enable 2FA
	user.TwoFactorSecret = req.Secret
	user.TwoFactorEnabled = true

	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to enable two-factor authentication",
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to commit transaction",
		})
	}

	// Set 2FA verified in session since user just successfully enabled it
	if err := service.SetSession2FAStatus(c, true); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to set session status",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"backup_codes": backupCodes,
		},
	})
}

// Disable2FA disables 2FA for a user with TOTP verification
func Disable2FA(c *fiber.Ctx) error {
	service := getTwoFactorService()
	var req struct {
		Token string `json:"token" validate:"required,len=6"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request format",
		})
	}

	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Validation failed: " + err.Error(),
		})
	}

	user, err := service.GetUserWithValidation(c.Locals("user_id"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	if !user.TwoFactorEnabled {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Two-factor authentication is not enabled",
		})
	}

	// Verify TOTP code
	if !service.ValidateTOTP(user.TwoFactorSecret, req.Token) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid TOTP code",
		})
	}

	// Begin transaction
	tx := service.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete all backup codes
	if err := tx.Where("user_id = ?", user.ID).Delete(&types.BackupCode{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to delete backup codes",
		})
	}

	// Disable 2FA
	user.TwoFactorSecret = ""
	user.TwoFactorEnabled = false

	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to disable two-factor authentication",
		})
	}

	// Clear 2FA session
	if err := service.SetSession2FAStatus(c, false); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to clear session",
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to commit transaction",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Two-factor authentication disabled successfully",
	})
}

// Verify2FA verifies a TOTP code and sets session status
func Verify2FA(c *fiber.Ctx) error {
	service := getTwoFactorService()
	var req struct {
		Token string `json:"token" validate:"required,len=6"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request format",
		})
	}

	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Validation failed: " + err.Error(),
		})
	}

	user, err := service.GetUserWithValidation(c.Locals("user_id"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	if !user.TwoFactorEnabled {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Two-factor authentication is not enabled",
		})
	}

	// Verify TOTP code
	if !service.ValidateTOTP(user.TwoFactorSecret, req.Token) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid TOTP code",
		})
	}

	// Set 2FA verified in session
	if err := service.SetSession2FAStatus(c, true); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to save session",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Two-factor authentication verified successfully",
	})
}

// VerifyBackupCode verifies a backup code and sets session status
func VerifyBackupCode(c *fiber.Ctx) error {
	service := getTwoFactorService()
	var req struct {
		Code string `json:"code" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request format",
		})
	}

	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Validation failed: " + err.Error(),
		})
	}

	user, err := service.GetUserWithValidation(c.Locals("user_id"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	if !user.TwoFactorEnabled {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Two-factor authentication is not enabled",
		})
	}

	// Validate backup code
	if err := service.ValidateBackupCode(user.ID, req.Code); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Set 2FA verified in session
	if err := service.SetSession2FAStatus(c, true); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to save session",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Backup code verified successfully",
	})
}

// GenerateBackupCodes generates new backup codes for a user
func GenerateBackupCodes(c *fiber.Ctx) error {
	service := getTwoFactorService()
	
	var req struct {
		Token string `json:"token" validate:"required,len=6"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request format",
		})
	}
	
	user, err := service.GetUserWithValidation(c.Locals("user_id"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	if !user.TwoFactorEnabled {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Two-factor authentication is not enabled",
		})
	}

	// Verify TOTP code
	if !service.ValidateTOTP(user.TwoFactorSecret, req.Token) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid TOTP code",
		})
	}

	// Generate new backup codes
	backupCodes, err := service.GenerateBackupCodes()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to generate backup codes",
		})
	}

	// Save backup codes
	if err := service.SaveBackupCodes(user.ID, backupCodes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to save backup codes",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"backup_codes": backupCodes,
		},
	})
}