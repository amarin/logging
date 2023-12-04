package logging

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	levelTraceName   = "trace"
	levelDebugName   = "debug"
	levelInfoName    = "info"
	levelWarnName    = "warn"
	levelErrorName   = "error"
	levelPanicName   = "panic"
	levelFatalName   = "fatal"
	levelUnknownName = "level"
)

// Level is a type for logging levels.
type Level int

// Logging levels constants.
const (
	// LevelTrace defines constant logging level for most verbose logging in debug env.
	LevelTrace Level = 0

	// LevelDebug defines constant logging level for detailed logging in debug env.
	LevelDebug Level = 1

	// LevelInfo defines constant logging level for info, warnings and errors reporting.
	// It is used as DefaultLevel for all new loggers when level is not set directly.
	LevelInfo Level = 2

	// LevelWarn defines constant logging level for warnings and errors reporting only.
	LevelWarn Level = 3

	// LevelError defines constant logging level for only errors reporting.
	LevelError Level = 4

	// LevelPanic defines constant logging level for application panic reporting.
	LevelPanic Level = 5

	// LevelFatal defines constant logging level for application panic reporting.
	LevelFatal Level = 6

	// DefaultLevel defines default level if omitted.
	DefaultLevel = LevelInfo
)

// IsEnabledForLevel detects if logger with such Level suitable to produce messages with specified level.
// Used to filter messages in Logger.Trace, Logger.Debug, Logger.Info, Logger.Warn, Logger.Error
// and formatting method companions Logger.Tracef, Logger.Debugf, Logger.Infof, Logger.Warnf, Logger.Errorf.
func (l Level) IsEnabledForLevel(level Level) bool {
	return level >= l
}

// String returns level string representation. Implements fmt.Stringer.
func (l Level) String() string {
	switch l {
	case LevelTrace:
		return levelTraceName
	case LevelDebug:
		return levelDebugName
	case LevelInfo:
		return levelInfoName
	case LevelWarn:
		return levelWarnName
	case LevelError:
		return levelErrorName
	case LevelPanic:
		return levelPanicName
	case LevelFatal:
		return levelFatalName
	default:
		return levelUnknownName + strconv.Itoa(int(l))
	}
}

// UnmarshalText makes it easy to configure logging levels using YAML,
// TOML or JSON files.
func (l *Level) UnmarshalText(text []byte) error {
	if l == nil {
		return fmt.Errorf("`%w: unmarshal into nil", Error)
	}

	switch strings.ToLower(string(text)) {
	case levelTraceName, "t", "trc":
		*l = LevelTrace
	case levelDebugName, "d", "dbg":
		*l = LevelDebug
	case levelInfoName, "i", "inf":
		*l = LevelInfo
	case levelWarnName, "warning", "w", "wrn":
		*l = LevelWarn
	case levelErrorName, "e", "err":
		*l = LevelError
	case levelPanicName, "p":
		*l = LevelPanic
	case levelFatalName, "f":
		*l = LevelFatal
	default:
		return fmt.Errorf("%w: unrecognized level: %q", Error, text)
	}

	return nil
}

// MarshalText marshals level value to text representation.
// Implements encoding.TextMarshaler.
func (l Level) MarshalText() (text []byte, err error) {
	switch l {
	case LevelTrace:
		return []byte(levelTraceName), nil
	case LevelDebug:
		return []byte(levelDebugName), nil
	case LevelInfo:
		return []byte(levelInfoName), nil
	case LevelWarn:
		return []byte(levelWarnName), nil
	case LevelError:
		return []byte(levelErrorName), nil
	case LevelPanic:
		return []byte(levelPanicName), nil
	case LevelFatal:
		return []byte(levelFatalName), nil
	default:
		return nil, fmt.Errorf("%w: unknown level value %v", Error, l)
	}
}

// WithLevel adds specified level to configuration.
func WithLevel(level Level) Option {
	return func(c *Config) {
		c.Level = level
	}
}
