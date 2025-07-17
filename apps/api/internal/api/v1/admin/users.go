package admin

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/middleware"
	"github.com/orvexcc/billing/api/internal/types"
	"github.com/orvexcc/billing/api/internal/utils"
)

type AdminUserRequest struct {
	Email                  string         `json:"email" validate:"required,email"`
	Username               string         `json:"username" validate:"required,min=3,max=50,alphanum"`
	FirstName              string         `json:"first_name" validate:"required,min=1,max=50"`
	LastName               string         `json:"last_name" validate:"required,min=1,max=50"`
	Password               string         `json:"password,omitempty" validate:"omitempty,min=8"`
	Role                   types.UserRole `json:"role" validate:"required,oneof=user admin super_admin"`
	EmailVerified          bool           `json:"email_verified"`
	AccountLocked          bool           `json:"account_locked"`
	EmailNotifications     bool           `json:"email_notifications"`
	MarketingEmails        bool           `json:"marketing_emails"`
	SecurityNotifications bool           `json:"security_notifications"`
	ProfilePublic          bool           `json:"profile_public"`
	ShowEmail              bool           `json:"show_email"`
	Bio                    string         `json:"bio,omitempty" validate:"max=500"`
}

type AdminUserUpdateRequest struct {
	Email                  string         `json:"email,omitempty" validate:"omitempty,email"`
	Username               string         `json:"username,omitempty" validate:"omitempty,min=3,max=50,alphanum"`
	FirstName              string         `json:"first_name,omitempty" validate:"omitempty,min=1,max=50"`
	LastName               string         `json:"last_name,omitempty" validate:"omitempty,min=1,max=50"`
	Password               string         `json:"password,omitempty" validate:"omitempty,min=8"`
	Role                   types.UserRole `json:"role,omitempty" validate:"omitempty,oneof=user admin super_admin"`
	EmailVerified          *bool          `json:"email_verified,omitempty"`
	AccountLocked          *bool          `json:"account_locked,omitempty"`
	EmailNotifications     *bool          `json:"email_notifications,omitempty"`
	MarketingEmails        *bool          `json:"marketing_emails,omitempty"`
	SecurityNotifications *bool          `json:"security_notifications,omitempty"`
	ProfilePublic          *bool          `json:"profile_public,omitempty"`
	ShowEmail              *bool          `json:"show_email,omitempty"`
	Bio                    string         `json:"bio,omitempty" validate:"max=500"`
}

type BulkUserAction struct {
	Action  string `json:"action" validate:"required,oneof=lock unlock verify_email delete change_role"`
	UserIDs []uint `json:"user_ids" validate:"required,min=1,max=100"`
	Role    types.UserRole `json:"role,omitempty" validate:"omitempty,oneof=user admin super_admin"`
}

func ListUsers(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)

	// Security: Log admin access
	logAdminActivity("users_list", adminUser.ID, adminUser.Email, fiber.Map{
		"ip": c.IP(),
	})

	// Query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "25"))
	search := c.Query("search")
	role := c.Query("role")
	status := c.Query("status") // active, locked, unverified
	sortBy := c.Query("sort", "created_at")
	sortOrder := c.Query("order", "desc")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 25
	}

	// Security: Validate sort parameters
	validSorts := map[string]bool{
		"created_at": true, "email": true, "username": true, "last_login_at": true,
	}
	if !validSorts[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" {
		sortOrder = "desc"
	}

	query := db.DB.Model(&types.User{})

	// Apply filters
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("email ILIKE ? OR username ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm)
	}

	if role != "" {
		query = query.Where("role = ?", role)
	}

	switch status {
	case "locked":
		query = query.Where("account_locked = ?", true)
	case "active":
		query = query.Where("account_locked = ?", false)
	case "unverified":
		query = query.Where("email_verified = ?", false)
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get users with pagination
	var users []types.User
	offset := (page - 1) * pageSize
	orderClause := sortBy + " " + sortOrder
	result := query.Order(orderClause).
		Offset(offset).
		Limit(pageSize).
		Find(&users)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}

	// Remove sensitive data
	for i := range users {
		users[i].PasswordHash = ""
		users[i].TwoFactorSecret = ""
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    users,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"sort_by":   sortBy,
			"sort_order": sortOrder,
		},
	})
}

