package user

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/middleware"
	"github.com/orvexcc/billing/api/internal/types"
)

type SessionResponse struct {
	ID           string    `json:"id"`
	UserID       uint      `json:"user_id"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	CreatedAt    time.Time `json:"created_at"`
	LastActivity time.Time `json:"last_activity"`
	IsCurrent    bool      `json:"is_current"`
}

func GetUserSessions(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Get current session ID
	sess, err := middleware.SessionStore.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Session error",
		})
	}
	currentSessionID := sess.ID()

	var sessions []types.Session
	result := db.DB.Where("user_id = ?", userID).Find(&sessions)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch sessions",
		})
	}

	// Convert to response format
	var sessionResponses []SessionResponse
	for _, session := range sessions {
		sessionResponses = append(sessionResponses, SessionResponse{
			ID:           session.ID,
			UserID:       session.UserID,
			IPAddress:    session.IPAddress,
			UserAgent:    session.UserAgent,
			CreatedAt:    session.CreatedAt,
			LastActivity: session.UpdatedAt,
			IsCurrent:    session.ID == currentSessionID,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    sessionResponses,
	})
}

func RevokeSession(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	sessionID := c.Params("sessionId")
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Session ID is required",
		})
	}

	// Get current session ID to prevent self-revocation
	sess, err := middleware.SessionStore.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Session error",
		})
	}
	currentSessionID := sess.ID()

	if sessionID == currentSessionID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot revoke current session",
		})
	}

	// Delete the session
	result := db.DB.Where("id = ? AND user_id = ?", sessionID, userID).Delete(&types.Session{})
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to revoke session",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Session not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Session revoked successfully",
	})
}

func RevokeAllOtherSessions(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Get current session ID
	sess, err := middleware.SessionStore.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Session error",
		})
	}
	currentSessionID := sess.ID()

	// Delete all sessions except current
	result := db.DB.Where("user_id = ? AND id != ?", userID, currentSessionID).Delete(&types.Session{})
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to revoke sessions",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "All other sessions revoked successfully",
		"revoked_count": result.RowsAffected,
	})
}