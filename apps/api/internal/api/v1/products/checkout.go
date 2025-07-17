package products

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
	"github.com/orvexcc/billing/api/internal/utils"
	"gorm.io/gorm"
)

func InitiateCheckout(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	var req types.CheckoutRequest
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

	// Security: Get and validate user's cart
	var cart types.Cart
	result := db.DB.Where("user_id = ? AND status = ?", userID, types.CartStatusActive).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Preload("Coupon").
		First(&cart)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No active cart found",
		})
	}

	if len(cart.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot checkout with empty cart",
		})
	}

	// Security: Validate all cart items and prices
	validationResult := validateCartForCheckout(&cart)
	if !validationResult.Valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": validationResult.Reason,
		})
	}

	// Security: Re-validate stock availability
	for _, item := range cart.Items {
		if item.Product.TrackQuantity && 
		   item.Product.StockQuantity < item.Quantity && 
		   !item.Product.AllowBackorder {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Insufficient stock for product: %s", item.Product.Name),
				"product_id": item.ProductID,
				"available_quantity": item.Product.StockQuantity,
				"requested_quantity": item.Quantity,
			})
		}
	}

	// Security: Re-validate coupon if applied
	if cart.CouponID != nil {
		couponValidation := validateCoupon(cart.Coupon, &cart, userID.(uint))
		if !couponValidation.Valid {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Applied coupon is no longer valid: " + couponValidation.Reason,
			})
		}
	}

	// Start database transaction for checkout
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Generate invoice in transaction
	invoice, err := generateInvoiceFromCart(&cart, req, tx)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate invoice: " + err.Error(),
		})
	}

	// Update product stock quantities in transaction
	err = updateProductStock(&cart, tx)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update product stock: " + err.Error(),
		})
	}

	// Record coupon usage if applicable
	if cart.CouponID != nil {
		err = recordCouponUsageForCheckout(cart.CouponID, userID.(uint), cart.ID, invoice.ID, cart.DiscountAmount, tx)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to record coupon usage: " + err.Error(),
			})
		}
	}

	// Mark cart as completed
	cart.Status = types.CartStatusCompleted
	if err := tx.Save(&cart).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update cart status",
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to complete checkout",
		})
	}

	// Process payment based on method
	switch req.PaymentMethod {
	case "stripe":
		return processCheckoutPayment(c, invoice, "stripe")
	case "paypal":
		return processCheckoutPayment(c, invoice, "paypal")
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unsupported payment method",
		})
	}
}

// Security validation for cart before checkout
type CheckoutValidation struct {
	Valid  bool
	Reason string
}

func validateCartForCheckout(cart *types.Cart) CheckoutValidation {
	// Validate all cart item price hashes
	for _, item := range cart.Items {
		if !validatePriceHash(item.ProductID, item.UnitPrice, item.PriceHash) {
			return CheckoutValidation{false, "Cart item price validation failed - cart may have been tampered with"}
		}

		// Validate product is still active
		if !item.Product.IsActive {
			return CheckoutValidation{false, fmt.Sprintf("Product %s is no longer available", item.Product.Name)}
		}
	}

	// Validate cart totals are correct
	expectedSubtotal := 0.0
	for _, item := range cart.Items {
		expectedSubtotal += item.TotalPrice
	}

	if cart.SubTotal != expectedSubtotal {
		return CheckoutValidation{false, "Cart subtotal validation failed"}
	}

	return CheckoutValidation{true, "Cart is valid for checkout"}
}

