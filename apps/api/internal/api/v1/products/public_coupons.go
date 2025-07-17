package products

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
)

// ListPublicCoupons returns public coupons that users can discover
func ListPublicCoupons(c *fiber.Ctx) error {
	// Query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	couponType := c.Query("type")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 20 {
		pageSize = 10
	}

	now := time.Now()
	query := db.DB.Model(&types.Coupon{}).
		Where("is_active = ? AND is_public = ?", true, true).
		Where("(starts_at IS NULL OR starts_at <= ?)", now).
		Where("(expires_at IS NULL OR expires_at > ?)", now).
		Where("(usage_limit IS NULL OR used_count < usage_limit)")

	// Apply type filter
	if couponType != "" {
		query = query.Where("type = ?", couponType)
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get public coupons
	var coupons []types.Coupon
	offset := (page - 1) * pageSize
	result := query.Select("id, code, name, description, type, value, min_order_amount, max_discount_amount, expires_at").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&coupons)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve coupons",
		})
	}

	// Remove sensitive fields for public display
	for i := range coupons {
		coupons[i].UsageLimit = nil
		coupons[i].UserUsageLimit = nil
		coupons[i].UsedCount = 0
		coupons[i].ApplicableProductIDs = nil
		coupons[i].ApplicableCategoryIDs = nil
		coupons[i].ExcludedProductIDs = nil
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    coupons,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetPublicCouponByCode allows users to check a specific public coupon
func GetPublicCouponByCode(c *fiber.Ctx) error {
	couponCode := c.Params("code")
	if couponCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Coupon code is required",
		})
	}

	var coupon types.Coupon
	now := time.Now()
	result := db.DB.Where("code = ? AND is_active = ? AND is_public = ?", couponCode, true, true).
		Where("(starts_at IS NULL OR starts_at <= ?)", now).
		Where("(expires_at IS NULL OR expires_at > ?)", now).
		Where("(usage_limit IS NULL OR used_count < usage_limit)").
		Select("id, code, name, description, type, value, min_order_amount, max_discount_amount, expires_at").
		First(&coupon)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Coupon not found or not available",
		})
	}

	// Remove sensitive fields
	coupon.UsageLimit = nil
	coupon.UserUsageLimit = nil
	coupon.UsedCount = 0
	coupon.ApplicableProductIDs = nil
	coupon.ApplicableCategoryIDs = nil
	coupon.ExcludedProductIDs = nil

	return c.JSON(fiber.Map{
		"success": true,
		"data":    coupon,
	})
}