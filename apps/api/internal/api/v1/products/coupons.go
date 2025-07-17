package products

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
	"github.com/orvexcc/billing/api/internal/utils"
)

func ApplyCoupon(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	var req types.ApplyCouponRequest
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

	// Security: Normalize coupon code
	couponCode := strings.TrimSpace(strings.ToUpper(req.CouponCode))
	if len(couponCode) < 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid coupon code",
		})
	}

	// Get user's active cart
	var cart types.Cart
	result := db.DB.Where("user_id = ? AND status = ?", userID, types.CartStatusActive).
		Preload("Items").
		Preload("Items.Product").
		First(&cart)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No active cart found",
		})
	}

	if len(cart.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot apply coupon to empty cart",
		})
	}

	// Security: Find and validate coupon
	var coupon types.Coupon
	result = db.DB.Where("code = ? AND is_active = ?", couponCode, true).First(&coupon)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Invalid or expired coupon code",
		})
	}

	// Validate coupon with comprehensive security checks
	validationResult := validateCoupon(&coupon, &cart, userID.(uint))
	if !validationResult.Valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": validationResult.Reason,
		})
	}

	// Remove existing coupon if any
	if cart.CouponID != nil {
		cart.CouponID = nil
		cart.CouponCode = ""
		cart.DiscountAmount = 0
	}

	// Apply new coupon
	cart.CouponID = &coupon.ID
	cart.CouponCode = coupon.Code

	// Update cart totals with new coupon
	updateCartTotals(&cart)

	// Reload cart with coupon details
	db.DB.Where("id = ?", cart.ID).Preload("Coupon").First(&cart)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Coupon applied successfully",
		"discount_amount": cart.DiscountAmount,
		"new_total": cart.Total,
		"coupon": cart.Coupon,
	})
}

func RemoveCoupon(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Get user's active cart
	var cart types.Cart
	result := db.DB.Where("user_id = ? AND status = ?", userID, types.CartStatusActive).First(&cart)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No active cart found",
		})
	}

	if cart.CouponID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No coupon applied to cart",
		})
	}

	// Remove coupon
	cart.CouponID = nil
	cart.CouponCode = ""
	cart.DiscountAmount = 0

	// Update cart totals
	updateCartTotals(&cart)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Coupon removed successfully",
		"new_total": cart.Total,
	})
}

func ValidateCouponCode(c *fiber.Ctx) error {
	couponCode := c.Params("code")
	if couponCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Coupon code is required",
		})
	}

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	// Security: Normalize coupon code
	couponCode = strings.TrimSpace(strings.ToUpper(couponCode))

	// Find coupon
	var coupon types.Coupon
	result := db.DB.Where("code = ? AND is_active = ?", couponCode, true).First(&coupon)
	if result.Error != nil {
		return c.JSON(fiber.Map{
			"valid": false,
			"reason": "Coupon not found or inactive",
		})
	}

	// Get user's cart for validation context
	var cart types.Cart
	result = db.DB.Where("user_id = ? AND status = ?", userID, types.CartStatusActive).
		Preload("Items").
		Preload("Items.Product").
		First(&cart)

	if result.Error != nil {
		return c.JSON(fiber.Map{
			"valid": false,
			"reason": "No active cart found",
		})
	}

	// Validate coupon
	validationResult := validateCoupon(&coupon, &cart, userID.(uint))

	response := fiber.Map{
		"valid": validationResult.Valid,
		"reason": validationResult.Reason,
	}

	if validationResult.Valid {
		// Calculate potential discount
		discount := calculateCouponDiscount(&coupon, cart.SubTotal)
		response["discount_amount"] = discount
		response["coupon"] = fiber.Map{
			"code": coupon.Code,
			"name": coupon.Name,
			"description": coupon.Description,
			"type": coupon.Type,
			"value": coupon.Value,
		}
	}

	return c.JSON(response)
}

// Comprehensive coupon validation structure
type CouponValidation struct {
	Valid  bool
	Reason string
}

// Security-focused coupon validation function
func validateCoupon(coupon *types.Coupon, cart *types.Cart, userID uint) CouponValidation {
	now := time.Now()

	// Check if coupon is active
	if !coupon.IsActive {
		return CouponValidation{false, "Coupon is not active"}
	}

	// Check date restrictions
	if coupon.StartsAt != nil && now.Before(*coupon.StartsAt) {
		return CouponValidation{false, "Coupon is not yet valid"}
	}

	if coupon.ExpiresAt != nil && now.After(*coupon.ExpiresAt) {
		return CouponValidation{false, "Coupon has expired"}
	}

	// Check minimum order amount
	if coupon.MinOrderAmount != nil && cart.SubTotal < *coupon.MinOrderAmount {
		return CouponValidation{false, "Order amount does not meet minimum requirement"}
	}

	// Check global usage limit
	if coupon.UsageLimit != nil && coupon.UsedCount >= *coupon.UsageLimit {
		return CouponValidation{false, "Coupon usage limit has been reached"}
	}

	// Check user-specific usage limit
	if coupon.UserUsageLimit != nil {
		var userUsageCount int64
		db.DB.Model(&types.CouponUsage{}).
			Where("coupon_id = ? AND user_id = ?", coupon.ID, userID).
			Count(&userUsageCount)

		if int(userUsageCount) >= *coupon.UserUsageLimit {
			return CouponValidation{false, "You have reached the usage limit for this coupon"}
		}
	}

	// Check product/category restrictions
	if len(coupon.ApplicableProductIDs) > 0 {
		hasApplicableProduct := false
		for _, item := range cart.Items {
			for _, productID := range coupon.ApplicableProductIDs {
				if item.ProductID == productID {
					hasApplicableProduct = true
					break
				}
			}
			if hasApplicableProduct {
				break
			}
		}
		if !hasApplicableProduct {
			return CouponValidation{false, "Coupon is not applicable to items in your cart"}
		}
	}

	if len(coupon.ApplicableCategoryIDs) > 0 {
		hasApplicableCategory := false
		for _, item := range cart.Items {
			for _, categoryID := range coupon.ApplicableCategoryIDs {
				if item.Product.CategoryID == categoryID {
					hasApplicableCategory = true
					break
				}
			}
			if hasApplicableCategory {
				break
			}
		}
		if !hasApplicableCategory {
			return CouponValidation{false, "Coupon is not applicable to product categories in your cart"}
		}
	}

	// Check excluded products
	if len(coupon.ExcludedProductIDs) > 0 {
		for _, item := range cart.Items {
			for _, excludedProductID := range coupon.ExcludedProductIDs {
				if item.ProductID == excludedProductID {
					return CouponValidation{false, "Coupon cannot be applied to one or more items in your cart"}
				}
			}
		}
	}

	return CouponValidation{true, "Coupon is valid"}
}

// Security function to record coupon usage
func recordCouponUsage(couponID uint, userID uint, cartID uint, discountAmount float64) error {
	usage := types.CouponUsage{
		CouponID:       couponID,
		UserID:         userID,
		CartID:         &cartID,
		DiscountAmount: discountAmount,
		UsedAt:         time.Now(),
	}

	if err := db.DB.Create(&usage).Error; err != nil {
		return err
	}

	// Increment coupon used count
	return db.DB.Model(&types.Coupon{}).
		Where("id = ?", couponID).
		UpdateColumn("used_count", db.DB.Raw("used_count + 1")).Error
}