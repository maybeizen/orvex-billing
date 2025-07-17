package products

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
)

func ListCategories(c *fiber.Ctx) error {
	// Query parameters for pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	showInactive := c.Query("show_inactive") == "true"

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := db.DB.Model(&types.Category{})
	
	// Filter active categories by default (security: don't expose inactive)
	if !showInactive {
		query = query.Where("is_active = ?", true)
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get categories with pagination
	var categories []types.Category
	offset := (page - 1) * pageSize
	result := query.Order("sort_order ASC, name ASC").
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
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetCategory(c *fiber.Ctx) error {
	categorySlug := c.Params("slug")
	if categorySlug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category slug is required",
		})
	}

	var category types.Category
	result := db.DB.Where("slug = ? AND is_active = ?", categorySlug, true).
		Preload("Products", "is_active = ?", true).
		First(&category)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    category,
	})
}

func GetCategoryProducts(c *fiber.Ctx) error {
	categorySlug := c.Params("slug")
	if categorySlug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category slug is required",
		})
	}

	// Verify category exists and is active
	var category types.Category
	result := db.DB.Where("slug = ? AND is_active = ?", categorySlug, true).First(&category)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	// Query parameters for pagination and filtering
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "12"))
	sortBy := c.Query("sort", "name") // name, price, created_at
	sortOrder := c.Query("order", "asc") // asc, desc
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	featured := c.Query("featured") == "true"

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 12
	}

	// Security: validate sort parameters
	validSorts := map[string]bool{
		"name": true, "price": true, "created_at": true, "updated_at": true,
	}
	if !validSorts[sortBy] {
		sortBy = "name"
	}
	if sortOrder != "desc" {
		sortOrder = "asc"
	}

	query := db.DB.Where("category_id = ? AND is_active = ?", category.ID, true)

	// Apply filters
	if minPrice != "" {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice != "" {
		query = query.Where("price <= ?", maxPrice)
	}
	if featured {
		query = query.Where("is_featured = ?", true)
	}

	// Get total count
	var total int64
	query.Model(&types.Product{}).Count(&total)

	// Get products with pagination
	var products []types.Product
	offset := (page - 1) * pageSize
	orderClause := sortBy + " " + sortOrder
	result = query.Preload("Category").
		Order(orderClause).
		Offset(offset).
		Limit(pageSize).
		Find(&products)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve products",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    products,
		"category": category,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"sort_by":   sortBy,
			"sort_order": sortOrder,
		},
	})
}