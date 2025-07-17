package admin

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
	"github.com/orvexcc/billing/api/internal/utils"
)

type AdminCouponRequest struct {
	Code                  string             `json:"code" validate:"required,min=3,max=50,alphanum"`
	Name                  string             `json:"name" validate:"required,min=1,max=255"`
	Description           string             `json:"description,omitempty" validate:"max=1000"`
	Type                  types.CouponType   `json:"type" validate:"required,oneof=percentage fixed free_shipping"`
	Value                 float64            `json:"value" validate:"required,min=0"`
	MinOrderAmount        *float64           `json:"min_order_amount,omitempty" validate:"omitempty,min=0"`
	MaxDiscountAmount     *float64           `json:"max_discount_amount,omitempty" validate:"omitempty,min=0"`
	UsageLimit            *int               `json:"usage_limit,omitempty" validate:"omitempty,min=1"`
	UserUsageLimit        *int               `json:"user_usage_limit,omitempty" validate:"omitempty,min=1"`
	StartsAt              *time.Time         `json:"starts_at,omitempty"`
	ExpiresAt             *time.Time         `json:"expires_at,omitempty"`
	IsActive              bool               `json:"is_active"`
	IsPublic              bool               `json:"is_public"`
	ApplicableProductIDs  []uint             `json:"applicable_product_ids,omitempty"`
	ApplicableCategoryIDs []uint             `json:"applicable_category_ids,omitempty"`
	ExcludedProductIDs    []uint             `json:"excluded_product_ids,omitempty"`
}

