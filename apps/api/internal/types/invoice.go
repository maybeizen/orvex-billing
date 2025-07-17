package types

import (
	"time"
	"gorm.io/gorm"
)

type Invoice struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	UserID          uint           `json:"user_id" gorm:"not null;index"`
	User            User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ServiceID       *uint          `json:"service_id" gorm:"index"`
	Service         *Service       `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
	
	// Invoice details
	InvoiceNumber   string         `json:"invoice_number" gorm:"uniqueIndex;not null"`
	Status          InvoiceStatus  `json:"status" gorm:"default:'pending'"`
	
	// Financial information
	SubTotal        float64        `json:"subtotal" gorm:"type:decimal(10,2);not null"`
	TaxAmount       float64        `json:"tax_amount" gorm:"type:decimal(10,2);default:0"`
	DiscountAmount  float64        `json:"discount_amount" gorm:"type:decimal(10,2);default:0"`
	Total           float64        `json:"total" gorm:"type:decimal(10,2);not null"`
	Currency        string         `json:"currency" gorm:"default:'USD'"`
	
	// Payment information
	PaymentMethod   string         `json:"payment_method"`
	PaymentID       string         `json:"payment_id"`
	StripePaymentID string         `json:"stripe_payment_id"`
	
	// Dates
	DueDate         time.Time      `json:"due_date" gorm:"not null"`
	PaidAt          *time.Time     `json:"paid_at"`
	
	// Invoice items
	Items           []InvoiceItem  `json:"items,omitempty" gorm:"foreignKey:InvoiceID"`
	
	// Metadata
	Notes           string         `json:"notes" gorm:"type:text"`
	
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type InvoiceStatus string

const (
	InvoiceStatusPending   InvoiceStatus = "pending"
	InvoiceStatusPaid      InvoiceStatus = "paid"
	InvoiceStatusOverdue   InvoiceStatus = "overdue"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
	InvoiceStatusRefunded  InvoiceStatus = "refunded"
)

type InvoiceItem struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	InvoiceID   uint    `json:"invoice_id" gorm:"not null;index"`
	Description string  `json:"description" gorm:"not null"`
	Quantity    int     `json:"quantity" gorm:"default:1"`
	UnitPrice   float64 `json:"unit_price" gorm:"type:decimal(10,2);not null"`
	Total       float64 `json:"total" gorm:"type:decimal(10,2);not null"`
	
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Request/Response types
type InvoiceListResponse struct {
	Success bool      `json:"success"`
	Data    []Invoice `json:"data"`
	Meta    struct {
		Total    int64 `json:"total"`
		Page     int   `json:"page"`
		PageSize int   `json:"page_size"`
	} `json:"meta"`
}

type PaymentRequest struct {
	PaymentMethod string `json:"payment_method" validate:"required,oneof=stripe paypal"`
	ReturnURL     string `json:"return_url,omitempty"`
}