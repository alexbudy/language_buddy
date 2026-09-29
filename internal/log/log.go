package log

import "log"

// DebugEnabled controls whether debug messages are logged.
var DebugEnabled bool

// DebugSQLEnabled controls whether SQL debug messages are logged.
var DebugSQLEnabled bool

// Debug logs a formatted debug message when debugging is enabled.
func Debug(format string, args ...any) {
	if DebugEnabled {
		log.Printf("[DEBUG] " + format, args...)
	}
}

// DebugSQL logs a formatted SQL debug message when SQL debugging is enabled.
func DebugSQL(format string, args ...any) {
	if DebugSQLEnabled {
		log.Printf("[DEBUG-SQL] " + format, args...)
	}
}