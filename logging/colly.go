package logging

import (
	"github.com/gocolly/colly/debug"
)

// Logger wrapper for Colly logs.
type CollyLogger struct {
	Log *Logger
}

// Initializes the logger.
func (l CollyLogger) Init() error {
	return nil
}

// Logs specific colly event.
func (l CollyLogger) Event(e *debug.Event) {
	l.Log.Debugf("%d [%6d - %s] %q\n", e.CollectorID, e.RequestID, e.Type, e.Values)
}
