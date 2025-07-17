package admin

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
	"github.com/orvexcc/billing/api/internal/utils"
)

type AdminCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Slug        string `json:"slug" validate:"required,min=1,max=255"`
	Description string `json:"description,omitempty"`
	ImageURL    string `json:"image_url,omitempty" validate:"omitempty,url"`
	IsActive    bool   `json:"is_active"`
	SortOrder   int    `json:"sort_order" validate:"min=0"`
}

type AdminCategoryUpdateRequest struct {
	Name        string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Slug        string `json:"slug,omitempty" validate:"omitempty,min=1,max=255"`
	Description string `json:"description,omitempty"`
	ImageURL    string `json:"image_url,omitempty" validate:"omitempty,url"`
	IsActive    *bool  `json:"is_active,omitempty"`
	SortOrder   *int   `json:"sort_order,omitempty" validate:"omitempty,min=0"`
}

func ListAdminCategories(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)

	// Security: Log admin access
	logAdminActivity("categories_list", adminUser.ID, adminUser.Email, fiber.Map{
		"ip": c.IP(),
	})

	// Query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "25"))
	search := c.Query("search")
	isActive := c.Query("is_active")
	sortBy := c.Query("sort", "sort_order")
	sortOrder := c.Query("order", "asc")
	includeProducts := c.Query("include_products") == "true"

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 25
	}

	// Security: Validate sort parameters
	validSorts := map[string]bool{
		"created_at": true, "name": true, "sort_order": true, "id": true,
	}
	if !validSorts[sortBy] {
		sortBy = "sort_order"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	query := db.DB.Model(&types.Category{})

	// Include products count or full products
	if includeProducts {
		query = query.Preload("Products")
	} else {
		// Add product count as a virtual field
		query = query.Select("categories.*, (SELECT COUNT(*) FROM products WHERE products.category_id = categories.id AND products.deleted_at IS NULL) as product_count")
	}

	// Apply filters
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("categories.name ILIKE ? OR categories.description ILIKE ?", searchTerm, searchTerm)
	}

	if isActive != "" {
		query = query.Where("categories.is_active = ?", isActive == "true")
	}

	// Get total count
	var total int64
	db.DB.Model(&types.Category{}).Count(&total)

	// Get categories with pagination
	var categories []types.Category
	offset := (page - 1) * pageSize
	orderClause := "categories." + sortBy + " " + sortOrder
	result := query.Order(orderClause).
		Offset(offset).
		Limit(pageSize).
		Find(&categories)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve categories",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    categories,
		"meta": fiber.Map{
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"sort_by":    sortBy,
			"sort_order": sortOrder,
		},
	})
}

func GetAdminCategory(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	categoryID := c.Params("id")

	if categoryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category ID is required",
		})
	}

	var category types.Category
	result := db.DB.Preload("Products").First(&category, categoryID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	// Get category statistics
	var stats struct {
		TotalProducts    int64   `json:"total_products"`
		ActiveProducts   int64   `json:"active_products"`
		FeaturedProducts int64   `json:"featured_products"`
		AvgPrice         float64 `json:"average_price"`
		TotalRevenue     float64 `json:"total_revenue"`
	}

	db.DB.Model(&types.Product{}).Where("category_id = ?", category.ID).Count(&stats.TotalProducts)
	db.DB.Model(&types.Product{}).Where("category_id = ? AND is_active = ?", category.ID, true).Count(&stats.ActiveProducts)
	db.DB.Model(&types.Product{}).Where("category_id = ? AND is_featured = ?", category.ID, true).Count(&stats.FeaturedProducts)
	
	// Calculate average price
	var avgResult struct {
		AvgPrice float64
	}
	db.DB.Model(&types.Product{}).
		Select("COALESCE(AVG(price), 0) as avg_price").
		Where("category_id = ? AND is_active = ?", category.ID, true).
		Scan(&avgResult)
	stats.AvgPrice = avgResult.AvgPrice

	// Calculate total revenue from this category (approximate)
	var revenueResult struct {
		TotalRevenue float64
	}
	db.DB.Model(&types.Invoice{}).
		Select("COALESCE(SUM(invoices.total), 0) as total_revenue").
		Joins("JOIN invoice_items ON invoices.id = invoice_items.invoice_id").
		Joins("JOIN products ON invoice_items.description ILIKE CONCAT('%', products.name, '%')").
		Where("products.category_id = ? AND invoices.status = ?", category.ID, types.InvoiceStatusPaid).
		Scan(&revenueResult)
	stats.TotalRevenue = revenueResult.TotalRevenue

	// Security: Log admin access
	logAdminActivity("category_view", adminUser.ID, adminUser.Email, fiber.Map{
		"target_category_id": category.ID,
		"category_name":      category.Name,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"category": category,
			"stats":    stats,
		},
	})
}

