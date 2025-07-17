package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
)

func GetUserByUUID(c *fiber.Ctx) error {
	uuidParam := c.Params("uuid")
	if uuidParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "UUID is required",
		})
	}

	// Parse UUID
	userUUID, err := uuid.Parse(uuidParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID format",
		})
	}

	var user types.User
	result := db.DB.Where("uuid = ?", userUUID).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Check if profile is public or if user is authenticated and accessing their own profile
	requestingUserID := c.Locals("user_id")
	if !user.ProfilePublic && (requestingUserID == nil || requestingUserID.(uint) != user.ID) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Profile is private",
		})
	}

	// Remove sensitive data
	user.PasswordHash = ""
	user.TwoFactorSecret = ""
	user.FailedAttempts = 0
	user.AccountLocked = false

	// Conditionally hide email
	if !user.ShowEmail && (requestingUserID == nil || requestingUserID.(uint) != user.ID) {
		user.Email = ""
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user":    user,
	})
}