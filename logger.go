package goweb

// LogLevel as defined in the RFC 5424 specification.
type LogLevel int

const (
	// LogLevelDebug as defined in the RFC 5424 specification.
	LogLevelDebug LogLevel = 1

	// LogLevelInfo as defined in the RFC 5424 specification.
	LogLevelInfo LogLevel = 2

	// LogLevelNotice as defined in the RFC 5424 specification.
	LogLevelNotice LogLevel = 3

	// LogLevelWarning as defined in the RFC 5424 specification.
	LogLevelWarning LogLevel = 4

	// LogLevelError as defined in the RFC 5424 specification.
	LogLevelError LogLevel = 5

	// LogLevelCritical as defined in the RFC 5424 specification.
	LogLevelCritical LogLevel = 6

	// LogLevelAlert as defined in the RFC 5424 specification.
	LogLevelAlert LogLevel = 7

	// LogLevelEmergency as defined in the RFC 5424 specification.
	LogLevelEmergency LogLevel = 8
)

// Logger is an interface that implements
// Log(level int, message interface{}).
type Logger interface {
	Log(ctx *Context, level LogLevel, messages ...interface{})
}

func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelNotice:
		return "NOTICE"
	case LogLevelWarning:
		return "WARNING"
	case LogLevelError:
		return "ERROR"
	case LogLevelCritical:
		return "CRITICAL"
	case LogLevelAlert:
		return "ALERT"
	}
	return "EMERGENCY"
}
