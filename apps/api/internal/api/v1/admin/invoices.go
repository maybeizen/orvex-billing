package admin

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
	"github.com/orvexcc/billing/api/internal/utils"
)

type AdminInvoiceRequest struct {
	UserID         uint                `json:"user_id" validate:"required"`
	ServiceID      *uint               `json:"service_id,omitempty"`
	InvoiceNumber  string              `json:"invoice_number" validate:"required,min=1,max=50"`
	Status         types.InvoiceStatus `json:"status" validate:"required,oneof=pending paid overdue cancelled refunded"`
	SubTotal       float64             `json:"subtotal" validate:"required,min=0"`
	TaxAmount      float64             `json:"tax_amount" validate:"min=0"`
	DiscountAmount float64             `json:"discount_amount" validate:"min=0"`
	Total          float64             `json:"total" validate:"required,min=0"`
	Currency       string              `json:"currency" validate:"required,len=3"`
	PaymentMethod  string              `json:"payment_method,omitempty"`
	DueDate        time.Time           `json:"due_date" validate:"required"`
	Notes          string              `json:"notes,omitempty" validate:"max=1000"`
	Items          []AdminInvoiceItem  `json:"items" validate:"required,min=1"`
}

type AdminInvoiceUpdateRequest struct {
	ServiceID      *uint               `json:"service_id,omitempty"`
	Status         types.InvoiceStatus `json:"status,omitempty" validate:"omitempty,oneof=pending paid overdue cancelled refunded"`
	SubTotal       *float64            `json:"subtotal,omitempty" validate:"omitempty,min=0"`
	TaxAmount      *float64            `json:"tax_amount,omitempty" validate:"omitempty,min=0"`
	DiscountAmount *float64            `json:"discount_amount,omitempty" validate:"omitempty,min=0"`
	Total          *float64            `json:"total,omitempty" validate:"omitempty,min=0"`
	Currency       string              `json:"currency,omitempty" validate:"omitempty,len=3"`
	PaymentMethod  string              `json:"payment_method,omitempty"`
	PaymentID      string              `json:"payment_id,omitempty"`
	DueDate        *time.Time          `json:"due_date,omitempty"`
	PaidAt         *time.Time          `json:"paid_at,omitempty"`
	Notes          string              `json:"notes,omitempty" validate:"max=1000"`
}

type AdminInvoiceItem struct {
	Description string  `json:"description" validate:"required,min=1,max=255"`
	Quantity    int     `json:"quantity" validate:"required,min=1"`
	UnitPrice   float64 `json:"unit_price" validate:"required,min=0"`
	Total       float64 `json:"total" validate:"required,min=0"`
}