func CreateAdminCategory(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)

	var req AdminCategoryRequest
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

	// Check if name already exists
	var existingCategory types.Category
	if result := db.DB.Where("name = ?", req.Name).First(&existingCategory); result.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Category name already exists",
		})
	}

	// Check if slug already exists
	if result := db.DB.Where("slug = ?", req.Slug).First(&existingCategory); result.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Category slug already exists",
		})
	}

	// Create category
	category := types.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		IsActive:    req.IsActive,
		SortOrder:   req.SortOrder,
	}

	result := db.DB.Create(&category)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create category",
		})
	}

	// Security: Log admin action
	logAdminActivity("category_create", adminUser.ID, adminUser.Email, fiber.Map{
		"created_category_id": category.ID,
		"category_name":       category.Name,
		"sort_order":          category.SortOrder,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Category created successfully",
		"data":    category,
	})
}

func UpdateAdminCategory(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	categoryID := c.Params("id")

	if categoryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category ID is required",
		})
	}

	var req AdminCategoryUpdateRequest
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

	var category types.Category
	result := db.DB.First(&category, categoryID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	// Track changes for logging
	changes := make(map[string]interface{})

	// Update fields if provided
	if req.Name != "" && req.Name != category.Name {
		// Check if name already exists
		var existingCategory types.Category
		if result := db.DB.Where("name = ? AND id != ?", req.Name, category.ID).First(&existingCategory); result.Error == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Category name already exists",
			})
		}
		category.Name = req.Name
		changes["name"] = req.Name
	}

	if req.Slug != "" && req.Slug != category.Slug {
		// Check if slug already exists
		var existingCategory types.Category
		if result := db.DB.Where("slug = ? AND id != ?", req.Slug, category.ID).First(&existingCategory); result.Error == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Category slug already exists",
			})
		}
		category.Slug = req.Slug
		changes["slug"] = req.Slug
	}

	if req.Description != "" {
		category.Description = req.Description
		changes["description"] = "updated"
	}
	if req.ImageURL != "" {
		category.ImageURL = req.ImageURL
		changes["image_url"] = req.ImageURL
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
		changes["is_active"] = *req.IsActive
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
		changes["sort_order"] = *req.SortOrder
	}

	if err := db.DB.Save(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update category",
		})
	}

	// Security: Log admin action
	logAdminActivity("category_update", adminUser.ID, adminUser.Email, fiber.Map{
		"target_category_id": category.ID,
		"category_name":      category.Name,
		"changes":            changes,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Category updated successfully",
		"data":    category,
	})
}

func DeleteAdminCategory(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	categoryID := c.Params("id")

	if categoryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category ID is required",
		})
	}

	var category types.Category
	result := db.DB.First(&category, categoryID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	// Check if category has products - DEPENDENCY VALIDATION
	var productCount int64
	db.DB.Model(&types.Product{}).Where("category_id = ?", category.ID).Count(&productCount)
	if productCount > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete category with existing products",
			"product_count": productCount,
			"message": "Please move or delete all products in this category first",
		})
	}

	// Soft delete category
	if err := db.DB.Delete(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete category",
		})
	}

	// Security: Log admin action
	logAdminActivity("category_delete", adminUser.ID, adminUser.Email, fiber.Map{
		"deleted_category_id": category.ID,
		"category_name":       category.Name,
		"product_count":       productCount,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Category deleted successfully",
	})
}

