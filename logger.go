package goweb

import (
	"fmt"
	"runtime"
)

const (
	// LogLevelDebug as defined in the RFC 5424 specification.
	LogLevelDebug = iota

	// LogLevelInfo as defined in the RFC 5424 specification.
	LogLevelInfo

	// LogLevelNotice as defined in the RFC 5424 specification.
	LogLevelNotice

	// LogLevelWarning as defined in the RFC 5424 specification.
	LogLevelWarning

	// LogLevelError as defined in the RFC 5424 specification.
	LogLevelError

	// LogLevelCritical as defined in the RFC 5424 specification.
	LogLevelCritical

	// LogLevelAlert as defined in the RFC 5424 specification.
	LogLevelAlert

	// LogLevelEmergency as defined in the RFC 5424 specification.
	LogLevelEmergency
)

// Logger outputs messages with Log() if ShouldLog()
// returns true.
type Logger interface {
	ShouldLog(level int) bool
	Log(level int, message interface{})
}

func getTitleAndColor(level int) (string, string) {
	switch level {
	case LogLevelDebug:
		return "Debug", "#aaaaaa"
	case LogLevelInfo:
		return "Info", "#439fe0"
	case LogLevelNotice:
		return "Notice", "#439fe0"
	case LogLevelWarning:
		return "Warning", "warning"
	case LogLevelError:
		return "Error", "danger"
	case LogLevelCritical:
		return "Critical", "danger"
	case LogLevelAlert:
		return "Alert", "danger"
	case LogLevelEmergency:
		return "Emergency", "danger"
	default:
		return "", ""
	}
}

func caller() string {
	_, f, l, _ := runtime.Caller(3)
	return fmt.Sprintf("%s#%d", f, l)
}
