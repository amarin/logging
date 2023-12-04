package logging

import "context"

// Logger defines default set of logging methods should provided by any logging backend.
type Logger interface {
	// Level returns current logger level.
	Level() Level

	// IsEnabledForLevel detects if internal logging level suitable to produce messages with specified logging.Level.
	// Used to filter messages in Trace, Debug, Info, Warn, Error
	// and formatting method companions Tracef, Debugf, Infof, Warnf, Errorf.
	IsEnabledForLevel(level Level) bool

	// Trace sends trace data onto logging.
	Trace(args ...interface{})

	// Tracef sends message template and filling arguments onto logging.
	Tracef(fmt string, args ...interface{})

	// Debug sends debug data onto logging.
	Debug(args ...interface{})

	// Debugf sends message template and filling arguments onto logging.
	Debugf(fmt string, args ...interface{})

	// Info sends info level data onto logging.
	Info(args ...interface{})

	// Infof sends message template and filling arguments onto logging.
	Infof(fmt string, args ...interface{})

	// Warn sends warn data onto logging.
	Warn(args ...interface{})

	// Warnf sends message template and filling arguments onto logging.
	Warnf(fmt string, args ...interface{})

	// Error sends error data onto logging.
	Error(args ...interface{})

	// Errorf sends message template and filling arguments onto logging.
	Errorf(fmt string, args ...interface{})

	// Fatal sends error data onto logging and calls os.exit(1).
	Fatal(args ...interface{})

	// Fatalf sends message template and filling arguments onto logging and calls os.exit(1).
	Fatalf(fmt string, args ...interface{})

	// WithKeys provides a new logger instance having specified key-value set.
	WithKeys(fields Keys) Logger

	// WithKey provides a new logger instance having specified single key-value pair set.
	WithKey(key string, value any) Logger

	// WithError provides a new logger instance having specified error key.
	WithError(err error) Logger

	// WithLevel provides a new logger instance inherit settings from parent except specified logging level.
	WithLevel(level Level) Logger

	// WithContext takes data from specified context. Uses configured ContextExtractorFunc's.
	WithContext(ctx context.Context) Logger
}
