package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
	"github.com/orvexcc/billing/api/internal/utils"
)

func UploadAvatar(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	var user types.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if user.AvatarURL != "" {
		utils.DeleteAvatar(user.AvatarURL)
	}

	avatarURL, err := utils.SaveAvatar(file, user.ID)
	if err != nil {
		if fileErr, ok := err.(utils.FileValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fileErr.Message,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save avatar",
		})
	}

	user.AvatarURL = avatarURL
	if err := db.DB.Save(&user).Error; err != nil {
		utils.DeleteAvatar(avatarURL)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user avatar",
		})
	}

	return c.JSON(fiber.Map{
		"success":    true,
		"message":    "Avatar uploaded successfully",
		"avatar_url": avatarURL,
	})
}

func DeleteAvatar(c *fiber.Ctx) error {
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

	if user.AvatarURL != "" {
		utils.DeleteAvatar(user.AvatarURL)
		
		user.AvatarURL = ""
		if err := db.DB.Save(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update user avatar",
			})
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Avatar deleted successfully",
	})
}