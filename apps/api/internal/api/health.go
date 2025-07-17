package api

import (
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/orvexcc/billing/api/internal/db"
	"github.com/orvexcc/billing/api/internal/types"
)

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Uptime    string            `json:"uptime"`
	Checks    map[string]Check  `json:"checks"`
	System    SystemInfo        `json:"system"`
}

type Check struct {
	Status  string        `json:"status"`
	Message string        `json:"message,omitempty"`
	Latency time.Duration `json:"latency,omitempty"`
}

type SystemInfo struct {
	GoVersion    string `json:"go_version"`
	NumGoroutine int    `json:"goroutines"`
	MemAlloc     uint64 `json:"memory_alloc_mb"`
	MemSys       uint64 `json:"memory_sys_mb"`
}

var startTime = time.Now()

func HealthCheck(c *fiber.Ctx) error {
	start := time.Now()
	checks := make(map[string]Check)
	overallStatus := "healthy"

	dbStart := time.Now()
	dbStatus := "healthy"
	dbMessage := "Database connection is working"
	
	var userCount int64
	if err := db.DB.Model(&types.User{}).Count(&userCount).Error; err != nil {
		dbStatus = "unhealthy"
		dbMessage = "Database query failed: " + err.Error()
		overallStatus = "unhealthy"
	}

	checks["database"] = Check{
		Status:  dbStatus,
		Message: dbMessage,
		Latency: time.Since(dbStart),
	}

	sessionStart := time.Now()
	sessionStatus := "healthy"
	sessionMessage := "Session storage is working"
	
	var sessionCount int64
	if err := db.DB.Model(&types.Session{}).Count(&sessionCount).Error; err != nil {
		sessionStatus = "unhealthy"
		sessionMessage = "Session storage query failed: " + err.Error()
		overallStatus = "unhealthy"
	}
	
	checks["sessions"] = Check{
		Status:  sessionStatus,
		Message: sessionMessage,
		Latency: time.Since(sessionStart),
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	response := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    time.Since(startTime).String(),
		Checks:    checks,
		System: SystemInfo{
			GoVersion:    runtime.Version(),
			NumGoroutine: runtime.NumGoroutine(),
			MemAlloc:     bToMb(m.Alloc),
			MemSys:       bToMb(m.Sys),
		},
	}

	statusCode := fiber.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = fiber.StatusServiceUnavailable
	}

	c.Set("X-Response-Time", time.Since(start).String())
	return c.Status(statusCode).JSON(response)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}