func ReorderCategories(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)

	type CategoryOrder struct {
		ID        uint `json:"id" validate:"required"`
		SortOrder int  `json:"sort_order" validate:"min=0"`
	}

	type ReorderRequest struct {
		Categories []CategoryOrder `json:"categories" validate:"required,min=1,max=100"`
	}

	var req ReorderRequest
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

	// Start transaction
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var successCount, failureCount int
	var errors []string

	for _, categoryOrder := range req.Categories {
		var category types.Category
		if result := tx.First(&category, categoryOrder.ID); result.Error != nil {
			failureCount++
			errors = append(errors, fmt.Sprintf("Category with ID %d not found", categoryOrder.ID))
			continue
		}

		category.SortOrder = categoryOrder.SortOrder
		if err := tx.Save(&category).Error; err != nil {
			failureCount++
			errors = append(errors, fmt.Sprintf("Failed to update category %s", category.Name))
		} else {
			successCount++
		}
	}

	if failureCount > 0 {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to reorder categories",
			"details": errors,
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save category order",
		})
	}

	// Security: Log admin action
	logAdminActivity("categories_reorder", adminUser.ID, adminUser.Email, fiber.Map{
		"categories_count": len(req.Categories),
		"success_count":    successCount,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Categories reordered successfully",
		"data": fiber.Map{
			"updated_count": successCount,
		},
	})
}

func BulkUpdateCategories(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)

	type BulkCategoryUpdate struct {
		CategoryIDs []uint `json:"category_ids" validate:"required,min=1,max=50"`
		Action      string `json:"action" validate:"required,oneof=activate deactivate delete"`
	}

	var req BulkCategoryUpdate
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

	// Get all categories for validation
	var categories []types.Category
	result := db.DB.Where("id IN ?", req.CategoryIDs).Find(&categories)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch categories",
		})
	}

	if len(categories) != len(req.CategoryIDs) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Some categories not found",
		})
	}

	var successCount, failureCount int
	var errors []string

	switch req.Action {
	case "activate":
		result := db.DB.Model(&types.Category{}).Where("id IN ?", req.CategoryIDs).Update("is_active", true)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to activate categories",
			})
		}
		successCount = int(result.RowsAffected)

	case "deactivate":
		result := db.DB.Model(&types.Category{}).Where("id IN ?", req.CategoryIDs).Update("is_active", false)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to deactivate categories",
			})
		}
		successCount = int(result.RowsAffected)

	case "delete":
		// Check each category for products before deletion
		for _, category := range categories {
			var productCount int64
			db.DB.Model(&types.Product{}).Where("category_id = ?", category.ID).Count(&productCount)
			if productCount > 0 {
				failureCount++
				errors = append(errors, fmt.Sprintf("Cannot delete category %s - has %d products", category.Name, productCount))
				continue
			}

			if err := db.DB.Delete(&category).Error; err != nil {
				failureCount++
				errors = append(errors, fmt.Sprintf("Failed to delete category %s", category.Name))
			} else {
				successCount++
			}
		}
	}

	// Security: Log bulk action
	logAdminActivity("categories_bulk_update", adminUser.ID, adminUser.Email, fiber.Map{
		"action":        req.Action,
		"category_ids":  req.CategoryIDs,
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

func GetCategoryProductsAdmin(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	categoryID := c.Params("id")

	if categoryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category ID is required",
		})
	}

	// Verify category exists
	var category types.Category
	if result := db.DB.First(&category, categoryID); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	// Query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "25"))
	isActive := c.Query("is_active")
	isFeatured := c.Query("is_featured")
	sortBy := c.Query("sort", "created_at")
	sortOrder := c.Query("order", "desc")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 25
	}

	query := db.DB.Model(&types.Product{}).Where("category_id = ?", categoryID)

	// Apply filters
	if isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}
	if isFeatured != "" {
		query = query.Where("is_featured = ?", isFeatured == "true")
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get products with pagination
	var products []types.Product
	offset := (page - 1) * pageSize
	orderClause := sortBy + " " + sortOrder
	result := query.Order(orderClause).
		Offset(offset).
		Limit(pageSize).
		Find(&products)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve products",
		})
	}

	// Security: Log admin access
	logAdminActivity("category_products_view", adminUser.ID, adminUser.Email, fiber.Map{
		"target_category_id": category.ID,
		"category_name":      category.Name,
		"products_count":     total,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"category": category,
			"products": products,
		},
		"meta": fiber.Map{
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"sort_by":    sortBy,
			"sort_order": sortOrder,
		},
	})
}