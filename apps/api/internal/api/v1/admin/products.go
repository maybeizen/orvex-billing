package admin

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
	"github.com/orvexcc/billing/api/internal/utils"
)

type AdminProductRequest struct {
	CategoryID        uint                     `json:"category_id" validate:"required"`
	Name              string                   `json:"name" validate:"required,min=1,max=255"`
	Slug              string                   `json:"slug" validate:"required,min=1,max=255"`
	Description       string                   `json:"description,omitempty"`
	ShortDesc         string                   `json:"short_description,omitempty" validate:"max=255"`
	Price             float64                  `json:"price" validate:"required,min=0"`
	ComparePrice      *float64                 `json:"compare_price,omitempty" validate:"omitempty,min=0"`
	Currency          string                   `json:"currency" validate:"required,len=3"`
	SKU               string                   `json:"sku,omitempty"`
	StockQuantity     int                      `json:"stock_quantity" validate:"min=0"`
	TrackQuantity     bool                     `json:"track_quantity"`
	AllowBackorder    bool                     `json:"allow_backorder"`
	LowStockThreshold int                      `json:"low_stock_threshold" validate:"min=0"`
	MaxPerUser        *int                     `json:"max_per_user,omitempty" validate:"omitempty,min=1"`
	MinQuantity       int                      `json:"min_quantity" validate:"min=1"`
	MaxQuantity       *int                     `json:"max_quantity,omitempty" validate:"omitempty,min=1"`
	IsActive          bool                     `json:"is_active"`
	IsDigital         bool                     `json:"is_digital"`
	RequiresShipping  bool                     `json:"requires_shipping"`
	IsFeatured        bool                     `json:"is_featured"`
	IsRecurring       bool                     `json:"is_recurring"`
	BillingCycle      *string                  `json:"billing_cycle,omitempty" validate:"omitempty,oneof=one_time monthly yearly"`
	TrialDays         *int                     `json:"trial_days,omitempty" validate:"omitempty,min=0"`
	SetupFee          *float64                 `json:"setup_fee,omitempty" validate:"omitempty,min=0"`
	Resources         types.ProductResources   `json:"resources,omitempty"`
	ConnectedPanels   types.ProductPanels      `json:"connected_panels,omitempty"`
	UpgradeOptions    types.ProductUpgrades    `json:"upgrade_options,omitempty"`
	Images            types.ProductImages      `json:"images,omitempty"`
	MetaTitle         string                   `json:"meta_title,omitempty" validate:"max=255"`
	MetaDesc          string                   `json:"meta_description,omitempty" validate:"max=500"`
	ServiceConfig     *types.ServiceConfig     `json:"service_config,omitempty"`
}

