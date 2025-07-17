package products

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
)

func ListProducts(c *fiber.Ctx) error {
	// Query parameters for pagination and filtering
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "12"))
	sortBy := c.Query("sort", "name")
	sortOrder := c.Query("order", "asc")
	categoryID := c.Query("category")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	featured := c.Query("featured") == "true"
	search := c.Query("search")

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

	query := db.DB.Where("is_active = ?", true)

	// Apply filters
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if minPrice != "" {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice != "" {
		query = query.Where("price <= ?", maxPrice)
	}
	if featured {
		query = query.Where("is_featured = ?", true)
	}
	if search != "" {
		// Security: use parameterized query to prevent SQL injection
		searchTerm := "%" + search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ? OR short_description ILIKE ?", 
			searchTerm, searchTerm, searchTerm)
	}

	// Get total count
	var total int64
	query.Model(&types.Product{}).Count(&total)

	// Get products with pagination
	var products []types.Product
	offset := (page - 1) * pageSize
	orderClause := sortBy + " " + sortOrder
	result := query.Preload("Category").
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
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"sort_by":   sortBy,
			"sort_order": sortOrder,
		},
	})
}

func GetProduct(c *fiber.Ctx) error {
	productSlug := c.Params("slug")
	if productSlug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product slug is required",
		})
	}

	var product types.Product
	result := db.DB.Where("slug = ? AND is_active = ?", productSlug, true).
		Preload("Category").
		First(&product)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Security: Check stock availability without exposing exact stock levels
	stockStatus := "in_stock"
	if product.TrackQuantity {
		if product.StockQuantity <= 0 && !product.AllowBackorder {
			stockStatus = "out_of_stock"
		} else if product.StockQuantity <= 5 {
			stockStatus = "low_stock"
		}
	}

	// Return product with calculated fields
	response := fiber.Map{
		"success": true,
		"data":    product,
		"stock_status": stockStatus,
		"available": product.StockQuantity > 0 || !product.TrackQuantity || product.AllowBackorder,
	}

	return c.JSON(response)
}

func GetFeaturedProducts(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "8"))
	if limit < 1 || limit > 20 {
		limit = 8
	}

	var products []types.Product
	result := db.DB.Where("is_active = ? AND is_featured = ?", true, true).
		Preload("Category").
		Order("updated_at DESC").
		Limit(limit).
		Find(&products)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve featured products",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    products,
	})
}

func SearchProducts(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query is required",
		})
	}

	// Security: limit search query length
	if len(query) > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query too long",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "12"))
	categoryID := c.Query("category")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 12
	}

	dbQuery := db.DB.Where("is_active = ?", true)
	
	// Security: use parameterized query to prevent SQL injection
	searchTerm := "%" + query + "%"
	dbQuery = dbQuery.Where("name ILIKE ? OR description ILIKE ? OR short_description ILIKE ? OR sku ILIKE ?", 
		searchTerm, searchTerm, searchTerm, searchTerm)

	if categoryID != "" {
		dbQuery = dbQuery.Where("category_id = ?", categoryID)
	}

	// Get total count
	var total int64
	dbQuery.Model(&types.Product{}).Count(&total)

	// Get products with pagination
	var products []types.Product
	offset := (page - 1) * pageSize
	result := dbQuery.Preload("Category").
		Order("name ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&products)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search products",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    products,
		"query":   query,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}