package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// FxLogger wraps the zap logger to implement the Logger interface
type FxLogger struct {
	*zap.SugaredLogger
}

// NewFxLogger creates a new FxLogger instance
func NewFxLogger(logger Logger) *FxLogger {
	// Assuming Logger has a method to get a SugaredLogger
	sugaredLogger, ok := logger.(interface{ SugaredLogger() *zap.SugaredLogger })
	if !ok {
		panic("logger does not implement SugaredLogger method")
	}
	return &FxLogger{SugaredLogger: sugaredLogger.SugaredLogger()} // Ensure FxLogger is properly initialized
}

// Log implements the Logger interface
func (l *FxLogger) Log(level zapcore.Level, msg string, fields ...interface{}) {
	switch level {
	case zapcore.DebugLevel:
		l.SugaredLogger.Debugw(msg, fields...)
	case zapcore.InfoLevel:
		l.SugaredLogger.Infow(msg, fields...)
	case zapcore.WarnLevel:
		l.SugaredLogger.Warnw(msg, fields...)
	case zapcore.ErrorLevel:
		l.SugaredLogger.Errorw(msg, fields...)
	default:
		l.SugaredLogger.Infow(msg, fields...) // Fallback to Info
	}
}

// LogLevel returns the log level
func (l *FxLogger) LogLevel() zapcore.Level {
	return zapcore.InfoLevel // Adjust as needed
}