type AdminProductUpdateRequest struct {
	CategoryID        *uint                    `json:"category_id,omitempty"`
	Name              string                   `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Slug              string                   `json:"slug,omitempty" validate:"omitempty,min=1,max=255"`
	Description       string                   `json:"description,omitempty"`
	ShortDesc         string                   `json:"short_description,omitempty" validate:"max=255"`
	Price             *float64                 `json:"price,omitempty" validate:"omitempty,min=0"`
	ComparePrice      *float64                 `json:"compare_price,omitempty" validate:"omitempty,min=0"`
	Currency          string                   `json:"currency,omitempty" validate:"omitempty,len=3"`
	SKU               string                   `json:"sku,omitempty"`
	StockQuantity     *int                     `json:"stock_quantity,omitempty" validate:"omitempty,min=0"`
	TrackQuantity     *bool                    `json:"track_quantity,omitempty"`
	AllowBackorder    *bool                    `json:"allow_backorder,omitempty"`
	LowStockThreshold *int                     `json:"low_stock_threshold,omitempty" validate:"omitempty,min=0"`
	MaxPerUser        *int                     `json:"max_per_user,omitempty" validate:"omitempty,min=1"`
	MinQuantity       *int                     `json:"min_quantity,omitempty" validate:"omitempty,min=1"`
	MaxQuantity       *int                     `json:"max_quantity,omitempty" validate:"omitempty,min=1"`
	IsActive          *bool                    `json:"is_active,omitempty"`
	IsDigital         *bool                    `json:"is_digital,omitempty"`
	RequiresShipping  *bool                    `json:"requires_shipping,omitempty"`
	IsFeatured        *bool                    `json:"is_featured,omitempty"`
	IsRecurring       *bool                    `json:"is_recurring,omitempty"`
	BillingCycle      *string                  `json:"billing_cycle,omitempty" validate:"omitempty,oneof=one_time monthly yearly"`
	TrialDays         *int                     `json:"trial_days,omitempty" validate:"omitempty,min=0"`
	SetupFee          *float64                 `json:"setup_fee,omitempty" validate:"omitempty,min=0"`
	Resources         *types.ProductResources  `json:"resources,omitempty"`
	ConnectedPanels   *types.ProductPanels     `json:"connected_panels,omitempty"`
	UpgradeOptions    *types.ProductUpgrades   `json:"upgrade_options,omitempty"`
	Images            *types.ProductImages     `json:"images,omitempty"`
	MetaTitle         string                   `json:"meta_title,omitempty" validate:"max=255"`
	MetaDesc          string                   `json:"meta_description,omitempty" validate:"max=500"`
	ServiceConfig     *types.ServiceConfig     `json:"service_config,omitempty"`
}

func ListAdminProducts(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)

	// Security: Log admin access
	logAdminActivity("products_list", adminUser.ID, adminUser.Email, fiber.Map{
		"ip": c.IP(),
	})

	// Query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "25"))
	search := c.Query("search")
	categoryID := c.Query("category_id")
	isActive := c.Query("is_active")
	isFeatured := c.Query("is_featured")
	isDigital := c.Query("is_digital")
	isRecurring := c.Query("is_recurring")
	stockStatus := c.Query("stock_status") // in_stock, low_stock, out_of_stock
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
		"created_at": true, "name": true, "price": true, "stock_quantity": true, "is_featured": true,
	}
	if !validSorts[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" {
		sortOrder = "desc"
	}

	query := db.DB.Model(&types.Product{}).Preload("Category")

	// Apply filters
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("products.name ILIKE ? OR products.description ILIKE ? OR products.sku ILIKE ?",
			searchTerm, searchTerm, searchTerm)
	}

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	if isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	if isFeatured != "" {
		query = query.Where("is_featured = ?", isFeatured == "true")
	}

	if isDigital != "" {
		query = query.Where("is_digital = ?", isDigital == "true")
	}

	if isRecurring != "" {
		query = query.Where("is_recurring = ?", isRecurring == "true")
	}

	// Stock status filtering
	switch stockStatus {
	case "in_stock":
		query = query.Where("(track_quantity = false OR stock_quantity > low_stock_threshold)")
	case "low_stock":
		query = query.Where("track_quantity = true AND stock_quantity > 0 AND stock_quantity <= low_stock_threshold")
	case "out_of_stock":
		query = query.Where("track_quantity = true AND stock_quantity = 0")
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get products with pagination
	var products []types.Product
	offset := (page - 1) * pageSize
	orderClause := "products." + sortBy + " " + sortOrder
	result := query.Order(orderClause).
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
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"sort_by":    sortBy,
			"sort_order": sortOrder,
		},
	})
}

func GetAdminProduct(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	productID := c.Params("id")

	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product ID is required",
		})
	}

	var product types.Product
	result := db.DB.Preload("Category").First(&product, productID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Security: Log admin access
	logAdminActivity("product_view", adminUser.ID, adminUser.Email, fiber.Map{
		"target_product_id": product.ID,
		"product_name":      product.Name,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"data":    product,
	})
}

func CreateAdminProduct(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)

	var req AdminProductRequest
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

	// Verify category exists
	var category types.Category
	if result := db.DB.First(&category, req.CategoryID); result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	// Check if slug already exists
	var existingProduct types.Product
	if result := db.DB.Where("slug = ?", req.Slug).First(&existingProduct); result.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Product slug already exists",
		})
	}

	// Check if SKU already exists (if provided)
	if req.SKU != "" {
		if result := db.DB.Where("sku = ?", req.SKU).First(&existingProduct); result.Error == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Product SKU already exists",
			})
		}
	}

	// Create product
	product := types.Product{
		CategoryID:        req.CategoryID,
		Name:              req.Name,
		Slug:              req.Slug,
		Description:       req.Description,
		ShortDesc:         req.ShortDesc,
		Price:             req.Price,
		ComparePrice:      req.ComparePrice,
		Currency:          req.Currency,
		SKU:               req.SKU,
		StockQuantity:     req.StockQuantity,
		TrackQuantity:     req.TrackQuantity,
		AllowBackorder:    req.AllowBackorder,
		LowStockThreshold: req.LowStockThreshold,
		MaxPerUser:        req.MaxPerUser,
		MinQuantity:       req.MinQuantity,
		MaxQuantity:       req.MaxQuantity,
		IsActive:          req.IsActive,
		IsDigital:         req.IsDigital,
		RequiresShipping:  req.RequiresShipping,
		IsFeatured:        req.IsFeatured,
		IsRecurring:       req.IsRecurring,
		BillingCycle:      req.BillingCycle,
		TrialDays:         req.TrialDays,
		SetupFee:          req.SetupFee,
		Resources:         req.Resources,
		ConnectedPanels:   req.ConnectedPanels,
		UpgradeOptions:    req.UpgradeOptions,
		Images:            req.Images,
		MetaTitle:         req.MetaTitle,
		MetaDesc:          req.MetaDesc,
		ServiceConfig:     req.ServiceConfig,
	}

	result := db.DB.Create(&product)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product",
		})
	}

	// Preload category data
	db.DB.Preload("Category").First(&product, product.ID)

	// Security: Log admin action
	logAdminActivity("product_create", adminUser.ID, adminUser.Email, fiber.Map{
		"created_product_id": product.ID,
		"product_name":       product.Name,
		"category_id":        product.CategoryID,
		"price":              product.Price,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Product created successfully",
		"data":    product,
	})
}

func UpdateAdminProduct(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	productID := c.Params("id")

	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product ID is required",
		})
	}

	var req AdminProductUpdateRequest
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

	var product types.Product
	result := db.DB.First(&product, productID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Track changes for logging
	changes := make(map[string]interface{})

	// Update fields if provided
	if req.CategoryID != nil {
		// Verify category exists
		var category types.Category
		if result := db.DB.First(&category, *req.CategoryID); result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Category not found",
			})
		}
		product.CategoryID = *req.CategoryID
		changes["category_id"] = *req.CategoryID
	}

	if req.Name != "" {
		product.Name = req.Name
		changes["name"] = req.Name
	}

	if req.Slug != "" && req.Slug != product.Slug {
		// Check if slug already exists
		var existingProduct types.Product
		if result := db.DB.Where("slug = ? AND id != ?", req.Slug, product.ID).First(&existingProduct); result.Error == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Product slug already exists",
			})
		}
		product.Slug = req.Slug
		changes["slug"] = req.Slug
	}

	if req.Description != "" {
		product.Description = req.Description
		changes["description"] = "updated"
	}
	if req.ShortDesc != "" {
		product.ShortDesc = req.ShortDesc
		changes["short_description"] = "updated"
	}
	if req.Price != nil {
		product.Price = *req.Price
		changes["price"] = *req.Price
	}
	if req.ComparePrice != nil {
		product.ComparePrice = req.ComparePrice
		changes["compare_price"] = *req.ComparePrice
	}
	if req.Currency != "" {
		product.Currency = req.Currency
		changes["currency"] = req.Currency
	}
	if req.SKU != "" && req.SKU != product.SKU {
		// Check if SKU already exists
		var existingProduct types.Product
		if result := db.DB.Where("sku = ? AND id != ?", req.SKU, product.ID).First(&existingProduct); result.Error == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Product SKU already exists",
			})
		}
		product.SKU = req.SKU
		changes["sku"] = req.SKU
	}
	if req.StockQuantity != nil {
		product.StockQuantity = *req.StockQuantity
		changes["stock_quantity"] = *req.StockQuantity
	}
	if req.TrackQuantity != nil {
		product.TrackQuantity = *req.TrackQuantity
		changes["track_quantity"] = *req.TrackQuantity
	}
	if req.AllowBackorder != nil {
		product.AllowBackorder = *req.AllowBackorder
		changes["allow_backorder"] = *req.AllowBackorder
	}
	if req.LowStockThreshold != nil {
		product.LowStockThreshold = *req.LowStockThreshold
		changes["low_stock_threshold"] = *req.LowStockThreshold
	}
	if req.MaxPerUser != nil {
		product.MaxPerUser = req.MaxPerUser
		changes["max_per_user"] = *req.MaxPerUser
	}
	if req.MinQuantity != nil {
		product.MinQuantity = *req.MinQuantity
		changes["min_quantity"] = *req.MinQuantity
	}
	if req.MaxQuantity != nil {
		product.MaxQuantity = req.MaxQuantity
		changes["max_quantity"] = *req.MaxQuantity
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
		changes["is_active"] = *req.IsActive
	}
	if req.IsDigital != nil {
		product.IsDigital = *req.IsDigital
		changes["is_digital"] = *req.IsDigital
	}
	if req.RequiresShipping != nil {
		product.RequiresShipping = *req.RequiresShipping
		changes["requires_shipping"] = *req.RequiresShipping
	}
	if req.IsFeatured != nil {
		product.IsFeatured = *req.IsFeatured
		changes["is_featured"] = *req.IsFeatured
	}
	if req.IsRecurring != nil {
		product.IsRecurring = *req.IsRecurring
		changes["is_recurring"] = *req.IsRecurring
	}
	if req.BillingCycle != nil {
		product.BillingCycle = req.BillingCycle
		changes["billing_cycle"] = *req.BillingCycle
	}
	if req.TrialDays != nil {
		product.TrialDays = req.TrialDays
		changes["trial_days"] = *req.TrialDays
	}
	if req.SetupFee != nil {
		product.SetupFee = req.SetupFee
		changes["setup_fee"] = *req.SetupFee
	}
	if req.Resources != nil {
		product.Resources = *req.Resources
		changes["resources"] = "updated"
	}
	if req.ConnectedPanels != nil {
		product.ConnectedPanels = *req.ConnectedPanels
		changes["connected_panels"] = "updated"
	}
	if req.UpgradeOptions != nil {
		product.UpgradeOptions = *req.UpgradeOptions
		changes["upgrade_options"] = "updated"
	}
	if req.Images != nil {
		product.Images = *req.Images
		changes["images"] = "updated"
	}
	if req.MetaTitle != "" {
		product.MetaTitle = req.MetaTitle
		changes["meta_title"] = req.MetaTitle
	}
	if req.MetaDesc != "" {
		product.MetaDesc = req.MetaDesc
		changes["meta_description"] = "updated"
	}
	if req.ServiceConfig != nil {
		product.ServiceConfig = req.ServiceConfig
		changes["service_config"] = "updated"
	}

	if err := db.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update product",
		})
	}

	// Preload category data
	db.DB.Preload("Category").First(&product, product.ID)

	// Security: Log admin action
	logAdminActivity("product_update", adminUser.ID, adminUser.Email, fiber.Map{
		"target_product_id": product.ID,
		"product_name":      product.Name,
		"changes":           changes,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Product updated successfully",
		"data":    product,
	})
}

func DeleteAdminProduct(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	productID := c.Params("id")

	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product ID is required",
		})
	}

	var product types.Product
	result := db.DB.First(&product, productID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Check if product is in any active carts
	var cartItemCount int64
	db.DB.Model(&types.CartItem{}).Where("product_id = ?", product.ID).Count(&cartItemCount)
	if cartItemCount > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete product that is in active carts",
			"cart_items": cartItemCount,
		})
	}

	// Check if product has active services
	var serviceCount int64
	db.DB.Model(&types.Service{}).Where("id IN (SELECT service_id FROM invoices WHERE id IN (SELECT invoice_id FROM invoice_items WHERE description ILIKE ?))", "%"+product.Name+"%").Count(&serviceCount)

	// Soft delete product
	if err := db.DB.Delete(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete product",
		})
	}

	// Security: Log admin action
	logAdminActivity("product_delete", adminUser.ID, adminUser.Email, fiber.Map{
		"deleted_product_id": product.ID,
		"product_name":       product.Name,
		"category_id":        product.CategoryID,
		"had_services":       serviceCount > 0,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Product deleted successfully",
	})
}

func DuplicateAdminProduct(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	productID := c.Params("id")

	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product ID is required",
		})
	}

	var originalProduct types.Product
	result := db.DB.First(&originalProduct, productID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Create duplicate
	duplicate := originalProduct
	duplicate.ID = 0 // Reset ID for new record
	duplicate.Name = originalProduct.Name + " (Copy)"
	duplicate.Slug = originalProduct.Slug + "-copy"
	duplicate.SKU = "" // Clear SKU to avoid conflicts
	duplicate.IsActive = false // Set as inactive by default

	// Ensure unique slug
	slugCounter := 1
	baseSlug := duplicate.Slug
	for {
		var existingProduct types.Product
		if result := db.DB.Where("slug = ?", duplicate.Slug).First(&existingProduct); result.Error != nil {
			break // Slug is unique
		}
		slugCounter++
		duplicate.Slug = baseSlug + "-" + strconv.Itoa(slugCounter)
	}

	if err := db.DB.Create(&duplicate).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to duplicate product",
		})
	}

	// Preload category data
	db.DB.Preload("Category").First(&duplicate, duplicate.ID)

	// Security: Log admin action
	logAdminActivity("product_duplicate", adminUser.ID, adminUser.Email, fiber.Map{
		"original_product_id":   originalProduct.ID,
		"duplicated_product_id": duplicate.ID,
		"product_name":          duplicate.Name,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Product duplicated successfully",
		"data":    duplicate,
	})
}

func BulkUpdateProducts(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)

	type BulkProductUpdate struct {
		ProductIDs []uint `json:"product_ids" validate:"required,min=1,max=100"`
		Action     string `json:"action" validate:"required,oneof=activate deactivate feature unfeature delete set_category adjust_price"`
		
		// Action-specific parameters
		CategoryID   *uint    `json:"category_id,omitempty"`
		PriceChange  *float64 `json:"price_change,omitempty"` // Can be positive or negative
		PricePercent *float64 `json:"price_percent,omitempty"` // Percentage change
	}

	var req BulkProductUpdate
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

	// Get all products for validation
	var products []types.Product
	result := db.DB.Where("id IN ?", req.ProductIDs).Find(&products)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch products",
		})
	}

	if len(products) != len(req.ProductIDs) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Some products not found",
		})
	}

	var successCount, failureCount int
	var errors []string
	var updateData map[string]interface{}

	// Prepare update data based on action
	switch req.Action {
	case "activate":
		updateData = map[string]interface{}{"is_active": true}
	case "deactivate":
		updateData = map[string]interface{}{"is_active": false}
	case "feature":
		updateData = map[string]interface{}{"is_featured": true}
	case "unfeature":
		updateData = map[string]interface{}{"is_featured": false}
	case "set_category":
		if req.CategoryID == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Category ID is required for set_category action",
			})
		}
		// Verify category exists
		var category types.Category
		if result := db.DB.First(&category, *req.CategoryID); result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Category not found",
			})
		}
		updateData = map[string]interface{}{"category_id": *req.CategoryID}
	case "delete":
		// Handle deletion separately
		for _, product := range products {
			// Check if product is in any active carts
			var cartItemCount int64
			db.DB.Model(&types.CartItem{}).Where("product_id = ?", product.ID).Count(&cartItemCount)
			if cartItemCount > 0 {
				failureCount++
				errors = append(errors, fmt.Sprintf("Cannot delete product %s (ID: %d) - in active carts", product.Name, product.ID))
				continue
			}

			if err := db.DB.Delete(&product).Error; err != nil {
				failureCount++
				errors = append(errors, fmt.Sprintf("Failed to delete product %s (ID: %d)", product.Name, product.ID))
			} else {
				successCount++
			}
		}
	case "adjust_price":
		// Handle price adjustment separately
		for _, product := range products {
			newPrice := product.Price
			
			if req.PriceChange != nil {
				newPrice += *req.PriceChange
			} else if req.PricePercent != nil {
				newPrice *= (1 + (*req.PricePercent / 100))
			} else {
				failureCount++
				errors = append(errors, fmt.Sprintf("No price change specified for product %s", product.Name))
				continue
			}

			if newPrice < 0 {
				failureCount++
				errors = append(errors, fmt.Sprintf("Price cannot be negative for product %s", product.Name))
				continue
			}

			if err := db.DB.Model(&product).Update("price", newPrice).Error; err != nil {
				failureCount++
				errors = append(errors, fmt.Sprintf("Failed to update price for product %s", product.Name))
			} else {
				successCount++
			}
		}
	}

	// Apply bulk update for non-special actions
	if req.Action != "delete" && req.Action != "adjust_price" && updateData != nil {
		result := db.DB.Model(&types.Product{}).Where("id IN ?", req.ProductIDs).Updates(updateData)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to perform bulk update",
			})
		}
		successCount = int(result.RowsAffected)
		failureCount = len(req.ProductIDs) - successCount
	}

	// Security: Log bulk action
	logAdminActivity("products_bulk_update", adminUser.ID, adminUser.Email, fiber.Map{
		"action":        req.Action,
		"product_ids":   req.ProductIDs,
		"success_count": successCount,
		"failure_count": failureCount,
		"category_id":   req.CategoryID,
		"price_change":  req.PriceChange,
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