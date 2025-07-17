package admin

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
	"github.com/orvexcc/billing/api/internal/utils"
)

type AdminServiceRequest struct {
	UserID       uint                `json:"user_id" validate:"required"`
	Name         string              `json:"name" validate:"required,min=1,max=255"`
	Description  string              `json:"description,omitempty" validate:"max=1000"`
	Status       types.ServiceStatus `json:"status" validate:"required,oneof=pending active suspended cancelled expired"`
	Price        float64             `json:"price" validate:"required,min=0"`
	Currency     string              `json:"currency" validate:"required,len=3"`
	BillingType  types.BillingType   `json:"billing_type" validate:"required,oneof=one_time monthly yearly"`
	Config       types.ServiceConfig `json:"config,omitempty"`
	NextBillDate *time.Time          `json:"next_bill_date,omitempty"`
	Notes        string              `json:"notes,omitempty" validate:"max=500"`
}

type AdminServiceUpdateRequest struct {
	Name         string              `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description  string              `json:"description,omitempty" validate:"max=1000"`
	Status       types.ServiceStatus `json:"status,omitempty" validate:"omitempty,oneof=pending active suspended cancelled expired"`
	Price        *float64            `json:"price,omitempty" validate:"omitempty,min=0"`
	Currency     string              `json:"currency,omitempty" validate:"omitempty,len=3"`
	BillingType  types.BillingType   `json:"billing_type,omitempty" validate:"omitempty,oneof=one_time monthly yearly"`
	Config       *types.ServiceConfig `json:"config,omitempty"`
	NextBillDate *time.Time          `json:"next_bill_date,omitempty"`
	Notes        string              `json:"notes,omitempty" validate:"max=500"`
}

type ServiceActionRequest struct {
	Action string `json:"action" validate:"required,oneof=activate suspend terminate cancel"`
	Reason string `json:"reason,omitempty" validate:"max=500"`
}

func ListAdminServices(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)

	// Security: Log admin access
	logAdminActivity("services_list", adminUser.ID, adminUser.Email, fiber.Map{
		"ip": c.IP(),
	})

	// Query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "25"))
	search := c.Query("search")
	status := c.Query("status")
	userID := c.Query("user_id")
	billingType := c.Query("billing_type")
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
		"created_at": true, "name": true, "status": true, "price": true, "next_bill_date": true,
	}
	if !validSorts[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" {
		sortOrder = "desc"
	}

	query := db.DB.Model(&types.Service{}).Preload("User")

	// Apply filters
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Joins("JOIN users ON services.user_id = users.id").
			Where("services.name ILIKE ? OR services.description ILIKE ? OR users.email ILIKE ? OR users.username ILIKE ?",
				searchTerm, searchTerm, searchTerm, searchTerm)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if billingType != "" {
		query = query.Where("billing_type = ?", billingType)
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get services with pagination
	var services []types.Service
	offset := (page - 1) * pageSize
	orderClause := sortBy + " " + sortOrder
	result := query.Order(orderClause).
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
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"sort_by":    sortBy,
			"sort_order": sortOrder,
		},
	})
}

func GetAdminService(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	serviceID := c.Params("id")

	if serviceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Service ID is required",
		})
	}

	var service types.Service
	result := db.DB.Preload("User").
		Preload("Invoices").
		Preload("Invoices.Items").
		First(&service, serviceID)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Service not found",
		})
	}

	// Security: Log admin access
	logAdminActivity("service_view", adminUser.ID, adminUser.Email, fiber.Map{
		"target_service_id": service.ID,
		"target_user_id":    service.UserID,
		"service_name":      service.Name,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"data":    service,
	})
}

func CreateAdminService(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)

	var req AdminServiceRequest
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

	// Verify user exists
	var user types.User
	if result := db.DB.First(&user, req.UserID); result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Create service
	service := types.Service{
		UserID:       req.UserID,
		Name:         req.Name,
		Description:  req.Description,
		Status:       req.Status,
		Price:        req.Price,
		Currency:     req.Currency,
		BillingType:  req.BillingType,
		Config:       req.Config,
		NextBillDate: req.NextBillDate,
	}

	// Set status timestamps
	now := time.Now()
	switch req.Status {
	case types.ServiceStatusActive:
		service.ActivatedAt = &now
	case types.ServiceStatusSuspended:
		service.SuspendedAt = &now
	case types.ServiceStatusCancelled:
		service.CancelledAt = &now
	}

	result := db.DB.Create(&service)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create service",
		})
	}

	// Preload user data
	db.DB.Preload("User").First(&service, service.ID)

	// Security: Log admin action
	logAdminActivity("service_create", adminUser.ID, adminUser.Email, fiber.Map{
		"created_service_id": service.ID,
		"target_user_id":     service.UserID,
		"service_name":       service.Name,
		"status":             service.Status,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Service created successfully",
		"data":    service,
	})
}

func UpdateAdminService(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	serviceID := c.Params("id")

	if serviceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Service ID is required",
		})
	}

	var req AdminServiceUpdateRequest
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

	var service types.Service
	result := db.DB.First(&service, serviceID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Service not found",
		})
	}

	// Track changes for logging
	changes := make(map[string]interface{})
	oldStatus := service.Status

	// Update fields if provided
	if req.Name != "" {
		service.Name = req.Name
		changes["name"] = req.Name
	}
	if req.Description != "" {
		service.Description = req.Description
		changes["description"] = req.Description
	}
	if req.Status != "" && req.Status != service.Status {
		service.Status = req.Status
		changes["status"] = req.Status

		// Update status timestamps
		now := time.Now()
		switch req.Status {
		case types.ServiceStatusActive:
			if oldStatus != types.ServiceStatusActive {
				service.ActivatedAt = &now
			}
		case types.ServiceStatusSuspended:
			if oldStatus != types.ServiceStatusSuspended {
				service.SuspendedAt = &now
			}
		case types.ServiceStatusCancelled:
			if oldStatus != types.ServiceStatusCancelled {
				service.CancelledAt = &now
			}
		}
	}
	if req.Price != nil {
		service.Price = *req.Price
		changes["price"] = *req.Price
	}
	if req.Currency != "" {
		service.Currency = req.Currency
		changes["currency"] = req.Currency
	}
	if req.BillingType != "" {
		service.BillingType = req.BillingType
		changes["billing_type"] = req.BillingType
	}
	if req.Config != nil {
		service.Config = *req.Config
		changes["config"] = "updated"
	}
	if req.NextBillDate != nil {
		service.NextBillDate = req.NextBillDate
		changes["next_bill_date"] = req.NextBillDate
	}

	if err := db.DB.Save(&service).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update service",
		})
	}

	// Preload user data
	db.DB.Preload("User").First(&service, service.ID)

	// Security: Log admin action
	logAdminActivity("service_update", adminUser.ID, adminUser.Email, fiber.Map{
		"target_service_id": service.ID,
		"target_user_id":    service.UserID,
		"service_name":      service.Name,
		"changes":           changes,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Service updated successfully",
		"data":    service,
	})
}

func DeleteAdminService(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	serviceID := c.Params("id")

	if serviceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Service ID is required",
		})
	}

	var service types.Service
	result := db.DB.First(&service, serviceID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Service not found",
		})
	}

	// Soft delete service
	if err := db.DB.Delete(&service).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete service",
		})
	}

	// Security: Log admin action
	logAdminActivity("service_delete", adminUser.ID, adminUser.Email, fiber.Map{
		"deleted_service_id": service.ID,
		"target_user_id":     service.UserID,
		"service_name":       service.Name,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Service deleted successfully",
	})
}

func ServiceAction(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	serviceID := c.Params("id")

	if serviceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Service ID is required",
		})
	}

	var req ServiceActionRequest
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

	var service types.Service
	result := db.DB.First(&service, serviceID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Service not found",
		})
	}

	now := time.Now()
	oldStatus := service.Status

	// Perform action
	switch req.Action {
	case "activate":
		if service.Status == types.ServiceStatusActive {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Service is already active",
			})
		}
		service.Status = types.ServiceStatusActive
		service.ActivatedAt = &now
		service.SuspendedAt = nil
		service.CancelledAt = nil

	case "suspend":
		if service.Status == types.ServiceStatusSuspended {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Service is already suspended",
			})
		}
		if service.Status == types.ServiceStatusCancelled {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot suspend cancelled service",
			})
		}
		service.Status = types.ServiceStatusSuspended
		service.SuspendedAt = &now

	case "terminate", "cancel":
		if service.Status == types.ServiceStatusCancelled {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Service is already cancelled",
			})
		}
		service.Status = types.ServiceStatusCancelled
		service.CancelledAt = &now

	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid action",
		})
	}

	if err := db.DB.Save(&service).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update service",
		})
	}

	// Preload user data
	db.DB.Preload("User").First(&service, service.ID)

	// Security: Log admin action
	logAdminActivity("service_action", adminUser.ID, adminUser.Email, fiber.Map{
		"target_service_id": service.ID,
		"target_user_id":    service.UserID,
		"service_name":      service.Name,
		"action":            req.Action,
		"old_status":        oldStatus,
		"new_status":        service.Status,
		"reason":            req.Reason,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Service " + req.Action + " completed successfully",
		"data":    service,
	})
}