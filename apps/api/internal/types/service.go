package types

import (
	"time"
	"gorm.io/gorm"
)

type Service struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	User        User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Status      ServiceStatus  `json:"status" gorm:"default:'pending'"`
	Price       float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	Currency    string         `json:"currency" gorm:"default:'USD'"`
	BillingType BillingType    `json:"billing_type" gorm:"default:'one_time'"`
	
	// Service configuration
	Config      ServiceConfig  `json:"config" gorm:"type:text"`
	
	// Billing information
	NextBillDate *time.Time    `json:"next_bill_date"`
	LastBillDate *time.Time    `json:"last_bill_date"`
	
	// Status tracking
	ActivatedAt  *time.Time    `json:"activated_at"`
	SuspendedAt  *time.Time    `json:"suspended_at"`
	CancelledAt  *time.Time    `json:"cancelled_at"`
	
	// Relationships
	Invoices     []Invoice     `json:"invoices,omitempty" gorm:"foreignKey:ServiceID"`
	
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type ServiceStatus string

const (
	ServiceStatusPending   ServiceStatus = "pending"
	ServiceStatusActive    ServiceStatus = "active"
	ServiceStatusSuspended ServiceStatus = "suspended"
	ServiceStatusCancelled ServiceStatus = "cancelled"
	ServiceStatusExpired   ServiceStatus = "expired"
)

type BillingType string

const (
	BillingTypeOneTime BillingType = "one_time"
	BillingTypeMonthly BillingType = "monthly"
	BillingTypeYearly  BillingType = "yearly"
)

type ServiceConfig struct {
	CPUCores int    `json:"cpu_cores,omitempty"`
	RAM      int    `json:"ram_gb,omitempty"`
	Storage  int    `json:"storage_gb,omitempty"`
	Features []string `json:"features,omitempty"`
	CustomData map[string]interface{} `json:"custom_data,omitempty"`
}