type InvoiceActionRequest struct {
	Action      string     `json:"action" validate:"required,oneof=mark_paid mark_overdue cancel refund"`
	PaymentID   string     `json:"payment_id,omitempty"`
	PaymentMethod string   `json:"payment_method,omitempty"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	Reason      string     `json:"reason,omitempty" validate:"max=500"`
}

func ListAdminInvoices(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)

	// Security: Log admin access
	logAdminActivity("invoices_list", adminUser.ID, adminUser.Email, fiber.Map{
		"ip": c.IP(),
	})

	// Query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "25"))
	search := c.Query("search")
	status := c.Query("status")
	userID := c.Query("user_id")
	serviceID := c.Query("service_id")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	dueDateFrom := c.Query("due_date_from")
	dueDateTo := c.Query("due_date_to")
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
		"created_at": true, "invoice_number": true, "status": true, "total": true, "due_date": true, "paid_at": true,
	}
	if !validSorts[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" {
		sortOrder = "desc"
	}

	query := db.DB.Model(&types.Invoice{}).
		Preload("User").
		Preload("Service").
		Preload("Items")

	// Apply filters
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Joins("JOIN users ON invoices.user_id = users.id").
			Where("invoices.invoice_number ILIKE ? OR invoices.notes ILIKE ? OR users.email ILIKE ? OR users.username ILIKE ?",
				searchTerm, searchTerm, searchTerm, searchTerm)
	}

	if status != "" {
		query = query.Where("invoices.status = ?", status)
	}

	if userID != "" {
		query = query.Where("invoices.user_id = ?", userID)
	}

	if serviceID != "" {
		query = query.Where("invoices.service_id = ?", serviceID)
	}

	// Date filters
	if dateFrom != "" {
		if parsedDate, err := time.Parse("2006-01-02", dateFrom); err == nil {
			query = query.Where("invoices.created_at >= ?", parsedDate)
		}
	}
	if dateTo != "" {
		if parsedDate, err := time.Parse("2006-01-02", dateTo); err == nil {
			query = query.Where("invoices.created_at <= ?", parsedDate.Add(24*time.Hour))
		}
	}
	if dueDateFrom != "" {
		if parsedDate, err := time.Parse("2006-01-02", dueDateFrom); err == nil {
			query = query.Where("invoices.due_date >= ?", parsedDate)
		}
	}
	if dueDateTo != "" {
		if parsedDate, err := time.Parse("2006-01-02", dueDateTo); err == nil {
			query = query.Where("invoices.due_date <= ?", parsedDate.Add(24*time.Hour))
		}
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get invoices with pagination
	var invoices []types.Invoice
	offset := (page - 1) * pageSize
	orderClause := "invoices." + sortBy + " " + sortOrder
	result := query.Order(orderClause).
		Offset(offset).
		Limit(pageSize).
		Find(&invoices)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve invoices",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    invoices,
		"meta": fiber.Map{
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"sort_by":    sortBy,
			"sort_order": sortOrder,
		},
	})
}

func GetAdminInvoice(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	invoiceID := c.Params("id")

	if invoiceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invoice ID is required",
		})
	}

	var invoice types.Invoice
	result := db.DB.Preload("User").
		Preload("Service").
		Preload("Items").
		First(&invoice, invoiceID)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Invoice not found",
		})
	}

	// Security: Log admin access
	logAdminActivity("invoice_view", adminUser.ID, adminUser.Email, fiber.Map{
		"target_invoice_id": invoice.ID,
		"target_user_id":    invoice.UserID,
		"invoice_number":    invoice.InvoiceNumber,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"data":    invoice,
	})
}

func CreateAdminInvoice(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)

	var req AdminInvoiceRequest
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

	// Verify service exists if provided
	if req.ServiceID != nil {
		var service types.Service
		if result := db.DB.First(&service, *req.ServiceID); result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Service not found",
			})
		}
	}

	// Check if invoice number already exists
	var existingInvoice types.Invoice
	if result := db.DB.Where("invoice_number = ?", req.InvoiceNumber).First(&existingInvoice); result.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Invoice number already exists",
		})
	}

	// Start transaction
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create invoice
	invoice := types.Invoice{
		UserID:         req.UserID,
		ServiceID:      req.ServiceID,
		InvoiceNumber:  req.InvoiceNumber,
		Status:         req.Status,
		SubTotal:       req.SubTotal,
		TaxAmount:      req.TaxAmount,
		DiscountAmount: req.DiscountAmount,
		Total:          req.Total,
		Currency:       req.Currency,
		PaymentMethod:  req.PaymentMethod,
		DueDate:        req.DueDate,
		Notes:          req.Notes,
	}

	// Set paid_at if status is paid
	if req.Status == types.InvoiceStatusPaid {
		now := time.Now()
		invoice.PaidAt = &now
	}

	if err := tx.Create(&invoice).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create invoice",
		})
	}

	// Create invoice items
	for _, item := range req.Items {
		invoiceItem := types.InvoiceItem{
			InvoiceID:   invoice.ID,
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Total:       item.Total,
		}
		if err := tx.Create(&invoiceItem).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create invoice items",
			})
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save invoice",
		})
	}

	// Preload data
	db.DB.Preload("User").Preload("Service").Preload("Items").First(&invoice, invoice.ID)

	// Security: Log admin action
	logAdminActivity("invoice_create", adminUser.ID, adminUser.Email, fiber.Map{
		"created_invoice_id": invoice.ID,
		"target_user_id":     invoice.UserID,
		"invoice_number":     invoice.InvoiceNumber,
		"total":              invoice.Total,
		"status":             invoice.Status,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Invoice created successfully",
		"data":    invoice,
	})
}

func UpdateAdminInvoice(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	invoiceID := c.Params("id")

	if invoiceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invoice ID is required",
		})
	}

	var req AdminInvoiceUpdateRequest
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

	var invoice types.Invoice
	result := db.DB.First(&invoice, invoiceID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Invoice not found",
		})
	}

	// Track changes for logging
	changes := make(map[string]interface{})
	oldStatus := invoice.Status

	// Update fields if provided
	if req.ServiceID != nil {
		// Verify service exists
		var service types.Service
		if result := db.DB.First(&service, *req.ServiceID); result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Service not found",
			})
		}
		invoice.ServiceID = req.ServiceID
		changes["service_id"] = *req.ServiceID
	}

	if req.Status != "" && req.Status != invoice.Status {
		invoice.Status = req.Status
		changes["status"] = req.Status

		// Handle status-specific updates
		if req.Status == types.InvoiceStatusPaid && invoice.PaidAt == nil {
			if req.PaidAt != nil {
				invoice.PaidAt = req.PaidAt
			} else {
				now := time.Now()
				invoice.PaidAt = &now
			}
		} else if req.Status != types.InvoiceStatusPaid {
			invoice.PaidAt = nil
		}
	}

	if req.SubTotal != nil {
		invoice.SubTotal = *req.SubTotal
		changes["subtotal"] = *req.SubTotal
	}
	if req.TaxAmount != nil {
		invoice.TaxAmount = *req.TaxAmount
		changes["tax_amount"] = *req.TaxAmount
	}
	if req.DiscountAmount != nil {
		invoice.DiscountAmount = *req.DiscountAmount
		changes["discount_amount"] = *req.DiscountAmount
	}
	if req.Total != nil {
		invoice.Total = *req.Total
		changes["total"] = *req.Total
	}
	if req.Currency != "" {
		invoice.Currency = req.Currency
		changes["currency"] = req.Currency
	}
	if req.PaymentMethod != "" {
		invoice.PaymentMethod = req.PaymentMethod
		changes["payment_method"] = req.PaymentMethod
	}
	if req.PaymentID != "" {
		invoice.PaymentID = req.PaymentID
		changes["payment_id"] = req.PaymentID
	}
	if req.DueDate != nil {
		invoice.DueDate = *req.DueDate
		changes["due_date"] = req.DueDate
	}
	if req.PaidAt != nil {
		invoice.PaidAt = req.PaidAt
		changes["paid_at"] = req.PaidAt
	}
	if req.Notes != "" {
		invoice.Notes = req.Notes
		changes["notes"] = "updated"
	}

	if err := db.DB.Save(&invoice).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update invoice",
		})
	}

	// Preload data
	db.DB.Preload("User").Preload("Service").Preload("Items").First(&invoice, invoice.ID)

	// Security: Log admin action
	logAdminActivity("invoice_update", adminUser.ID, adminUser.Email, fiber.Map{
		"target_invoice_id": invoice.ID,
		"target_user_id":    invoice.UserID,
		"invoice_number":    invoice.InvoiceNumber,
		"old_status":        oldStatus,
		"changes":           changes,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Invoice updated successfully",
		"data":    invoice,
	})
}

func DeleteAdminInvoice(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	invoiceID := c.Params("id")

	if invoiceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invoice ID is required",
		})
	}

	var invoice types.Invoice
	result := db.DB.First(&invoice, invoiceID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Invoice not found",
		})
	}

	// Security: Prevent deletion of paid invoices without super admin
	if invoice.Status == types.InvoiceStatusPaid && adminUser.Role != types.UserRoleSuper {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only super admins can delete paid invoices",
		})
	}

	// Start transaction to delete invoice and items
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete invoice items first
	if err := tx.Where("invoice_id = ?", invoice.ID).Delete(&types.InvoiceItem{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete invoice items",
		})
	}

	// Delete invoice
	if err := tx.Delete(&invoice).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete invoice",
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to complete deletion",
		})
	}

	// Security: Log admin action
	logAdminActivity("invoice_delete", adminUser.ID, adminUser.Email, fiber.Map{
		"deleted_invoice_id": invoice.ID,
		"target_user_id":     invoice.UserID,
		"invoice_number":     invoice.InvoiceNumber,
		"total":              invoice.Total,
		"status":             invoice.Status,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Invoice deleted successfully",
	})
}

func InvoiceAction(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	invoiceID := c.Params("id")

	if invoiceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invoice ID is required",
		})
	}

	var req InvoiceActionRequest
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

	var invoice types.Invoice
	result := db.DB.First(&invoice, invoiceID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Invoice not found",
		})
	}

	oldStatus := invoice.Status

	// Perform action
	switch req.Action {
	case "mark_paid":
		if invoice.Status == types.InvoiceStatusPaid {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invoice is already paid",
			})
		}
		invoice.Status = types.InvoiceStatusPaid
		if req.PaidAt != nil {
			invoice.PaidAt = req.PaidAt
		} else {
			now := time.Now()
			invoice.PaidAt = &now
		}
		if req.PaymentID != "" {
			invoice.PaymentID = req.PaymentID
		}
		if req.PaymentMethod != "" {
			invoice.PaymentMethod = req.PaymentMethod
		}

	case "mark_overdue":
		if invoice.Status == types.InvoiceStatusPaid {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot mark paid invoice as overdue",
			})
		}
		invoice.Status = types.InvoiceStatusOverdue

	case "cancel":
		if invoice.Status == types.InvoiceStatusPaid {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot cancel paid invoice",
			})
		}
		invoice.Status = types.InvoiceStatusCancelled

	case "refund":
		if invoice.Status != types.InvoiceStatusPaid {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Only paid invoices can be refunded",
			})
		}
		// Security: Only super admins can process refunds
		if adminUser.Role != types.UserRoleSuper {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Only super admins can process refunds",
			})
		}
		invoice.Status = types.InvoiceStatusRefunded

	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid action",
		})
	}

	if err := db.DB.Save(&invoice).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update invoice",
		})
	}

	// Preload data
	db.DB.Preload("User").Preload("Service").Preload("Items").First(&invoice, invoice.ID)

	// Security: Log admin action
	logAdminActivity("invoice_action", adminUser.ID, adminUser.Email, fiber.Map{
		"target_invoice_id": invoice.ID,
		"target_user_id":    invoice.UserID,
		"invoice_number":    invoice.InvoiceNumber,
		"action":            req.Action,
		"old_status":        oldStatus,
		"new_status":        invoice.Status,
		"reason":            req.Reason,
		"payment_id":        req.PaymentID,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Invoice " + req.Action + " completed successfully",
		"data":    invoice,
	})
}