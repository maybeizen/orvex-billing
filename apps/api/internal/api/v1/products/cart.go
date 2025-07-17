package products

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
	"github.com/orvexcc/billing/api/internal/utils"
)

// Security function to generate price hash
func generatePriceHash(productID uint, price float64) string {
	data := fmt.Sprintf("%d:%f", productID, price)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// Security function to validate price hasn't been tampered with
func validatePriceHash(productID uint, price float64, hash string) bool {
	expectedHash := generatePriceHash(productID, price)
	return expectedHash == hash
}

func GetCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	var cart types.Cart
	result := db.DB.Where("user_id = ? AND status = ?", userID, types.CartStatusActive).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		Preload("Coupon").
		First(&cart)

	if result.Error != nil {
		// Create new cart if none exists
		cart = types.Cart{
			UserID:   userID.(uint),
			Status:   types.CartStatusActive,
			Currency: "USD",
		}
		if err := db.DB.Create(&cart).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create cart",
			})
		}
	}

	// Calculate and update cart totals
	updateCartTotals(&cart)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    cart,
	})
}

func AddToCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	var req types.AddToCartRequest
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

	// Security: Verify product exists and is active
	var product types.Product
	result := db.DB.Where("id = ? AND is_active = ?", req.ProductID, true).First(&product)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found or inactive",
		})
	}

	// Security: Check stock availability
	if product.TrackQuantity && product.StockQuantity < req.Quantity && !product.AllowBackorder {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Insufficient stock available",
			"available_quantity": product.StockQuantity,
		})
	}

	// Get or create user's active cart
	var cart types.Cart
	result = db.DB.Where("user_id = ? AND status = ?", userID, types.CartStatusActive).First(&cart)
	if result.Error != nil {
		cart = types.Cart{
			UserID:   userID.(uint),
			Status:   types.CartStatusActive,
			Currency: "USD",
		}
		if err := db.DB.Create(&cart).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create cart",
			})
		}
	}

	// Check if item already exists in cart
	var existingItem types.CartItem
	itemExists := db.DB.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&existingItem).Error == nil

	if itemExists {
		// Security: Validate existing price hash
		if !validatePriceHash(existingItem.ProductID, existingItem.UnitPrice, existingItem.PriceHash) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Price validation failed - cart may have been tampered with",
			})
		}

		// Update quantity
		newQuantity := existingItem.Quantity + req.Quantity
		
		// Security: Check total quantity against stock
		if product.TrackQuantity && product.StockQuantity < newQuantity && !product.AllowBackorder {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Insufficient stock for requested quantity",
				"available_quantity": product.StockQuantity,
				"current_in_cart": existingItem.Quantity,
			})
		}

		existingItem.Quantity = newQuantity
		existingItem.TotalPrice = existingItem.UnitPrice * float64(newQuantity)
		
		if err := db.DB.Save(&existingItem).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update cart item",
			})
		}
	} else {
		// Add new item to cart
		priceHash := generatePriceHash(product.ID, product.Price)
		
		cartItem := types.CartItem{
			CartID:     cart.ID,
			ProductID:  product.ID,
			Quantity:   req.Quantity,
			UnitPrice:  product.Price,
			TotalPrice: product.Price * float64(req.Quantity),
			PriceHash:  priceHash,
		}

		if err := db.DB.Create(&cartItem).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to add item to cart",
			})
		}
	}

	// Update cart totals
	updateCartTotals(&cart)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Product added to cart successfully",
	})
}