func GetUser(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	userID := c.Params("id")

	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	var user types.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Security: Log admin access
	logAdminActivity("user_view", adminUser.ID, adminUser.Email, fiber.Map{
		"target_user_id": user.ID,
		"target_email":   user.Email,
	})

	// Remove sensitive data
	user.PasswordHash = ""
	user.TwoFactorSecret = ""

	// Get related data
	var services []types.Service
	db.DB.Where("user_id = ?", user.ID).Find(&services)

	var invoices []types.Invoice
	db.DB.Where("user_id = ?", user.ID).Order("created_at DESC").Limit(10).Find(&invoices)

	var sessions []types.Session
	db.DB.Where("user_id = ?", user.ID).Order("created_at DESC").Limit(5).Find(&sessions)

	var carts []types.Cart
	db.DB.Where("user_id = ?", user.ID).Order("created_at DESC").Limit(3).Find(&carts)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user":     user,
			"services": services,
			"invoices": invoices,
			"sessions": sessions,
			"carts":    carts,
		},
	})
}

func CreateUser(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)

	var req AdminUserRequest
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

	// Security: Only super admins can create admin users
	if req.Role == types.UserRoleAdmin || req.Role == types.UserRoleSuper {
		if adminUser.Role != types.UserRoleSuper {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Only super admins can create admin users",
			})
		}
	}

	// Check if email or username already exists
	var existingUser types.User
	if result := db.DB.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser); result.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User with this email or username already exists",
		})
	}

	// Hash password
	hashedPassword, err := middleware.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Create user
	user := types.User{
		UUID:                   uuid.New(),
		Email:                  req.Email,
		Username:               req.Username,
		PasswordHash:           hashedPassword,
		FirstName:              req.FirstName,
		LastName:               req.LastName,
		Bio:                    req.Bio,
		Role:                   req.Role,
		EmailVerified:          req.EmailVerified,
		AccountLocked:          req.AccountLocked,
		EmailNotifications:     req.EmailNotifications,
		MarketingEmails:        req.MarketingEmails,
		SecurityNotifications: req.SecurityNotifications,
		ProfilePublic:          req.ProfilePublic,
		ShowEmail:              req.ShowEmail,
	}

	result := db.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Security: Log admin action
	logAdminActivity("user_create", adminUser.ID, adminUser.Email, fiber.Map{
		"created_user_id": user.ID,
		"created_email":   user.Email,
		"role":            user.Role,
	})

	// Remove sensitive data
	user.PasswordHash = ""
	user.TwoFactorSecret = ""

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "User created successfully",
		"data":    user,
	})
}

