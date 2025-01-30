package logging

import (
	"fmt"

	"go.uber.org/fx/fxevent"
)

// FxEventLogger is a custom logger that implements fxevent.Logger
type FxEventLogger struct {
	Logger Logger // Your existing logger interface
}

// LogEvent implements the fxevent.Logger interface
func (l *FxEventLogger) LogEvent(event fxevent.Event) {
	// Implement your logging logic here
	l.Logger.Debug("Event logged", String("event", fmt.Sprintf("%v", event)))
}
