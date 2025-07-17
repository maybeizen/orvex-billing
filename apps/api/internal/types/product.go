package types

import (
	"time"
	"gorm.io/gorm"
)

type Category struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;uniqueIndex"`
	Slug        string    `json:"slug" gorm:"not null;uniqueIndex"`
	Description string    `json:"description" gorm:"type:text"`
	ImageURL    string    `json:"image_url"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	
	// Relationships
	Products    []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
	
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Product struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CategoryID  uint           `json:"category_id" gorm:"not null;index"`
	Category    Category       `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	
	// Basic product information
	Name        string         `json:"name" gorm:"not null"`
	Slug        string         `json:"slug" gorm:"not null;uniqueIndex"`
	Description string         `json:"description" gorm:"type:text"`
	ShortDesc   string         `json:"short_description" gorm:"type:varchar(255)"`
	
	// Pricing and inventory
	Price       float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	ComparePrice *float64      `json:"compare_price,omitempty" gorm:"type:decimal(10,2)"`
	Currency    string         `json:"currency" gorm:"default:'USD'"`
	SKU         string         `json:"sku" gorm:"uniqueIndex"`
	
	// Stock management
	StockQuantity    int       `json:"stock_quantity" gorm:"default:0"`
	TrackQuantity    bool      `json:"track_quantity" gorm:"default:true"`
	AllowBackorder   bool      `json:"allow_backorder" gorm:"default:false"`
	LowStockThreshold int      `json:"low_stock_threshold" gorm:"default:5"`
	
	// Per-user limitations
	MaxPerUser       *int      `json:"max_per_user,omitempty" gorm:"default:null"`
	MinQuantity      int       `json:"min_quantity" gorm:"default:1"`
	MaxQuantity      *int      `json:"max_quantity,omitempty" gorm:"default:null"`
	
	// Product status and visibility
	IsActive         bool      `json:"is_active" gorm:"default:true"`
	IsDigital        bool      `json:"is_digital" gorm:"default:false"`
	RequiresShipping bool      `json:"requires_shipping" gorm:"default:true"`
	IsFeatured       bool      `json:"is_featured" gorm:"default:false"`
	IsRecurring      bool      `json:"is_recurring" gorm:"default:false"`
	
	// Billing and subscription
	BillingCycle     *string   `json:"billing_cycle,omitempty" gorm:"type:varchar(20)"` // monthly, yearly, one_time
	TrialDays        *int      `json:"trial_days,omitempty"`
	SetupFee         *float64  `json:"setup_fee,omitempty" gorm:"type:decimal(10,2)"`
	
	// Resources and specifications
	Resources        ProductResources `json:"resources" gorm:"type:text"`
	
	// Panel integrations
	ConnectedPanels  ProductPanels    `json:"connected_panels" gorm:"type:text"`
	
	// Product upgrades and downgrades
	UpgradeOptions   ProductUpgrades  `json:"upgrade_options" gorm:"type:text"`
	
	// Media and SEO
	Images      ProductImages  `json:"images" gorm:"type:text"`
	MetaTitle   string         `json:"meta_title"`
	MetaDesc    string         `json:"meta_description"`
	
	// Service configuration for digital products
	ServiceConfig *ServiceConfig `json:"service_config,omitempty" gorm:"type:text"`
	
	// Relationships
	CartItems   []CartItem     `json:"cart_items,omitempty" gorm:"foreignKey:ProductID"`
	
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type ProductImages struct {
	Main      string   `json:"main,omitempty"`
	Gallery   []string `json:"gallery,omitempty"`
	Thumbnail string   `json:"thumbnail,omitempty"`
}

type ProductResources struct {
	CPU           *int     `json:"cpu_cores,omitempty"`
	RAM           *int     `json:"ram_gb,omitempty"`
	Storage       *int     `json:"storage_gb,omitempty"`
	Bandwidth     *int     `json:"bandwidth_gb,omitempty"`
	Databases     *int     `json:"databases,omitempty"`
	EmailAccounts *int     `json:"email_accounts,omitempty"`
	Domains       *int     `json:"domains,omitempty"`
	Subdomains    *int     `json:"subdomains,omitempty"`
	Features      []string `json:"features,omitempty"`
	CustomSpecs   map[string]interface{} `json:"custom_specs,omitempty"`
}

type ProductPanels struct {
	Pterodactyl  *PterodactylConfig `json:"pterodactyl,omitempty"`
	cPanel       *CPanelConfig      `json:"cpanel,omitempty"`
	Plesk        *PleskConfig       `json:"plesk,omitempty"`
	CustomPanels map[string]interface{} `json:"custom_panels,omitempty"`
}

type PterodactylConfig struct {
	ServerType   string            `json:"server_type,omitempty"`
	NodeID       *int              `json:"node_id,omitempty"`
	NestID       *int              `json:"nest_id,omitempty"`
	EggID        *int              `json:"egg_id,omitempty"`
	Allocations  int               `json:"allocations,omitempty"`
	Environment  map[string]string `json:"environment,omitempty"`
	StartupCmd   string            `json:"startup_cmd,omitempty"`
}

type CPanelConfig struct {
	Package      string `json:"package,omitempty"`
	QuotaGB      *int   `json:"quota_gb,omitempty"`
	BandwidthGB  *int   `json:"bandwidth_gb,omitempty"`
	EmailQuota   *int   `json:"email_quota,omitempty"`
	DatabaseQuota *int  `json:"database_quota,omitempty"`
}

type PleskConfig struct {
	ServicePlan  string `json:"service_plan,omitempty"`
	DiskSpaceGB  *int   `json:"disk_space_gb,omitempty"`
	TrafficGB    *int   `json:"traffic_gb,omitempty"`
	DomainLimit  *int   `json:"domain_limit,omitempty"`
	EmailLimit   *int   `json:"email_limit,omitempty"`
}

type ProductUpgrades struct {
	UpgradeTo    []ProductUpgrade `json:"upgrade_to,omitempty"`
	DowngradeTo  []ProductUpgrade `json:"downgrade_to,omitempty"`
	AutoUpgrade  *AutoUpgradeRule `json:"auto_upgrade,omitempty"`
}

type ProductUpgrade struct {
	ProductID    uint    `json:"product_id"`
	ProductName  string  `json:"product_name,omitempty"`
	PriceDiff    float64 `json:"price_difference"`
	Prorated     bool    `json:"prorated"`
	Description  string  `json:"description,omitempty"`
}

type AutoUpgradeRule struct {
	Enabled      bool   `json:"enabled"`
	Trigger      string `json:"trigger"` // resource_usage, time_based, manual
	Threshold    int    `json:"threshold,omitempty"` // percentage for resource usage
	UpgradeToID  uint   `json:"upgrade_to_id"`
}

// Request/Response types
type CategoryResponse struct {
	Success bool       `json:"success"`
	Data    []Category `json:"data"`
	Meta    struct {
		Total    int64 `json:"total"`
		Page     int   `json:"page"`
		PageSize int   `json:"page_size"`
	} `json:"meta,omitempty"`
}

type ProductResponse struct {
	Success bool      `json:"success"`
	Data    []Product `json:"data"`
	Meta    struct {
		Total    int64 `json:"total"`
		Page     int   `json:"page"`
		PageSize int   `json:"page_size"`
	} `json:"meta,omitempty"`
}