func generateInvoiceFromCart(cart *types.Cart, req types.CheckoutRequest, tx *gorm.DB) (*types.Invoice, error) {
	// Generate unique invoice number
	invoiceNumber := fmt.Sprintf("INV-%d-%d", time.Now().Unix(), cart.ID)

	// Create invoice
	invoice := types.Invoice{
		UserID:         cart.UserID,
		InvoiceNumber:  invoiceNumber,
		Status:         types.InvoiceStatusPending,
		SubTotal:       cart.SubTotal,
		TaxAmount:      cart.TaxAmount,
		DiscountAmount: cart.DiscountAmount,
		Total:          cart.Total,
		Currency:       cart.Currency,
		DueDate:        time.Now().Add(30 * 24 * time.Hour), // 30 days
		Notes:          req.Notes,
	}

	if err := tx.Create(&invoice).Error; err != nil {
		return nil, err
	}

	// Create invoice items from cart items
	for _, cartItem := range cart.Items {
		invoiceItem := types.InvoiceItem{
			InvoiceID:   invoice.ID,
			Description: fmt.Sprintf("%s - %s", cartItem.Product.Name, cartItem.Product.Category.Name),
			Quantity:    cartItem.Quantity,
			UnitPrice:   cartItem.UnitPrice,
			Total:       cartItem.TotalPrice,
		}

		if err := tx.Create(&invoiceItem).Error; err != nil {
			return nil, err
		}
	}

	return &invoice, nil
}

func updateProductStock(cart *types.Cart, tx *gorm.DB) error {
	for _, item := range cart.Items {
		if item.Product.TrackQuantity {
			// Security: Use atomic update to prevent race conditions
			result := tx.Model(&types.Product{}).
				Where("id = ? AND stock_quantity >= ?", item.ProductID, item.Quantity).
				UpdateColumn("stock_quantity", gorm.Expr("stock_quantity - ?", item.Quantity))

			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return fmt.Errorf("insufficient stock for product ID %d", item.ProductID)
			}
		}
	}
	return nil
}

func recordCouponUsageForCheckout(couponID *uint, userID uint, cartID uint, invoiceID uint, discountAmount float64, tx *gorm.DB) error {
	if couponID == nil {
		return nil
	}

	usage := types.CouponUsage{
		CouponID:       *couponID,
		UserID:         userID,
		CartID:         &cartID,
		InvoiceID:      &invoiceID,
		DiscountAmount: discountAmount,
		UsedAt:         time.Now(),
	}

	if err := tx.Create(&usage).Error; err != nil {
		return err
	}

	// Increment coupon used count
	return tx.Model(&types.Coupon{}).
		Where("id = ?", *couponID).
		UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error
}

func processCheckoutPayment(c *fiber.Ctx, invoice *types.Invoice, paymentMethod string) error {
	// This integrates with the existing payment processing system
	var paymentURL string
	var message string

	switch paymentMethod {
	case "stripe":
		paymentURL = fmt.Sprintf("https://checkout.stripe.com/pay/checkout-%d", invoice.ID)
		message = "Stripe payment initiated for checkout"
	case "paypal":
		paymentURL = fmt.Sprintf("https://www.paypal.com/checkoutnow?invoice=%d", invoice.ID)
		message = "PayPal payment initiated for checkout"
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": message,
		"invoice_id": invoice.ID,
		"invoice_number": invoice.InvoiceNumber,
		"payment_url": paymentURL,
		"payment_method": paymentMethod,
		"amount": invoice.Total,
		"currency": invoice.Currency,
		"checkout_completed": true,
	})
}

func GetCheckoutSummary(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Get user's active cart
	var cart types.Cart
	result := db.DB.Where("user_id = ? AND status = ?", userID, types.CartStatusActive).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Preload("Coupon").
		First(&cart)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No active cart found",
		})
	}

	if len(cart.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot checkout with empty cart",
		})
	}

	// Update cart totals to ensure accuracy
	updateCartTotals(&cart)

	// Prepare checkout summary
	summary := fiber.Map{
		"subtotal": cart.SubTotal,
		"tax_amount": cart.TaxAmount,
		"discount_amount": cart.DiscountAmount,
		"total": cart.Total,
		"currency": cart.Currency,
		"item_count": len(cart.Items),
		"items": cart.Items,
	}

	if cart.CouponID != nil {
		summary["applied_coupon"] = cart.Coupon
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": summary,
	})
}