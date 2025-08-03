package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	WhatsApp WhatsAppConfig `json:"whatsapp"`
	Log      LogConfig      `json:"log"`
	Security SecurityConfig `json:"security"`
	Features FeaturesConfig `json:"features"`
	Auth     AuthConfig     `json:"auth"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Host         string          `json:"host"`
	Port         int             `json:"port"`
	ReadTimeout  time.Duration   `json:"read_timeout"`
	WriteTimeout time.Duration   `json:"write_timeout"`
	IdleTimeout  time.Duration   `json:"idle_timeout"`
	CORS         CORSConfig      `json:"cors"`
	RateLimit    RateLimitConfig `json:"rate_limit"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Driver          string        `json:"driver"`            // "sqlite", "postgres"
	URL             string        `json:"url"`               // Connection string
	AutoMigrate     bool          `json:"auto_migrate"`      // Auto-run migrations
	MaxOpenConns    int           `json:"max_open_conns"`    // Max open connections
	MaxIdleConns    int           `json:"max_idle_conns"`    // Max idle connections
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"` // Connection max lifetime

	// PostgreSQL specific settings
	PostgreSQL PostgreSQLConfig `json:"postgresql"`

	// SQLite specific settings
	SQLite SQLiteConfig `json:"sqlite"`
}

// PostgreSQLConfig represents PostgreSQL specific configuration
type PostgreSQLConfig struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Database   string `json:"database"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	SSLMode    string `json:"ssl_mode"`    // disable, require, verify-ca, verify-full
	TimeZone   string `json:"timezone"`    // Default timezone
	SearchPath string `json:"search_path"` // Schema search path
	AppName    string `json:"app_name"`    // Application name for connection
}

// SQLiteConfig represents SQLite specific configuration
type SQLiteConfig struct {
	Path        string `json:"path"`         // Database file path
	ForeignKeys bool   `json:"foreign_keys"` // Enable foreign key constraints
	JournalMode string `json:"journal_mode"` // WAL, DELETE, TRUNCATE, PERSIST, MEMORY, OFF
	Synchronous string `json:"synchronous"`  // OFF, NORMAL, FULL, EXTRA
	CacheSize   int    `json:"cache_size"`   // Page cache size
	TempStore   string `json:"temp_store"`   // DEFAULT, FILE, MEMORY
	MmapSize    int64  `json:"mmap_size"`    // Memory-mapped I/O size
}

// WhatsAppConfig represents WhatsApp configuration
type WhatsAppConfig struct {
	LogLevel       string        `json:"log_level"`
	QRTimeout      time.Duration `json:"qr_timeout"`
	ReconnectDelay time.Duration `json:"reconnect_delay"`
	MaxReconnects  int           `json:"max_reconnects"`
}

// LogConfig represents logging configuration
type LogConfig struct {
	Level         string `json:"level"`
	Output        string `json:"output"`         // "console", "file", or "dual"
	ConsoleFormat string `json:"console_format"` // "console" (colorido) or "json" (estruturado)
	FileFormat    string `json:"file_format"`    // "console" or "json" (recomendado: json)
	TimeFormat    string `json:"time_format"`
	Caller        bool   `json:"caller"`
	StackTrace    bool   `json:"stack_trace"`
	FilePath      string `json:"file_path"`   // Path for file logging
	MaxSize       int    `json:"max_size"`    // Maximum size in MB before rotation
	MaxBackups    int    `json:"max_backups"` // Maximum number of backup files
	MaxAge        int    `json:"max_age"`     // Maximum age in days
	PrettyJSON    bool   `json:"pretty_json"` // Enable pretty printing for JSON file output
}

// CORSConfig represents CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string `json:"allowed_origins"`
	AllowedMethods   []string `json:"allowed_methods"`
	AllowedHeaders   []string `json:"allowed_headers"`
	AllowCredentials bool     `json:"allow_credentials"`
	MaxAge           int      `json:"max_age"`
}

// RateLimitConfig represents rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int `json:"requests_per_minute"`
	BurstSize         int `json:"burst_size"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Enabled    bool            `json:"enabled"`
	Type       string          `json:"type"` // "api_key" or "basic"
	APIKeys    []string        `json:"api_keys"`
	HeaderName string          `json:"header_name"`
	BasicAuth  BasicAuthConfig `json:"basic_auth"`
}

// BasicAuthConfig represents basic authentication configuration
type BasicAuthConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
	JWTSecret string `json:"jwt_secret,omitempty"`
	APIKey    string `json:"api_key,omitempty"`
}