func UpdateCartItem(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	itemID := c.Params("itemId")
	if itemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Item ID is required",
		})
	}

	var req types.UpdateCartItemRequest
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

	// Security: Verify item belongs to user's cart
	var cartItem types.CartItem
	result := db.DB.Joins("JOIN carts ON cart_items.cart_id = carts.id").
		Where("cart_items.id = ? AND carts.user_id = ? AND carts.status = ?", 
			itemID, userID, types.CartStatusActive).
		Preload("Product").
		First(&cartItem)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Cart item not found",
		})
	}

	// Security: Validate price hash
	if !validatePriceHash(cartItem.ProductID, cartItem.UnitPrice, cartItem.PriceHash) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Price validation failed - item may have been tampered with",
		})
	}

	if req.Quantity == 0 {
		// Remove item from cart
		if err := db.DB.Delete(&cartItem).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to remove item from cart",
			})
		}
	} else {
		// Security: Check stock availability
		if cartItem.Product.TrackQuantity && 
		   cartItem.Product.StockQuantity < req.Quantity && 
		   !cartItem.Product.AllowBackorder {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Insufficient stock available",
				"available_quantity": cartItem.Product.StockQuantity,
			})
		}

		// Update quantity
		cartItem.Quantity = req.Quantity
		cartItem.TotalPrice = cartItem.UnitPrice * float64(req.Quantity)
		
		if err := db.DB.Save(&cartItem).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update cart item",
			})
		}
	}

	// Update cart totals
	var cart types.Cart
	db.DB.First(&cart, cartItem.CartID)
	updateCartTotals(&cart)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Cart item updated successfully",
	})
}

func RemoveFromCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	itemID := c.Params("itemId")
	if itemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Item ID is required",
		})
	}

	// Security: Verify item belongs to user's cart
	var cartItem types.CartItem
	result := db.DB.Joins("JOIN carts ON cart_items.cart_id = carts.id").
		Where("cart_items.id = ? AND carts.user_id = ? AND carts.status = ?", 
			itemID, userID, types.CartStatusActive).
		First(&cartItem)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Cart item not found",
		})
	}

	cartID := cartItem.CartID
	
	if err := db.DB.Delete(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove item from cart",
		})
	}

	// Update cart totals
	var cart types.Cart
	db.DB.First(&cart, cartID)
	updateCartTotals(&cart)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Item removed from cart successfully",
	})
}

func ClearCart(c *fiber.Ctx) error {
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

	// Remove all items from cart
	if err := db.DB.Where("cart_id = ?", cart.ID).Delete(&types.CartItem{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to clear cart",
		})
	}

	// Reset cart totals and remove coupon
	cart.SubTotal = 0
	cart.TaxAmount = 0
	cart.DiscountAmount = 0
	cart.Total = 0
	cart.CouponID = nil
	cart.CouponCode = ""
	
	if err := db.DB.Save(&cart).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update cart",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Cart cleared successfully",
	})
}

// Helper function to calculate and update cart totals
func updateCartTotals(cart *types.Cart) {
	var items []types.CartItem
	db.DB.Where("cart_id = ?", cart.ID).Find(&items)

	var subtotal float64
	for _, item := range items {
		subtotal += item.TotalPrice
	}

	cart.SubTotal = subtotal
	
	// Calculate tax (simplified - 8% tax rate)
	taxRate := 0.08
	cart.TaxAmount = subtotal * taxRate

	// Apply coupon discount if applicable
	cart.DiscountAmount = 0
	if cart.CouponID != nil {
		var coupon types.Coupon
		if db.DB.First(&coupon, *cart.CouponID).Error == nil {
			cart.DiscountAmount = calculateCouponDiscount(&coupon, subtotal)
		}
	}

	cart.Total = subtotal + cart.TaxAmount - cart.DiscountAmount
	if cart.Total < 0 {
		cart.Total = 0
	}

	db.DB.Save(cart)
}

// Helper function to calculate coupon discount
func calculateCouponDiscount(coupon *types.Coupon, subtotal float64) float64 {
	if !coupon.IsActive {
		return 0
	}

	// Check date restrictions
	now := time.Now()
	if coupon.StartsAt != nil && now.Before(*coupon.StartsAt) {
		return 0
	}
	if coupon.ExpiresAt != nil && now.After(*coupon.ExpiresAt) {
		return 0
	}

	// Check minimum order amount
	if coupon.MinOrderAmount != nil && subtotal < *coupon.MinOrderAmount {
		return 0
	}

	var discount float64
	switch coupon.Type {
	case types.CouponTypePercentage:
		discount = subtotal * (coupon.Value / 100)
	case types.CouponTypeFixed:
		discount = coupon.Value
	case types.CouponTypeFreeShip:
		// For free shipping, return a fixed shipping cost discount
		discount = 10.00 // Simplified shipping cost
	}

	// Apply maximum discount limit
	if coupon.MaxDiscountAmount != nil && discount > *coupon.MaxDiscountAmount {
		discount = *coupon.MaxDiscountAmount
	}

	return discount
}