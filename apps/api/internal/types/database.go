package types

type DatabaseConfig struct {
	Driver string
	Host string
	Port string
	Database string
	Username string
	Password string
	SSLMode string
}

type DatabaseInterface interface {
	// Add database methods here as needed
}