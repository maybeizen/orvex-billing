package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/middleware"
	"github.com/orvexcc/billing/api/internal/types"
	"github.com/orvexcc/billing/api/internal/utils"
)

func Register(c *fiber.Ctx) error {
	var req types.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.AuthResponse{
			Success: false,
			Message: "Invalid request format",
		})
	}

	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.AuthResponse{
			Success: false,
			Message: "Validation failed: " + err.Error(),
		})
	}

	var existingUser types.User
	if result := db.DB.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser); result.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(types.AuthResponse{
			Success: false,
			Message: "User with this email or username already exists",
		})
	}

	hashedPassword, err := middleware.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.AuthResponse{
			Success: false,
			Message: "Password hashing failed",
		})
	}

	user := types.User{
		UUID:         uuid.New(),
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	}

	result := db.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.AuthResponse{
			Success: false,
			Message: "User creation failed",
		})
	}

	user.PasswordHash = ""
	return c.Status(fiber.StatusCreated).JSON(types.AuthResponse{
		Success: true,
		Message: "Registration successful",
		User:    &user,
	})
}