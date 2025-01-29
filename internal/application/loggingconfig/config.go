package loggingconfig

// LoggerConfigInterface defines the configuration needed for logging
type LoggerConfigInterface interface {
	GetEnv() string
}

// LoggerConfig defines the actual configuration structure
type LoggerConfig struct {
	Level string // Log level
}