type AdminCouponUpdateRequest struct {
	Code                  string             `json:"code,omitempty" validate:"omitempty,min=3,max=50,alphanum"`
	Name                  string             `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description           string             `json:"description,omitempty" validate:"max=1000"`
	Type                  types.CouponType   `json:"type,omitempty" validate:"omitempty,oneof=percentage fixed free_shipping"`
	Value                 *float64           `json:"value,omitempty" validate:"omitempty,min=0"`
	MinOrderAmount        *float64           `json:"min_order_amount,omitempty" validate:"omitempty,min=0"`
	MaxDiscountAmount     *float64           `json:"max_discount_amount,omitempty" validate:"omitempty,min=0"`
	UsageLimit            *int               `json:"usage_limit,omitempty" validate:"omitempty,min=1"`
	UserUsageLimit        *int               `json:"user_usage_limit,omitempty" validate:"omitempty,min=1"`
	StartsAt              *time.Time         `json:"starts_at,omitempty"`
	ExpiresAt             *time.Time         `json:"expires_at,omitempty"`
	IsActive              *bool              `json:"is_active,omitempty"`
	IsPublic              *bool              `json:"is_public,omitempty"`
	ApplicableProductIDs  []uint             `json:"applicable_product_ids,omitempty"`
	ApplicableCategoryIDs []uint             `json:"applicable_category_ids,omitempty"`
	ExcludedProductIDs    []uint             `json:"excluded_product_ids,omitempty"`
}

type CouponUsageStats struct {
	TotalUsage       int64   `json:"total_usage"`
	UniqueUsers      int64   `json:"unique_users"`
	TotalDiscount    float64 `json:"total_discount"`
	AvgDiscountValue float64 `json:"avg_discount_value"`
	RecentUsage      []types.CouponUsage `json:"recent_usage"`
}

func ListAdminCoupons(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)

	// Security: Log admin access
	logAdminActivity("coupons_list", adminUser.ID, adminUser.Email, fiber.Map{
		"ip": c.IP(),
	})

	// Query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "25"))
	search := c.Query("search")
	couponType := c.Query("type")
	isActive := c.Query("is_active")
	isPublic := c.Query("is_public")
	status := c.Query("status") // active, expired, upcoming, depleted
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
		"created_at": true, "code": true, "name": true, "value": true, "used_count": true, "expires_at": true,
	}
	if !validSorts[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" {
		sortOrder = "desc"
	}

	query := db.DB.Model(&types.Coupon{})

	// Apply filters
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("code ILIKE ? OR name ILIKE ? OR description ILIKE ?",
			searchTerm, searchTerm, searchTerm)
	}

	if couponType != "" {
		query = query.Where("type = ?", couponType)
	}

	if isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	if isPublic != "" {
		query = query.Where("is_public = ?", isPublic == "true")
	}

	// Status filtering
	now := time.Now()
	switch status {
	case "active":
		query = query.Where("is_active = ? AND (starts_at IS NULL OR starts_at <= ?) AND (expires_at IS NULL OR expires_at > ?) AND (usage_limit IS NULL OR used_count < usage_limit)",
			true, now, now)
	case "expired":
		query = query.Where("expires_at IS NOT NULL AND expires_at <= ?", now)
	case "upcoming":
		query = query.Where("starts_at IS NOT NULL AND starts_at > ?", now)
	case "depleted":
		query = query.Where("usage_limit IS NOT NULL AND used_count >= usage_limit")
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get coupons with pagination
	var coupons []types.Coupon
	offset := (page - 1) * pageSize
	orderClause := sortBy + " " + sortOrder
	result := query.Order(orderClause).
		Offset(offset).
		Limit(pageSize).
		Find(&coupons)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve coupons",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    coupons,
		"meta": fiber.Map{
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"sort_by":    sortBy,
			"sort_order": sortOrder,
		},
	})
}

func GetAdminCoupon(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	couponID := c.Params("id")

	if couponID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Coupon ID is required",
		})
	}

	var coupon types.Coupon
	result := db.DB.First(&coupon, couponID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Coupon not found",
		})
	}

	// Get usage statistics
	var stats CouponUsageStats

	// Total usage count and unique users
	db.DB.Model(&types.CouponUsage{}).Where("coupon_id = ?", coupon.ID).Count(&stats.TotalUsage)
	db.DB.Model(&types.CouponUsage{}).Where("coupon_id = ?", coupon.ID).Distinct("user_id").Count(&stats.UniqueUsers)

	// Total discount amount
	var discountResult struct {
		TotalDiscount float64
		AvgDiscount   float64
	}
	db.DB.Model(&types.CouponUsage{}).
		Select("COALESCE(SUM(discount_amount), 0) as total_discount, COALESCE(AVG(discount_amount), 0) as avg_discount").
		Where("coupon_id = ?", coupon.ID).
		Scan(&discountResult)
	stats.TotalDiscount = discountResult.TotalDiscount
	stats.AvgDiscountValue = discountResult.AvgDiscount

	// Recent usage
	db.DB.Where("coupon_id = ?", coupon.ID).
		Preload("User").
		Order("used_at DESC").
		Limit(10).
		Find(&stats.RecentUsage)

	// Security: Log admin access
	logAdminActivity("coupon_view", adminUser.ID, adminUser.Email, fiber.Map{
		"target_coupon_id": coupon.ID,
		"coupon_code":      coupon.Code,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"coupon": coupon,
			"stats":  stats,
		},
	})
}

func CreateAdminCoupon(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)

	var req AdminCouponRequest
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

	// Validate percentage coupon
	if req.Type == types.CouponTypePercentage && req.Value > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Percentage discount cannot exceed 100%",
		})
	}

	// Validate date range
	if req.StartsAt != nil && req.ExpiresAt != nil && req.StartsAt.After(*req.ExpiresAt) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Start date cannot be after expiration date",
		})
	}

	// Check if coupon code already exists
	var existingCoupon types.Coupon
	if result := db.DB.Where("code = ?", req.Code).First(&existingCoupon); result.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Coupon code already exists",
		})
	}

	// Validate product and category restrictions
	if len(req.ApplicableProductIDs) > 0 {
		var productCount int64
		db.DB.Model(&types.Product{}).Where("id IN ?", req.ApplicableProductIDs).Count(&productCount)
		if productCount != int64(len(req.ApplicableProductIDs)) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Some specified products do not exist",
			})
		}
	}

	if len(req.ApplicableCategoryIDs) > 0 {
		var categoryCount int64
		db.DB.Model(&types.Category{}).Where("id IN ?", req.ApplicableCategoryIDs).Count(&categoryCount)
		if categoryCount != int64(len(req.ApplicableCategoryIDs)) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Some specified categories do not exist",
			})
		}
	}

	// Create coupon
	coupon := types.Coupon{
		Code:                  req.Code,
		Name:                  req.Name,
		Description:           req.Description,
		Type:                  req.Type,
		Value:                 req.Value,
		MinOrderAmount:        req.MinOrderAmount,
		MaxDiscountAmount:     req.MaxDiscountAmount,
		UsageLimit:            req.UsageLimit,
		UserUsageLimit:        req.UserUsageLimit,
		StartsAt:              req.StartsAt,
		ExpiresAt:             req.ExpiresAt,
		IsActive:              req.IsActive,
		IsPublic:              req.IsPublic,
		ApplicableProductIDs:  req.ApplicableProductIDs,
		ApplicableCategoryIDs: req.ApplicableCategoryIDs,
		ExcludedProductIDs:    req.ExcludedProductIDs,
	}

	result := db.DB.Create(&coupon)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create coupon",
		})
	}

	// Security: Log admin action
	logAdminActivity("coupon_create", adminUser.ID, adminUser.Email, fiber.Map{
		"created_coupon_id": coupon.ID,
		"coupon_code":       coupon.Code,
		"coupon_type":       coupon.Type,
		"value":             coupon.Value,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Coupon created successfully",
		"data":    coupon,
	})
}

func UpdateAdminCoupon(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	couponID := c.Params("id")

	if couponID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Coupon ID is required",
		})
	}

	var req AdminCouponUpdateRequest
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

	var coupon types.Coupon
	result := db.DB.First(&coupon, couponID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Coupon not found",
		})
	}

	// Check if coupon has been used and prevent certain changes
	if coupon.UsedCount > 0 {
		// Prevent changing critical fields if coupon has been used
		if req.Code != "" && req.Code != coupon.Code {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot change coupon code after it has been used",
			})
		}
		if req.Type != "" && req.Type != coupon.Type {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot change coupon type after it has been used",
			})
		}
		if req.Value != nil && *req.Value != coupon.Value {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot change coupon value after it has been used",
			})
		}
	}

	// Track changes for logging
	changes := make(map[string]interface{})

	// Update fields if provided
	if req.Code != "" && req.Code != coupon.Code {
		// Check if new code already exists
		var existingCoupon types.Coupon
		if result := db.DB.Where("code = ? AND id != ?", req.Code, coupon.ID).First(&existingCoupon); result.Error == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Coupon code already exists",
			})
		}
		coupon.Code = req.Code
		changes["code"] = req.Code
	}

	if req.Name != "" {
		coupon.Name = req.Name
		changes["name"] = req.Name
	}
	if req.Description != "" {
		coupon.Description = req.Description
		changes["description"] = "updated"
	}
	if req.Type != "" {
		coupon.Type = req.Type
		changes["type"] = req.Type
	}
	if req.Value != nil {
		// Validate percentage coupon
		if coupon.Type == types.CouponTypePercentage && *req.Value > 100 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Percentage discount cannot exceed 100%",
			})
		}
		coupon.Value = *req.Value
		changes["value"] = *req.Value
	}
	if req.MinOrderAmount != nil {
		coupon.MinOrderAmount = req.MinOrderAmount
		changes["min_order_amount"] = *req.MinOrderAmount
	}
	if req.MaxDiscountAmount != nil {
		coupon.MaxDiscountAmount = req.MaxDiscountAmount
		changes["max_discount_amount"] = *req.MaxDiscountAmount
	}
	if req.UsageLimit != nil {
		coupon.UsageLimit = req.UsageLimit
		changes["usage_limit"] = *req.UsageLimit
	}
	if req.UserUsageLimit != nil {
		coupon.UserUsageLimit = req.UserUsageLimit
		changes["user_usage_limit"] = *req.UserUsageLimit
	}
	if req.StartsAt != nil {
		coupon.StartsAt = req.StartsAt
		changes["starts_at"] = req.StartsAt
	}
	if req.ExpiresAt != nil {
		coupon.ExpiresAt = req.ExpiresAt
		changes["expires_at"] = req.ExpiresAt
	}
	if req.IsActive != nil {
		coupon.IsActive = *req.IsActive
		changes["is_active"] = *req.IsActive
	}
	if req.IsPublic != nil {
		coupon.IsPublic = *req.IsPublic
		changes["is_public"] = *req.IsPublic
	}
	if len(req.ApplicableProductIDs) > 0 {
		coupon.ApplicableProductIDs = req.ApplicableProductIDs
		changes["applicable_products"] = "updated"
	}
	if len(req.ApplicableCategoryIDs) > 0 {
		coupon.ApplicableCategoryIDs = req.ApplicableCategoryIDs
		changes["applicable_categories"] = "updated"
	}
	if len(req.ExcludedProductIDs) > 0 {
		coupon.ExcludedProductIDs = req.ExcludedProductIDs
		changes["excluded_products"] = "updated"
	}

	// Validate date range
	if coupon.StartsAt != nil && coupon.ExpiresAt != nil && coupon.StartsAt.After(*coupon.ExpiresAt) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Start date cannot be after expiration date",
		})
	}

	if err := db.DB.Save(&coupon).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update coupon",
		})
	}

	// Security: Log admin action
	logAdminActivity("coupon_update", adminUser.ID, adminUser.Email, fiber.Map{
		"target_coupon_id": coupon.ID,
		"coupon_code":      coupon.Code,
		"changes":          changes,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Coupon updated successfully",
		"data":    coupon,
	})
}

func DeleteAdminCoupon(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	couponID := c.Params("id")

	if couponID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Coupon ID is required",
		})
	}

	var coupon types.Coupon
	result := db.DB.First(&coupon, couponID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Coupon not found",
		})
	}

	// Check if coupon is currently applied to any active carts
	var activeCartCount int64
	db.DB.Model(&types.Cart{}).Where("coupon_id = ? AND status = ?", coupon.ID, types.CartStatusActive).Count(&activeCartCount)
	if activeCartCount > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete coupon that is currently applied to active carts",
			"active_carts": activeCartCount,
		})
	}

	// Start transaction to delete coupon and related data
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete coupon usage records first
	if err := tx.Where("coupon_id = ?", coupon.ID).Delete(&types.CouponUsage{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete coupon usage records",
		})
	}

	// Delete the coupon
	if err := tx.Delete(&coupon).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete coupon",
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to complete deletion",
		})
	}

	// Security: Log admin action
	logAdminActivity("coupon_delete", adminUser.ID, adminUser.Email, fiber.Map{
		"deleted_coupon_id": coupon.ID,
		"coupon_code":       coupon.Code,
		"used_count":        coupon.UsedCount,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Coupon deleted successfully",
	})
}

func BulkUpdateCoupons(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)

	type BulkCouponUpdate struct {
		CouponIDs []uint `json:"coupon_ids" validate:"required,min=1,max=50"`
		Action    string `json:"action" validate:"required,oneof=activate deactivate make_public make_private delete extend_expiry"`
		
		// Action-specific parameters
		ExpirationExtensionDays *int `json:"expiration_extension_days,omitempty"`
	}

	var req BulkCouponUpdate
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

	// Get all coupons for validation
	var coupons []types.Coupon
	result := db.DB.Where("id IN ?", req.CouponIDs).Find(&coupons)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch coupons",
		})
	}

	if len(coupons) != len(req.CouponIDs) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Some coupons not found",
		})
	}

	var successCount, failureCount int
	var errors []string

	switch req.Action {
	case "activate":
		result := db.DB.Model(&types.Coupon{}).Where("id IN ?", req.CouponIDs).Update("is_active", true)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to activate coupons",
			})
		}
		successCount = int(result.RowsAffected)

	case "deactivate":
		result := db.DB.Model(&types.Coupon{}).Where("id IN ?", req.CouponIDs).Update("is_active", false)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to deactivate coupons",
			})
		}
		successCount = int(result.RowsAffected)

	case "make_public":
		result := db.DB.Model(&types.Coupon{}).Where("id IN ?", req.CouponIDs).Update("is_public", true)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to make coupons public",
			})
		}
		successCount = int(result.RowsAffected)

	case "make_private":
		result := db.DB.Model(&types.Coupon{}).Where("id IN ?", req.CouponIDs).Update("is_public", false)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to make coupons private",
			})
		}
		successCount = int(result.RowsAffected)

	case "extend_expiry":
		if req.ExpirationExtensionDays == nil || *req.ExpirationExtensionDays <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Expiration extension days must be provided and positive",
			})
		}

		for _, coupon := range coupons {
			var newExpiryDate time.Time
			if coupon.ExpiresAt != nil {
				newExpiryDate = coupon.ExpiresAt.AddDate(0, 0, *req.ExpirationExtensionDays)
			} else {
				newExpiryDate = time.Now().AddDate(0, 0, *req.ExpirationExtensionDays)
			}

			if err := db.DB.Model(&coupon).Update("expires_at", newExpiryDate).Error; err != nil {
				failureCount++
				errors = append(errors, fmt.Sprintf("Failed to extend expiry for coupon %s", coupon.Code))
			} else {
				successCount++
			}
		}

	case "delete":
		for _, coupon := range coupons {
			// Check if coupon is in active carts
			var activeCartCount int64
			db.DB.Model(&types.Cart{}).Where("coupon_id = ? AND status = ?", coupon.ID, types.CartStatusActive).Count(&activeCartCount)
			if activeCartCount > 0 {
				failureCount++
				errors = append(errors, fmt.Sprintf("Cannot delete coupon %s - in active carts", coupon.Code))
				continue
			}

			// Delete usage records first
			if err := db.DB.Where("coupon_id = ?", coupon.ID).Delete(&types.CouponUsage{}).Error; err != nil {
				failureCount++
				errors = append(errors, fmt.Sprintf("Failed to delete usage records for coupon %s", coupon.Code))
				continue
			}

			// Delete coupon
			if err := db.DB.Delete(&coupon).Error; err != nil {
				failureCount++
				errors = append(errors, fmt.Sprintf("Failed to delete coupon %s", coupon.Code))
			} else {
				successCount++
			}
		}
	}

	// Security: Log bulk action
	logAdminActivity("coupons_bulk_update", adminUser.ID, adminUser.Email, fiber.Map{
		"action":        req.Action,
		"coupon_ids":    req.CouponIDs,
		"success_count": successCount,
		"failure_count": failureCount,
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

func GetCouponUsageHistory(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	couponID := c.Params("id")

	if couponID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Coupon ID is required",
		})
	}

	// Verify coupon exists
	var coupon types.Coupon
	if result := db.DB.First(&coupon, couponID); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Coupon not found",
		})
	}

	// Query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "25"))
	userID := c.Query("user_id")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 25
	}

	query := db.DB.Model(&types.CouponUsage{}).
		Where("coupon_id = ?", coupon.ID).
		Preload("User")

	// Apply filters
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if dateFrom != "" {
		if parsedDate, err := time.Parse("2006-01-02", dateFrom); err == nil {
			query = query.Where("used_at >= ?", parsedDate)
		}
	}
	if dateTo != "" {
		if parsedDate, err := time.Parse("2006-01-02", dateTo); err == nil {
			query = query.Where("used_at <= ?", parsedDate.Add(24*time.Hour))
		}
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get usage records with pagination
	var usageRecords []types.CouponUsage
	offset := (page - 1) * pageSize
	result := query.Order("used_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&usageRecords)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve usage history",
		})
	}

	// Security: Log admin access
	logAdminActivity("coupon_usage_history", adminUser.ID, adminUser.Email, fiber.Map{
		"target_coupon_id": coupon.ID,
		"coupon_code":      coupon.Code,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"coupon": coupon,
			"usage":  usageRecords,
		},
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}