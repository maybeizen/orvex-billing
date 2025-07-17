package services

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
	"github.com/orvexcc/billing/api/internal/utils"
)

func ListInvoices(c *fiber.Ctx) error {
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
	serviceID := c.Query("service_id")

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
	
	// Filter by service if provided
	if serviceID != "" {
		query = query.Where("service_id = ?", serviceID)
	}

	// Get total count
	var total int64
	query.Model(&types.Invoice{}).Count(&total)

	// Get invoices with pagination
	var invoices []types.Invoice
	offset := (page - 1) * pageSize
	result := query.Preload("Service").
		Preload("Items").
		Order("created_at DESC").
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
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetInvoice(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	invoiceID := c.Params("id")
	if invoiceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invoice ID is required",
		})
	}

	var invoice types.Invoice
	result := db.DB.Where("id = ? AND user_id = ?", invoiceID, userID).
		Preload("Service").
		Preload("Items").
		First(&invoice)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Invoice not found",
		})
	}

	// Check if invoice is overdue
	if invoice.Status == types.InvoiceStatusPending && 
		invoice.DueDate.Before(time.Now()) {
		invoice.Status = types.InvoiceStatusOverdue
		db.DB.Save(&invoice)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    invoice,
	})
}

func PayInvoice(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	invoiceID := c.Params("id")
	if invoiceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invoice ID is required",
		})
	}

	var req types.PaymentRequest
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
	result := db.DB.Where("id = ? AND user_id = ?", invoiceID, userID).
		Preload("Service").
		First(&invoice)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Invoice not found",
		})
	}

	// Check if invoice can be paid
	if invoice.Status == types.InvoiceStatusPaid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invoice is already paid",
		})
	}

	if invoice.Status == types.InvoiceStatusCancelled {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot pay a cancelled invoice",
		})
	}

	// Process payment based on method
	switch req.PaymentMethod {
	case "stripe":
		return processStripePayment(c, &invoice, req.ReturnURL)
	case "paypal":
		return processPayPalPayment(c, &invoice, req.ReturnURL)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unsupported payment method",
		})
	}
}

func processStripePayment(c *fiber.Ctx, invoice *types.Invoice, returnURL string) error {
	// This is a placeholder for Stripe integration
	// In a real implementation, you would:
	// 1. Create a Stripe payment intent
	// 2. Return the client secret for frontend processing
	// 3. Handle webhooks for payment confirmation

	paymentURL := "https://checkout.stripe.com/pay/example"
	if returnURL != "" {
		paymentURL += "?return_url=" + returnURL
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Stripe payment initiated",
		"payment_url": paymentURL,
		"payment_method": "stripe",
		"amount": invoice.Total,
		"currency": invoice.Currency,
		"invoice_id": invoice.ID,
	})
}

func processPayPalPayment(c *fiber.Ctx, invoice *types.Invoice, returnURL string) error {
	// This is a placeholder for PayPal integration
	// In a real implementation, you would:
	// 1. Create a PayPal order
	// 2. Return the approval URL
	// 3. Handle return URLs and payment confirmation

	paymentURL := "https://www.paypal.com/checkoutnow?example"
	if returnURL != "" {
		paymentURL += "&return_url=" + returnURL
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "PayPal payment initiated",
		"payment_url": paymentURL,
		"payment_method": "paypal",
		"amount": invoice.Total,
		"currency": invoice.Currency,
		"invoice_id": invoice.ID,
	})
}

// Webhook endpoint for payment confirmations
func PaymentWebhook(c *fiber.Ctx) error {
	// This would handle webhooks from payment processors
	// to confirm successful payments and update invoice status

	var webhookData struct {
		PaymentID    string  `json:"payment_id"`
		InvoiceID    uint    `json:"invoice_id"`
		Amount       float64 `json:"amount"`
		Status       string  `json:"status"`
		PaymentMethod string `json:"payment_method"`
	}

	if err := c.BodyParser(&webhookData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid webhook data",
		})
	}

	if webhookData.Status == "completed" {
		var invoice types.Invoice
		result := db.DB.First(&invoice, webhookData.InvoiceID)
		if result.Error != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Invoice not found",
			})
		}

		// Update invoice status
		now := time.Now()
		invoice.Status = types.InvoiceStatusPaid
		invoice.PaidAt = &now
		invoice.PaymentID = webhookData.PaymentID
		invoice.PaymentMethod = webhookData.PaymentMethod

		if err := db.DB.Save(&invoice).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update invoice",
			})
		}

		// Activate associated service if pending
		if invoice.ServiceID != nil {
			var service types.Service
			if db.DB.First(&service, *invoice.ServiceID).Error == nil {
				if service.Status == types.ServiceStatusPending {
					service.Status = types.ServiceStatusActive
					now := time.Now()
					service.ActivatedAt = &now
					db.DB.Save(&service)
				}
			}
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Webhook processed",
	})
}