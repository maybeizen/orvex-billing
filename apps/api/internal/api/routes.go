package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/api/v1/admin"
	"github.com/orvexcc/billing/api/internal/api/v1/auth"
	"github.com/orvexcc/billing/api/internal/api/v1/products"
	"github.com/orvexcc/billing/api/internal/api/v1/services"
	"github.com/orvexcc/billing/api/internal/api/v1/user"
	"github.com/orvexcc/billing/api/internal/middleware"
)

func RegisterRoutes(app *fiber.App) {
	// Add CORS middleware globally
	app.Use(middleware.CORS())
	
	// Static file serving for uploads
	app.Static("/uploads", "./uploads")
	
	api := app.Group("/api")
	
	api.Get("/health", HealthCheck)

	v1 := api.Group("/v1")
	
	authGroup := v1.Group("/auth")
	authGroup.Post("/login", auth.Login)
	authGroup.Post("/register", auth.Register)
	authGroup.Post("/logout", middleware.RequireAuth, auth.Logout)
	authGroup.Post("/logout-all", middleware.RequireAuth, auth.LogoutAll)
	authGroup.Post("/2fa/verify", middleware.RequireAuth, auth.Verify2FA)

	protected := v1.Group("/", middleware.RequireAuth)
	
	userGroup := protected.Group("/user")
	
	// Profile management
	userGroup.Get("/profile", user.GetProfile)
	userGroup.Put("/profile", user.UpdateProfile)
	userGroup.Put("/username", user.UpdateUsername)
	userGroup.Put("/email", user.UpdateEmail)
	userGroup.Put("/password", user.ChangePassword)
	
	// Avatar management
	userGroup.Post("/avatar", user.UploadAvatar)
	userGroup.Delete("/avatar", user.DeleteAvatar)
	
	// Preferences and settings
	userGroup.Put("/notifications", user.UpdateNotificationPreferences)
	userGroup.Put("/privacy", user.UpdatePrivacySettings)
	
	// Two-factor authentication
	userGroup.Post("/2fa/setup", auth.Setup2FA)
	userGroup.Post("/2fa/enable", auth.Enable2FA)
	userGroup.Post("/2fa/verify", auth.Verify2FA)
	userGroup.Post("/2fa/backup-verify", auth.VerifyBackupCode)
	userGroup.Post("/2fa/disable", auth.Disable2FA)
	userGroup.Post("/2fa/backup-codes", auth.GenerateBackupCodes)
	
	// Account deletion
	userGroup.Delete("/account", user.DeleteAccount)

	// Public user lookup by UUID (optional authentication)
	api.Get("/users/:uuid", middleware.OptionalAuth, user.GetUserByUUID)

	// Admin routes (protected by admin role)
	adminGroup := protected.Group("/admin", middleware.RequireAdmin)
	
	// Admin dashboard
	adminGroup.Get("/dashboard", admin.GetDashboard)
	
	// Admin user management
	adminGroup.Get("/users", admin.ListUsers)
	adminGroup.Get("/users/:id", admin.GetUser)
	adminGroup.Post("/users", admin.CreateUser)
	adminGroup.Put("/users/:id", admin.UpdateUser)
	adminGroup.Delete("/users/:id", admin.DeleteUser)
	
	// Bulk user actions
	adminGroup.Post("/users/bulk", admin.BulkUserActions)
	
	// User impersonation
	adminGroup.Post("/users/:id/impersonate", admin.ImpersonateUser)
	adminGroup.Post("/impersonation/stop", admin.StopImpersonation)
	
	// Admin service management
	adminGroup.Get("/services", admin.ListAdminServices)
	adminGroup.Get("/services/:id", admin.GetAdminService)
	adminGroup.Post("/services", admin.CreateAdminService)
	adminGroup.Put("/services/:id", admin.UpdateAdminService)
	adminGroup.Delete("/services/:id", admin.DeleteAdminService)
	adminGroup.Post("/services/:id/action", admin.ServiceAction)
	
	// Admin invoice management
	adminGroup.Get("/invoices", admin.ListAdminInvoices)
	adminGroup.Get("/invoices/:id", admin.GetAdminInvoice)
	adminGroup.Post("/invoices", admin.CreateAdminInvoice)
	adminGroup.Put("/invoices/:id", admin.UpdateAdminInvoice)
	adminGroup.Delete("/invoices/:id", admin.DeleteAdminInvoice)
	adminGroup.Post("/invoices/:id/action", admin.InvoiceAction)

	// Admin product management
	adminGroup.Get("/products", admin.ListAdminProducts)
	adminGroup.Get("/products/:id", admin.GetAdminProduct)
	adminGroup.Post("/products", admin.CreateAdminProduct)
	adminGroup.Put("/products/:id", admin.UpdateAdminProduct)
	adminGroup.Delete("/products/:id", admin.DeleteAdminProduct)
	adminGroup.Post("/products/:id/duplicate", admin.DuplicateAdminProduct)
	adminGroup.Post("/products/bulk", admin.BulkUpdateProducts)

	// Admin category management
	adminGroup.Get("/categories", admin.ListAdminCategories)
	adminGroup.Get("/categories/:id", admin.GetAdminCategory)
	adminGroup.Post("/categories", admin.CreateAdminCategory)
	adminGroup.Put("/categories/:id", admin.UpdateAdminCategory)
	adminGroup.Delete("/categories/:id", admin.DeleteAdminCategory)
	adminGroup.Post("/categories/reorder", admin.ReorderCategories)
	adminGroup.Post("/categories/bulk", admin.BulkUpdateCategories)
	adminGroup.Get("/categories/:id/products", admin.GetCategoryProductsAdmin)

	// Admin coupon management
	adminGroup.Get("/coupons", admin.ListAdminCoupons)
	adminGroup.Get("/coupons/:id", admin.GetAdminCoupon)
	adminGroup.Post("/coupons", admin.CreateAdminCoupon)
	adminGroup.Put("/coupons/:id", admin.UpdateAdminCoupon)
	adminGroup.Delete("/coupons/:id", admin.DeleteAdminCoupon)
	adminGroup.Post("/coupons/bulk", admin.BulkUpdateCoupons)
	adminGroup.Get("/coupons/:id/usage", admin.GetCouponUsageHistory)

	announcementsGroup := api.Group("/announcements")
	announcementsGroup.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Announcements list endpoint"})
	})
	announcementsGroup.Get("/:id", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Announcement detail endpoint"})
	})

	// Product and Category endpoints (public)
	categoriesGroup := api.Group("/categories")
	categoriesGroup.Get("/", products.ListCategories)
	categoriesGroup.Get("/:slug", products.GetCategory)
	categoriesGroup.Get("/:slug/products", products.GetCategoryProducts)

	productsGroup := api.Group("/products")
	productsGroup.Get("/", products.ListProducts)
	productsGroup.Get("/featured", products.GetFeaturedProducts)
	productsGroup.Get("/search", products.SearchProducts)
	productsGroup.Get("/:slug", products.GetProduct)

	// Cart endpoints (protected)
	cartGroup := api.Group("/cart", middleware.RequireAuth)
	cartGroup.Get("/", products.GetCart)
	cartGroup.Post("/items", products.AddToCart)
	cartGroup.Put("/items/:itemId", products.UpdateCartItem)
	cartGroup.Delete("/items/:itemId", products.RemoveFromCart)
	cartGroup.Delete("/clear", products.ClearCart)

	// Coupon endpoints (protected)
	cartGroup.Post("/coupons", products.ApplyCoupon)
	cartGroup.Delete("/coupons", products.RemoveCoupon)
	cartGroup.Get("/coupons/:code/validate", products.ValidateCouponCode)

	// Public coupon endpoints (anyone can view)
	api.Get("/coupons", products.ListPublicCoupons)
	api.Get("/coupons/:code", products.GetPublicCouponByCode)

	// Checkout endpoints (protected)
	checkoutGroup := api.Group("/checkout", middleware.RequireAuth)
	checkoutGroup.Get("/summary", products.GetCheckoutSummary)
	checkoutGroup.Post("/", products.InitiateCheckout)

	// Services endpoints
	servicesGroup := api.Group("/services", middleware.RequireAuth)
	servicesGroup.Get("/", services.ListServices)
	servicesGroup.Get("/:id", services.GetService)
	servicesGroup.Get("/:id/status", services.GetServiceStatus)
	
	// Invoice endpoints
	invoicesGroup := api.Group("/invoices", middleware.RequireAuth)
	invoicesGroup.Get("/", services.ListInvoices)
	invoicesGroup.Get("/:id", services.GetInvoice)
	invoicesGroup.Post("/:id/pay", services.PayInvoice)
	
	// Webhook endpoints (no authentication for payment processors)
	api.Post("/webhooks/payment", services.PaymentWebhook)
}