package services

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
	"gorm.io/gorm"
)

func ListServices(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Query parameters for pagination and filtering
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := db.DB.Where("user_id = ?", userID)
	
	// Filter by status if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	var total int64
	query.Model(&types.Service{}).Count(&total)

	// Get services with pagination
	var services []types.Service
	offset := (page - 1) * pageSize
	result := query.Preload("Invoices").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&services)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve services",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    services,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetService(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	serviceID := c.Params("id")
	if serviceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Service ID is required",
		})
	}

	var service types.Service
	result := db.DB.Where("id = ? AND user_id = ?", serviceID, userID).
		Preload("Invoices", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("Invoices.Items").
		First(&service)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Service not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    service,
	})
}

func GetServiceStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	serviceID := c.Params("id")
	if serviceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Service ID is required",
		})
	}

	var service types.Service
	result := db.DB.Where("id = ? AND user_id = ?", serviceID, userID).
		First(&service)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Service not found",
		})
	}

	// Calculate status information
	statusInfo := fiber.Map{
		"service_id":    service.ID,
		"status":        service.Status,
		"activated_at":  service.ActivatedAt,
		"suspended_at":  service.SuspendedAt,
		"cancelled_at":  service.CancelledAt,
		"next_bill_date": service.NextBillDate,
		"last_bill_date": service.LastBillDate,
		"billing_type":  service.BillingType,
	}

	// Add status-specific information
	switch service.Status {
	case types.ServiceStatusActive:
		statusInfo["message"] = "Service is active and running"
		statusInfo["health"] = "healthy"
	case types.ServiceStatusPending:
		statusInfo["message"] = "Service is pending activation"
		statusInfo["health"] = "pending"
	case types.ServiceStatusSuspended:
		statusInfo["message"] = "Service is temporarily suspended"
		statusInfo["health"] = "suspended"
	case types.ServiceStatusCancelled:
		statusInfo["message"] = "Service has been cancelled"
		statusInfo["health"] = "cancelled"
	case types.ServiceStatusExpired:
		statusInfo["message"] = "Service has expired"
		statusInfo["health"] = "expired"
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    statusInfo,
	})
}