// FeaturesConfig represents feature flags
type FeaturesConfig struct {
	EnableMetrics  bool `json:"enable_metrics"`
	EnableWebhooks bool `json:"enable_webhooks"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Host:         getEnvString("SERVER_HOST", "localhost"),
			Port:         getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getEnvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
			CORS: CORSConfig{
				AllowedOrigins:   getEnvStringSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
				AllowedMethods:   getEnvStringSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
				AllowedHeaders:   getEnvStringSlice("CORS_ALLOWED_HEADERS", []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}),
				AllowCredentials: getEnvBool("CORS_ALLOW_CREDENTIALS", false),
				MaxAge:           getEnvInt("CORS_MAX_AGE", 86400),
			},
			RateLimit: RateLimitConfig{
				RequestsPerMinute: getEnvInt("RATE_LIMIT_REQUESTS", 100),
				BurstSize:         getEnvInt("RATE_LIMIT_BURST_SIZE", 10),
			},
		},
		Database: DatabaseConfig{
			Driver:          getEnvString("DB_DRIVER", "sqlite3"),
			URL:             getEnvString("DB_URL", "./data/wazmeow.db"),
			AutoMigrate:     getEnvBool("DB_AUTO_MIGRATE", true),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			PostgreSQL: PostgreSQLConfig{
				Host:       getEnvString("POSTGRES_HOST", "localhost"),
				Port:       getEnvInt("POSTGRES_PORT", 5432),
				Database:   getEnvString("POSTGRES_DB", "wazmeow"),
				Username:   getEnvString("POSTGRES_USER", "wazmeow"),
				Password:   getEnvString("POSTGRES_PASSWORD", ""),
				SSLMode:    getEnvString("POSTGRES_SSL_MODE", "disable"),
				TimeZone:   getEnvString("POSTGRES_TIMEZONE", "UTC"),
				SearchPath: getEnvString("POSTGRES_SEARCH_PATH", "public"),
				AppName:    getEnvString("POSTGRES_APP_NAME", "wazmeow"),
			},
			SQLite: SQLiteConfig{
				Path:        getEnvString("SQLITE_PATH", "./data/wazmeow.db"),
				ForeignKeys: getEnvBool("SQLITE_FOREIGN_KEYS", true),
				JournalMode: getEnvString("SQLITE_JOURNAL_MODE", "WAL"),
				Synchronous: getEnvString("SQLITE_SYNCHRONOUS", "NORMAL"),
				CacheSize:   getEnvInt("SQLITE_CACHE_SIZE", 1000),
				TempStore:   getEnvString("SQLITE_TEMP_STORE", "memory"),
				MmapSize:    getEnvInt64("SQLITE_MMAP_SIZE", 268435456),
			},
		},
		WhatsApp: WhatsAppConfig{
			LogLevel:       getEnvString("WHATSAPP_LOG_LEVEL", "INFO"),
			QRTimeout:      getEnvDuration("WHATSAPP_QR_TIMEOUT", 5*time.Minute),
			ReconnectDelay: getEnvDuration("WHATSAPP_RECONNECT_DELAY", 5*time.Second),
			MaxReconnects:  getEnvInt("WHATSAPP_MAX_RECONNECTS", 3),
		},
		Log: LogConfig{
			Level:         getEnvString("LOG_LEVEL", "info"),
			Output:        getEnvString("LOG_OUTPUT", "console"),
			ConsoleFormat: getEnvString("LOG_CONSOLE_FORMAT", "console"),
			FileFormat:    getEnvString("LOG_FILE_FORMAT", "json"),
			TimeFormat:    getEnvString("LOG_TIME_FORMAT", time.RFC3339),
			Caller:        getEnvBool("LOG_CALLER", false),
			StackTrace:    getEnvBool("LOG_STACK_TRACE", false),
			FilePath:      getEnvString("LOG_FILE_PATH", "./logs/wazmeow.log"),
			MaxSize:       getEnvInt("LOG_MAX_SIZE", 100),
			MaxBackups:    getEnvInt("LOG_MAX_BACKUPS", 3),
			MaxAge:        getEnvInt("LOG_MAX_AGE", 28),
		},
		Security: SecurityConfig{
			JWTSecret: getEnvString("JWT_SECRET", ""),
			APIKey:    getEnvString("API_KEY", ""),
		},
		Features: FeaturesConfig{
			EnableMetrics:  getEnvBool("ENABLE_METRICS", false),
			EnableWebhooks: getEnvBool("ENABLE_WEBHOOKS", false),
		},
		Auth: AuthConfig{
			Enabled:    getEnvBool("AUTH_ENABLED", false),
			Type:       getEnvString("AUTH_TYPE", "api_key"),
			APIKeys:    getEnvStringSlice("AUTH_API_KEYS", []string{}),
			HeaderName: getEnvString("AUTH_HEADER_NAME", "X-API-Key"),
			BasicAuth: BasicAuthConfig{
				Username: getEnvString("AUTH_BASIC_USERNAME", ""),
				Password: getEnvString("AUTH_BASIC_PASSWORD", ""),
			},
		},
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Database.Driver == "" {
		return fmt.Errorf("database driver is required")
	}

	if c.Database.URL == "" {
		return fmt.Errorf("database URL is required")
	}

	if c.Log.Level == "" {
		return fmt.Errorf("log level is required")
	}

	validLogLevels := []string{"debug", "info", "warn", "error", "fatal"}
	if !contains(validLogLevels, c.Log.Level) {
		return fmt.Errorf("invalid log level: %s", c.Log.Level)
	}

	validLogOutputs := []string{"console", "file", "dual"}
	if !contains(validLogOutputs, c.Log.Output) {
		return fmt.Errorf("invalid log output: %s", c.Log.Output)
	}

	validLogFormats := []string{"json", "console"}
	if !contains(validLogFormats, c.Log.ConsoleFormat) {
		return fmt.Errorf("invalid console log format: %s", c.Log.ConsoleFormat)
	}

	if !contains(validLogFormats, c.Log.FileFormat) {
		return fmt.Errorf("invalid file log format: %s", c.Log.FileFormat)
	}

	return nil
}

// GetServerAddress returns the server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return getEnvString("ENVIRONMENT", "development") == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return getEnvString("ENVIRONMENT", "development") == "production"
}

// Helper functions for environment variable parsing

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvStringSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
