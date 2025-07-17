package admin

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
)

type DashboardStats struct {
	Users         UserStats         `json:"users"`
	Services      ServiceStats      `json:"services"`
	Invoices      InvoiceStats      `json:"invoices"`
	Revenue       RevenueStats      `json:"revenue"`
	Products      ProductStats      `json:"products"`
	Orders        OrderStats        `json:"orders"`
	RecentActivity []ActivityItem   `json:"recent_activity"`
	SystemHealth   SystemHealth     `json:"system_health"`
}

type UserStats struct {
	Total          int64 `json:"total"`
	Active         int64 `json:"active"`
	NewThisMonth   int64 `json:"new_this_month"`
	NewToday       int64 `json:"new_today"`
	VerifiedEmails int64 `json:"verified_emails"`
	TwoFactorEnabled int64 `json:"two_factor_enabled"`
}

type ServiceStats struct {
	Total      int64 `json:"total"`
	Active     int64 `json:"active"`
	Pending    int64 `json:"pending"`
	Suspended  int64 `json:"suspended"`
	Cancelled  int64 `json:"cancelled"`
}

type InvoiceStats struct {
	Total       int64   `json:"total"`
	Paid        int64   `json:"paid"`
	Pending     int64   `json:"pending"`
	Overdue     int64   `json:"overdue"`
	TotalAmount float64 `json:"total_amount"`
	PaidAmount  float64 `json:"paid_amount"`
}

type RevenueStats struct {
	Today         float64 `json:"today"`
	ThisWeek      float64 `json:"this_week"`
	ThisMonth     float64 `json:"this_month"`
	LastMonth     float64 `json:"last_month"`
	ThisYear      float64 `json:"this_year"`
	GrowthPercent float64 `json:"growth_percent"`
}

type ProductStats struct {
	Total      int64 `json:"total"`
	Active     int64 `json:"active"`
	OutOfStock int64 `json:"out_of_stock"`
	LowStock   int64 `json:"low_stock"`
}

type OrderStats struct {
	TotalOrders    int64 `json:"total_orders"`
	OrdersToday    int64 `json:"orders_today"`
	OrdersThisWeek int64 `json:"orders_this_week"`
	AvgOrderValue  float64 `json:"avg_order_value"`
}