func UpdateUser(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	userID := c.Params("id")

	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	var req AdminUserUpdateRequest
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

	// Security: Role change restrictions
	if req.Role != "" && req.Role != user.Role {
		// Only super admins can change roles
		if adminUser.Role != types.UserRoleSuper {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Only super admins can change user roles",
			})
		}
		
		// Cannot change super admin role
		if user.Role == types.UserRoleSuper && adminUser.ID != user.ID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Cannot change super admin role",
			})
		}
	}

	// Security: Self-modification restrictions
	if adminUser.ID == user.ID {
		if req.AccountLocked != nil && *req.AccountLocked {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Cannot lock your own account",
			})
		}
		if req.Role != "" && req.Role != user.Role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Cannot change your own role",
			})
		}
	}

	// Update fields if provided
	changes := make(map[string]interface{})
	
	if req.Email != "" && req.Email != user.Email {
		// Check if email already exists
		var existingUser types.User
		if result := db.DB.Where("email = ? AND id != ?", req.Email, user.ID).First(&existingUser); result.Error == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email already exists",
			})
		}
		user.Email = req.Email
		user.EmailVerified = false // Reset verification on email change
		changes["email"] = req.Email
	}

	if req.Username != "" && req.Username != user.Username {
		// Check if username already exists
		var existingUser types.User
		if result := db.DB.Where("username = ? AND id != ?", req.Username, user.ID).First(&existingUser); result.Error == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Username already exists",
			})
		}
		user.Username = req.Username
		changes["username"] = req.Username
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
		changes["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
		changes["last_name"] = req.LastName
	}
	if req.Bio != "" {
		user.Bio = req.Bio
		changes["bio"] = req.Bio
	}
	if req.Role != "" {
		user.Role = req.Role
		changes["role"] = req.Role
	}
	if req.EmailVerified != nil {
		user.EmailVerified = *req.EmailVerified
		changes["email_verified"] = *req.EmailVerified
	}
	if req.AccountLocked != nil {
		user.AccountLocked = *req.AccountLocked
		changes["account_locked"] = *req.AccountLocked
	}
	if req.EmailNotifications != nil {
		user.EmailNotifications = *req.EmailNotifications
	}
	if req.MarketingEmails != nil {
		user.MarketingEmails = *req.MarketingEmails
	}
	if req.SecurityNotifications != nil {
		user.SecurityNotifications = *req.SecurityNotifications
	}
	if req.ProfilePublic != nil {
		user.ProfilePublic = *req.ProfilePublic
	}
	if req.ShowEmail != nil {
		user.ShowEmail = *req.ShowEmail
	}

	// Update password if provided
	if req.Password != "" {
		hashedPassword, err := middleware.HashPassword(req.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to hash password",
			})
		}
		user.PasswordHash = hashedPassword
		changes["password"] = "updated"
		
		// Invalidate all user sessions on password change
		db.DB.Where("user_id = ?", user.ID).Delete(&types.Session{})
	}

	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	// Security: Log admin action
	logAdminActivity("user_update", adminUser.ID, adminUser.Email, fiber.Map{
		"target_user_id": user.ID,
		"target_email":   user.Email,
		"changes":        changes,
	})

	// Remove sensitive data
	user.PasswordHash = ""
	user.TwoFactorSecret = ""

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User updated successfully",
		"data":    user,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	userID := c.Params("id")

	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	var user types.User
	result := db.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Security: Cannot delete super admin or yourself
	if user.Role == types.UserRoleSuper {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Cannot delete super admin user",
		})
	}

	if adminUser.ID == user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Cannot delete your own account",
		})
	}

	// Soft delete user
	if err := db.DB.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	// Delete user sessions
	db.DB.Where("user_id = ?", user.ID).Delete(&types.Session{})

	// Delete user avatar if exists
	if user.AvatarURL != "" {
		utils.DeleteAvatar(user.AvatarURL)
	}

	// Security: Log admin action
	logAdminActivity("user_delete", adminUser.ID, adminUser.Email, fiber.Map{
		"deleted_user_id": user.ID,
		"deleted_email":   user.Email,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User deleted successfully",
	})
}

