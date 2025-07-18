package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/middleware"
	"github.com/orvexcc/billing/api/internal/types"
	"github.com/orvexcc/billing/api/internal/utils"
)

func GetProfile(c *fiber.Ctx) error {
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

	user.PasswordHash = ""
	user.TwoFactorSecret = ""

	// Check if 2FA is verified for this session
	twoFactorVerified := false
	if user.TwoFactorEnabled {
		sess, err := middleware.SessionStore.Get(c)
		if err == nil {
			if verified := sess.Get("2fa_verified"); verified != nil {
				if verifiedBool, ok := verified.(bool); ok {
					twoFactorVerified = verifiedBool
				}
			}
		}
	} else {
		// If 2FA is not enabled, consider it verified
		twoFactorVerified = true
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user":    user,
		"two_factor_verified": twoFactorVerified,
	})
}

func UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	var req types.ProfileUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed: " + err.Error(),
		})
	}

	var user types.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Bio = req.Bio

	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	user.PasswordHash = ""
	user.TwoFactorSecret = ""

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile updated successfully",
		"user":    user,
	})
}

func UpdateUsername(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	var req types.UsernameUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed: " + err.Error(),
		})
	}

	var existingUser types.User
	if result := db.DB.Where("username = ? AND id != ?", req.Username, userID).First(&existingUser); result.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Username is already taken",
		})
	}

	var user types.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	user.Username = req.Username
	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update username",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Username updated successfully",
	})
}

func UpdateEmail(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	var req types.EmailUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed: " + err.Error(),
		})
	}

	var existingUser types.User
	if result := db.DB.Where("email = ? AND id != ?", req.Email, userID).First(&existingUser); result.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email is already taken",
		})
	}

	var user types.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	user.Email = req.Email
	user.EmailVerified = false
	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update email",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Email updated successfully. Please verify your new email address.",
	})
}