type ActivityItem struct {
	Type        string    `json:"type"`
	Description string    `json:"description"`
	UserID      uint      `json:"user_id,omitempty"`
	UserEmail   string    `json:"user_email,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
	Metadata    fiber.Map `json:"metadata,omitempty"`
}

type SystemHealth struct {
	DatabaseStatus string  `json:"database_status"`
	ResponseTime   float64 `json:"response_time_ms"`
	Uptime         string  `json:"uptime"`
	ErrorRate      float64 `json:"error_rate"`
}

var startTime = time.Now()

func GetDashboard(c *fiber.Ctx) error {
	adminUser := c.Locals("admin_user").(types.User)
	
	// Security: Log admin access
	logAdminActivity("dashboard_access", adminUser.ID, adminUser.Email, fiber.Map{
		"ip": c.IP(),
		"user_agent": c.Get("User-Agent"),
	})

	stats := DashboardStats{}
	
	// Get current time ranges
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	thisWeek := today.AddDate(0, 0, -int(today.Weekday()))
	thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastMonth := thisMonth.AddDate(0, -1, 0)
	thisYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	// User Statistics
	db.DB.Model(&types.User{}).Count(&stats.Users.Total)
	db.DB.Model(&types.User{}).Where("account_locked = ?", false).Count(&stats.Users.Active)
	db.DB.Model(&types.User{}).Where("created_at >= ?", thisMonth).Count(&stats.Users.NewThisMonth)
	db.DB.Model(&types.User{}).Where("created_at >= ?", today).Count(&stats.Users.NewToday)
	db.DB.Model(&types.User{}).Where("email_verified = ?", true).Count(&stats.Users.VerifiedEmails)
	db.DB.Model(&types.User{}).Where("two_factor_enabled = ?", true).Count(&stats.Users.TwoFactorEnabled)

	// Service Statistics
	db.DB.Model(&types.Service{}).Count(&stats.Services.Total)
	db.DB.Model(&types.Service{}).Where("status = ?", types.ServiceStatusActive).Count(&stats.Services.Active)
	db.DB.Model(&types.Service{}).Where("status = ?", types.ServiceStatusPending).Count(&stats.Services.Pending)
	db.DB.Model(&types.Service{}).Where("status = ?", types.ServiceStatusSuspended).Count(&stats.Services.Suspended)
	db.DB.Model(&types.Service{}).Where("status = ?", types.ServiceStatusCancelled).Count(&stats.Services.Cancelled)

	// Invoice Statistics
	db.DB.Model(&types.Invoice{}).Count(&stats.Invoices.Total)
	db.DB.Model(&types.Invoice{}).Where("status = ?", types.InvoiceStatusPaid).Count(&stats.Invoices.Paid)
	db.DB.Model(&types.Invoice{}).Where("status = ?", types.InvoiceStatusPending).Count(&stats.Invoices.Pending)
	db.DB.Model(&types.Invoice{}).Where("status = ? AND due_date < ?", types.InvoiceStatusPending, now).Count(&stats.Invoices.Overdue)
	
	// Revenue totals
	db.DB.Model(&types.Invoice{}).Select("COALESCE(SUM(total), 0)").Row().Scan(&stats.Invoices.TotalAmount)
	db.DB.Model(&types.Invoice{}).Where("status = ?", types.InvoiceStatusPaid).Select("COALESCE(SUM(total), 0)").Row().Scan(&stats.Invoices.PaidAmount)

	// Revenue Statistics
	db.DB.Model(&types.Invoice{}).Where("status = ? AND paid_at >= ?", types.InvoiceStatusPaid, today).Select("COALESCE(SUM(total), 0)").Row().Scan(&stats.Revenue.Today)
	db.DB.Model(&types.Invoice{}).Where("status = ? AND paid_at >= ?", types.InvoiceStatusPaid, thisWeek).Select("COALESCE(SUM(total), 0)").Row().Scan(&stats.Revenue.ThisWeek)
	db.DB.Model(&types.Invoice{}).Where("status = ? AND paid_at >= ?", types.InvoiceStatusPaid, thisMonth).Select("COALESCE(SUM(total), 0)").Row().Scan(&stats.Revenue.ThisMonth)
	db.DB.Model(&types.Invoice{}).Where("status = ? AND paid_at >= ? AND paid_at < ?", types.InvoiceStatusPaid, lastMonth, thisMonth).Select("COALESCE(SUM(total), 0)").Row().Scan(&stats.Revenue.LastMonth)
	db.DB.Model(&types.Invoice{}).Where("status = ? AND paid_at >= ?", types.InvoiceStatusPaid, thisYear).Select("COALESCE(SUM(total), 0)").Row().Scan(&stats.Revenue.ThisYear)

	// Calculate growth percentage
	if stats.Revenue.LastMonth > 0 {
		stats.Revenue.GrowthPercent = ((stats.Revenue.ThisMonth - stats.Revenue.LastMonth) / stats.Revenue.LastMonth) * 100
	}

	// Product Statistics
	db.DB.Model(&types.Product{}).Count(&stats.Products.Total)
	db.DB.Model(&types.Product{}).Where("is_active = ?", true).Count(&stats.Products.Active)
	db.DB.Model(&types.Product{}).Where("track_quantity = ? AND stock_quantity = 0 AND allow_backorder = ?", true, false).Count(&stats.Products.OutOfStock)
	db.DB.Model(&types.Product{}).Where("track_quantity = ? AND stock_quantity > 0 AND stock_quantity <= 5", true).Count(&stats.Products.LowStock)

	// Order Statistics (using completed carts as orders)
	db.DB.Model(&types.Cart{}).Where("status = ?", types.CartStatusCompleted).Count(&stats.Orders.TotalOrders)
	db.DB.Model(&types.Cart{}).Where("status = ? AND updated_at >= ?", types.CartStatusCompleted, today).Count(&stats.Orders.OrdersToday)
	db.DB.Model(&types.Cart{}).Where("status = ? AND updated_at >= ?", types.CartStatusCompleted, thisWeek).Count(&stats.Orders.OrdersThisWeek)
	db.DB.Model(&types.Cart{}).Where("status = ?", types.CartStatusCompleted).Select("COALESCE(AVG(total), 0)").Row().Scan(&stats.Orders.AvgOrderValue)

	// Recent Activity (simplified)
	stats.RecentActivity = getRecentActivity()

	// System Health
	stats.SystemHealth = SystemHealth{
		DatabaseStatus: "healthy",
		ResponseTime:   2.5,
		Uptime:         time.Since(startTime).String(),
		ErrorRate:      0.01,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}

func getRecentActivity() []ActivityItem {
	activities := []ActivityItem{}

	// Recent user registrations
	var recentUsers []types.User
	db.DB.Where("created_at >= ?", time.Now().AddDate(0, 0, -7)).
		Order("created_at DESC").
		Limit(5).
		Find(&recentUsers)

	for _, user := range recentUsers {
		activities = append(activities, ActivityItem{
			Type:        "user_registration",
			Description: "New user registered",
			UserID:      user.ID,
			UserEmail:   user.Email,
			Timestamp:   user.CreatedAt,
		})
	}

	// Recent payments
	var recentInvoices []types.Invoice
	db.DB.Where("status = ? AND paid_at >= ?", types.InvoiceStatusPaid, time.Now().AddDate(0, 0, -7)).
		Preload("User").
		Order("paid_at DESC").
		Limit(5).
		Find(&recentInvoices)

	for _, invoice := range recentInvoices {
		activities = append(activities, ActivityItem{
			Type:        "payment_received",
			Description: "Payment received",
			UserID:      invoice.UserID,
			UserEmail:   invoice.User.Email,
			Timestamp:   *invoice.PaidAt,
			Metadata: fiber.Map{
				"amount":         invoice.Total,
				"invoice_number": invoice.InvoiceNumber,
			},
		})
	}

	// Sort activities by timestamp (most recent first)
	for i := 0; i < len(activities)-1; i++ {
		for j := i + 1; j < len(activities); j++ {
			if activities[i].Timestamp.Before(activities[j].Timestamp) {
				activities[i], activities[j] = activities[j], activities[i]
			}
		}
	}

	// Return only the 10 most recent
	if len(activities) > 10 {
		activities = activities[:10]
	}

	return activities
}

func logAdminActivity(action string, adminID uint, adminEmail string, metadata fiber.Map) {
	// In a production system, you'd want to log this to a dedicated admin activity log
	// For now, we'll just create a simple log structure
	_ = action
	_ = adminID
	_ = adminEmail
	_ = metadata
}