func BulkUserActions(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)

	var req BulkUserAction
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

	// Security: Limit bulk actions to prevent abuse
	if len(req.UserIDs) > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Maximum 100 users can be processed at once",
		})
	}

	// Get all users for validation
	var users []types.User
	result := db.DB.Where("id IN ?", req.UserIDs).Find(&users)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	if len(users) != len(req.UserIDs) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Some users not found",
		})
	}

	// Security validations
	for _, user := range users {
		// Cannot perform actions on super admin (unless you are super admin)
		if user.Role == types.UserRoleSuper && adminUser.Role != types.UserRoleSuper {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Cannot perform actions on super admin users",
			})
		}

		// Cannot perform actions on yourself for certain operations
		if user.ID == adminUser.ID {
			switch req.Action {
			case "lock", "delete":
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Cannot perform this action on your own account",
				})
			case "change_role":
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Cannot change your own role",
				})
			}
		}
	}

	// Security: Role changes require super admin
	if req.Action == "change_role" && adminUser.Role != types.UserRoleSuper {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only super admins can change user roles",
		})
	}

	// Perform bulk action
	var successCount, failureCount int
	var errors []string

	for _, user := range users {
		switch req.Action {
		case "lock":
			user.AccountLocked = true
			if err := db.DB.Save(&user).Error; err != nil {
				failureCount++
				errors = append(errors, fmt.Sprintf("Failed to lock user %s", user.Email))
			} else {
				successCount++
				// Invalidate user sessions
				db.DB.Where("user_id = ?", user.ID).Delete(&types.Session{})
			}

		case "unlock":
			user.AccountLocked = false
			user.FailedAttempts = 0
			if err := db.DB.Save(&user).Error; err != nil {
				failureCount++
				errors = append(errors, fmt.Sprintf("Failed to unlock user %s", user.Email))
			} else {
				successCount++
			}

		case "verify_email":
			user.EmailVerified = true
			if err := db.DB.Save(&user).Error; err != nil {
				failureCount++
				errors = append(errors, fmt.Sprintf("Failed to verify email for user %s", user.Email))
			} else {
				successCount++
			}

		case "delete":
			// Security: Cannot delete super admin
			if user.Role == types.UserRoleSuper {
				failureCount++
				errors = append(errors, fmt.Sprintf("Cannot delete super admin user %s", user.Email))
				continue
			}

			if err := db.DB.Delete(&user).Error; err != nil {
				failureCount++
				errors = append(errors, fmt.Sprintf("Failed to delete user %s", user.Email))
			} else {
				successCount++
				// Clean up user data
				db.DB.Where("user_id = ?", user.ID).Delete(&types.Session{})
				if user.AvatarURL != "" {
					utils.DeleteAvatar(user.AvatarURL)
				}
			}

		case "change_role":
			if req.Role == "" {
				failureCount++
				errors = append(errors, fmt.Sprintf("Role is required for user %s", user.Email))
				continue
			}

			user.Role = req.Role
			if err := db.DB.Save(&user).Error; err != nil {
				failureCount++
				errors = append(errors, fmt.Sprintf("Failed to change role for user %s", user.Email))
			} else {
				successCount++
			}

		default:
			failureCount++
			errors = append(errors, fmt.Sprintf("Invalid action for user %s", user.Email))
		}
	}

	// Security: Log bulk action
	logAdminActivity("bulk_user_action", adminUser.ID, adminUser.Email, fiber.Map{
		"action":        req.Action,
		"user_ids":      req.UserIDs,
		"success_count": successCount,
		"failure_count": failureCount,
		"role":          req.Role,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Bulk action completed: %d successful, %d failed", successCount, failureCount),
		"data": fiber.Map{
			"success_count": successCount,
			"failure_count": failureCount,
			"errors":        errors,
		},
	})
}

func ImpersonateUser(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	userID := c.Params("id")

	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	var targetUser types.User
	result := db.DB.First(&targetUser, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Security: Cannot impersonate super admin or yourself
	if targetUser.Role == types.UserRoleSuper {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Cannot impersonate super admin user",
		})
	}

	if adminUser.ID == targetUser.ID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot impersonate yourself",
		})
	}

	// Security: Account must not be locked
	if targetUser.AccountLocked {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Cannot impersonate locked account",
		})
	}

	// Create impersonation session
	sessionID, err := middleware.GenerateSecureToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate session",
		})
	}

	// Create special impersonation session with metadata
	impersonationSession := &types.Session{
		ID:        sessionID,
		UserID:    targetUser.ID,
		IPAddress: c.IP(),
		UserAgent: c.Get("User-Agent") + " [ADMIN_IMPERSONATION]",
		ExpiresAt: time.Now().Add(1 * time.Hour), // Short session for impersonation
	}

	if err := db.DB.Create(impersonationSession).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create impersonation session",
		})
	}

	// Security: Log impersonation
	logAdminActivity("user_impersonation", adminUser.ID, adminUser.Email, fiber.Map{
		"target_user_id": targetUser.ID,
		"target_email":   targetUser.Email,
		"session_id":     sessionID,
		"ip":             c.IP(),
	})

	// Set session cookie
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  impersonationSession.ExpiresAt,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User impersonation started",
		"data": fiber.Map{
			"user":       targetUser,
			"session_id": sessionID,
			"expires_at": impersonationSession.ExpiresAt,
		},
	})
}

func StopImpersonation(c *fiber.Ctx) error {
	sessionID := c.Cookies("session_id")
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No active session found",
		})
	}

	// Find and validate session
	var session types.Session
	result := db.DB.Where("id = ?", sessionID).First(&session)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Session not found",
		})
	}

	// Check if this is an impersonation session
	if !strings.Contains(session.UserAgent, "[ADMIN_IMPERSONATION]") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Not an impersonation session",
		})
	}

	// Delete impersonation session
	if err := db.DB.Delete(&session).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to end impersonation session",
		})
	}

	// Clear session cookie
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	// Security: Log impersonation end
	logAdminActivity("impersonation_ended", session.UserID, "", fiber.Map{
		"session_id": sessionID,
		"ip":         c.IP(),
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Impersonation ended successfully",
	})
}