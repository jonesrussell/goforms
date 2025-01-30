package loggingconfig

// LoggerConfigInterface defines the methods for logger configuration
type LoggerConfigInterface interface {
	GetEnv() string
	// Add other configuration methods as needed
}

// Config holds the configuration for logging
type Config struct {
	env string
}

// NewConfig creates a new configuration instance
func NewConfig() LoggerConfigInterface {
	// Here you can load the environment variable or set a default
	return &Config{env: "development"} // Example: default to development
}

// GetEnv returns the current environment
func (c *Config) GetEnv() string {
	return c.env
}

// LoggerConfig holds the configuration for logging
type LoggerConfig struct {
	Level string // Add other fields as necessary
}

// GetEnv returns the environment level for logging
func (lc *LoggerConfig) GetEnv() string {
	return lc.Level
}
