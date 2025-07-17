package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/middleware"
	"github.com/orvexcc/billing/api/internal/types"
	"github.com/orvexcc/billing/api/internal/utils"
)

func UpdateNotificationPreferences(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	var req types.NotificationPreferencesRequest
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

	user.EmailNotifications = req.EmailNotifications
	user.MarketingEmails = req.MarketingEmails
	user.SecurityNotifications = req.SecurityNotifications

	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update notification preferences",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notification preferences updated successfully",
	})
}

func UpdatePrivacySettings(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	var req types.PrivacySettingsRequest
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

	user.ProfilePublic = req.ProfilePublic
	user.ShowEmail = req.ShowEmail

	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update privacy settings",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Privacy settings updated successfully",
	})
}

func DeleteAccount(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	var req types.AccountDeletionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Note: Skipping validation for this example since the validator 
	// "eq=DELETE_MY_ACCOUNT" might not work as expected
	if req.Confirmation != "DELETE_MY_ACCOUNT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Confirmation text must be 'DELETE_MY_ACCOUNT'",
		})
	}

	var user types.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if !middleware.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Password is incorrect",
		})
	}

	if user.AvatarURL != "" {
		utils.DeleteAvatar(user.AvatarURL)
	}

	db.DB.Where("user_id = ?", userID).Delete(&types.Session{})

	if err := db.DB.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete account",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Account deleted successfully",
	})
}