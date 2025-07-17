package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/middleware"
	"github.com/orvexcc/billing/api/internal/types"
)

func Logout(c *fiber.Ctx) error {
	sess, err := middleware.SessionStore.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Session error",
		})
	}

	userID := sess.Get("user_id")
	if userID != nil {
		db.DB.Where("user_id = ?", userID).Delete(&types.Session{})
	}

	if err := sess.Destroy(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Session destruction failed",
		})
	}

	return c.JSON(types.AuthResponse{
		Success: true,
		Message: "Logout successful",
	})
}

func LogoutAll(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	db.DB.Where("user_id = ?", userID).Delete(&types.Session{})

	sess, err := middleware.SessionStore.Get(c)
	if err == nil {
		sess.Destroy()
	}

	return c.JSON(types.AuthResponse{
		Success: true,
		Message: "All sessions terminated",
	})
}