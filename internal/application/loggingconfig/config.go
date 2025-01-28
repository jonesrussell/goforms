package loggingconfig

// LoggerConfig defines the configuration needed for logging
type LoggerConfig interface {
	GetEnv() string
}
