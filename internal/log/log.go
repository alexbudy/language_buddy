package log

import "log"

// DebugEnabled controls whether debug messages are logged.
var DebugEnabled bool

// Debug logs a formatted debug message when debugging is enabled.
func Debug(format string, args ...any) {
	if DebugEnabled {
		log.Printf("[DEBUG] " + format, args...)
	}
}