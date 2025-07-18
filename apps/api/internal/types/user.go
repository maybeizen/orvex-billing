package types

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UUID            uuid.UUID `json:"uuid" gorm:"type:uuid;uniqueIndex;not null"`
	Email           string    `json:"email" gorm:"uniqueIndex;not null"`
	Username        string    `json:"username" gorm:"uniqueIndex;not null"`
	PasswordHash    string    `json:"-" gorm:"not null"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Bio             string    `json:"bio" gorm:"type:text"`
	AvatarURL       string    `json:"avatar_url"`
	EmailVerified   bool      `json:"email_verified" gorm:"default:false"`
	TwoFactorSecret string    `json:"-"`
	TwoFactorEnabled bool     `json:"two_factor_enabled" gorm:"default:false"`
	AccountLocked   bool      `json:"account_locked" gorm:"default:false"`
	FailedAttempts  int       `json:"-" gorm:"default:0"`
	
	// User role and permissions
	Role            UserRole  `json:"role" gorm:"default:'user'"`
	
	// Notification preferences
	EmailNotifications      bool `json:"email_notifications" gorm:"default:true"`
	MarketingEmails        bool `json:"marketing_emails" gorm:"default:false"`
	SecurityNotifications  bool `json:"security_notifications" gorm:"default:true"`
	
	ProfilePublic          bool `json:"profile_public" gorm:"default:false"`
	ShowEmail              bool `json:"show_email" gorm:"default:false"`
	
	LastLoginAt     *time.Time `json:"last_login_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type Session struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BackupCode struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id" gorm:"not null;index"`
	User      User       `json:"user" gorm:"foreignKey:UserID"`
	Code      string     `json:"code" gorm:"not null;index"`
	Used      bool       `json:"used" gorm:"default:false"`
	UsedAt    *time.Time `json:"used_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
	UserRoleSuper UserRole = "super_admin"
)
