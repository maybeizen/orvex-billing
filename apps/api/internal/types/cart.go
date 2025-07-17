package types

import (
	"time"
	"gorm.io/gorm"
)

type Cart struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id" gorm:"not null;index"`
	User      User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	SessionID string     `json:"session_id" gorm:"index"` // For guest carts
	
	// Cart totals (calculated fields)
	SubTotal      float64 `json:"subtotal" gorm:"type:decimal(10,2);default:0"`
	TaxAmount     float64 `json:"tax_amount" gorm:"type:decimal(10,2);default:0"`
	DiscountAmount float64 `json:"discount_amount" gorm:"type:decimal(10,2);default:0"`
	Total         float64 `json:"total" gorm:"type:decimal(10,2);default:0"`
	Currency      string  `json:"currency" gorm:"default:'USD'"`
	
	// Applied coupon
	CouponID    *uint   `json:"coupon_id,omitempty"`
	Coupon      *Coupon `json:"coupon,omitempty" gorm:"foreignKey:CouponID"`
	CouponCode  string  `json:"coupon_code"`
	
	// Cart status and expiration
	Status      CartStatus `json:"status" gorm:"default:'active'"`
	ExpiresAt   *time.Time `json:"expires_at"`
	
	// Relationships
	Items       []CartItem `json:"items,omitempty" gorm:"foreignKey:CartID"`
	
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type CartStatus string

const (
	CartStatusActive    CartStatus = "active"
	CartStatusAbandoned CartStatus = "abandoned"
	CartStatusCompleted CartStatus = "completed"
	CartStatusExpired   CartStatus = "expired"
)

type CartItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	CartID    uint    `json:"cart_id" gorm:"not null;index"`
	Cart      Cart    `json:"cart,omitempty" gorm:"foreignKey:CartID"`
	ProductID uint    `json:"product_id" gorm:"not null;index"`
	Product   Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	
	// Item details (snapshot at time of adding to cart)
	Quantity     int     `json:"quantity" gorm:"not null;default:1"`
	UnitPrice    float64 `json:"unit_price" gorm:"type:decimal(10,2);not null"`
	TotalPrice   float64 `json:"total_price" gorm:"type:decimal(10,2);not null"`
	
	// Security: prevent price manipulation
	PriceHash    string  `json:"-" gorm:"not null"` // Hash of product price at time of adding
	
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Coupon struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	Code        string      `json:"code" gorm:"not null;uniqueIndex"`
	Name        string      `json:"name" gorm:"not null"`
	Description string      `json:"description" gorm:"type:text"`
	
	// Discount configuration
	Type        CouponType  `json:"type" gorm:"not null"`
	Value       float64     `json:"value" gorm:"type:decimal(10,2);not null"`
	
	// Usage limits and validation
	MinOrderAmount *float64   `json:"min_order_amount,omitempty" gorm:"type:decimal(10,2)"`
	MaxDiscountAmount *float64 `json:"max_discount_amount,omitempty" gorm:"type:decimal(10,2)"`
	UsageLimit     *int       `json:"usage_limit,omitempty"`
	UsedCount      int        `json:"used_count" gorm:"default:0"`
	UserUsageLimit *int       `json:"user_usage_limit,omitempty"`
	
	// Date restrictions
	StartsAt    *time.Time  `json:"starts_at,omitempty"`
	ExpiresAt   *time.Time  `json:"expires_at,omitempty"`
	
	// Status and visibility
	IsActive    bool        `json:"is_active" gorm:"default:true"`
	IsPublic    bool        `json:"is_public" gorm:"default:false"`
	
	// Product/Category restrictions
	ApplicableProductIDs  []uint `json:"applicable_product_ids,omitempty" gorm:"type:text"`
	ApplicableCategoryIDs []uint `json:"applicable_category_ids,omitempty" gorm:"type:text"`
	ExcludedProductIDs    []uint `json:"excluded_product_ids,omitempty" gorm:"type:text"`
	
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type CouponType string

const (
	CouponTypePercentage CouponType = "percentage"
	CouponTypeFixed      CouponType = "fixed"
	CouponTypeFreeShip   CouponType = "free_shipping"
)

type CouponUsage struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	CouponID uint      `json:"coupon_id" gorm:"not null;index"`
	Coupon   Coupon    `json:"coupon,omitempty" gorm:"foreignKey:CouponID"`
	UserID   uint      `json:"user_id" gorm:"not null;index"`
	User     User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CartID   *uint     `json:"cart_id,omitempty"`
	InvoiceID *uint    `json:"invoice_id,omitempty"`
	
	DiscountAmount float64   `json:"discount_amount" gorm:"type:decimal(10,2);not null"`
	UsedAt         time.Time `json:"used_at" gorm:"not null"`
	
	CreatedAt      time.Time `json:"created_at"`
}

// Request/Response types
type AddToCartRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,min=1,max=100"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" validate:"required,min=0,max=100"`
}

type ApplyCouponRequest struct {
	CouponCode string `json:"coupon_code" validate:"required,min=3,max=50"`
}

type CheckoutRequest struct {
	PaymentMethod string `json:"payment_method" validate:"required,oneof=stripe paypal"`
	BillingAddress CheckoutAddress `json:"billing_address,omitempty"`
	ShippingAddress CheckoutAddress `json:"shipping_address,omitempty"`
	Notes         string `json:"notes,omitempty" validate:"max=500"`
}

type CheckoutAddress struct {
	FirstName string `json:"first_name" validate:"required,max=50"`
	LastName  string `json:"last_name" validate:"required,max=50"`
	Company   string `json:"company,omitempty" validate:"max=100"`
	Address1  string `json:"address1" validate:"required,max=100"`
	Address2  string `json:"address2,omitempty" validate:"max=100"`
	City      string `json:"city" validate:"required,max=50"`
	State     string `json:"state" validate:"required,max=50"`
	PostalCode string `json:"postal_code" validate:"required,max=20"`
	Country   string `json:"country" validate:"required,len=2"`
	Phone     string `json:"phone,omitempty" validate:"max=20"`
}