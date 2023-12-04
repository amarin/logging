package logging

import (
	"context"
	"sync"
)

// Default logger provides logging methods with default setting in case customisation is not required.
// It still allows to make descendant fully-functional loggers when required.
type Default struct {
	mu     sync.Mutex
	logger Logger
}

// Level returns current logger level. Actually it always DefaultLevel unless forked using WithKeys or WithLevel.
func (l *Default) Level() Level {
	return DefaultLevel
}

// IsEnabledForLevel detects if internal logging level suitable to produce messages with specified logging.Level.
// Used to filter messages in Trace, Debug, Info, Warn, Error
// and formatting method companions Tracef, Debugf, Infof, Warnf, Errorf.
func (l *Default) IsEnabledForLevel(level Level) bool {
	l.checkInnerLogger()
	return l.logger.IsEnabledForLevel(level)
}

// WithKeys provides a new logger instance having specified key-value set.
func (l *Default) WithKeys(fields Keys) Logger {
	l.checkInnerLogger()
	return l.logger.WithKeys(fields)
}

// WithLevel provides a new logger instance inherit settings from parent except specified logging level.
func (l *Default) WithLevel(level Level) Logger {
	l.checkInnerLogger()
	return l.logger.WithLevel(level)
}

// WithName provides a new named logger instance with default settings.
func (l *Default) WithName(name string) Logger {
	l.checkInnerLogger()
	return NewNamedLogger(name, useConfig.Level)
}

// WithContext provides a new logger instance with context data attached.
func (l *Default) WithContext(ctx context.Context) Logger {
	l.checkInnerLogger()
	return NewLoggerCtx(ctx, useConfig.Level)
}

// Trace sends trace level data onto logging.
func (l *Default) Trace(args ...interface{}) {
	l.checkInnerLogger()
	l.logger.Trace(args...)
}

// Tracef sends trace level message template formatted with specified arguments.
func (l *Default) Tracef(fmt string, args ...interface{}) {
	l.checkInnerLogger()
	l.logger.Tracef(fmt, args)
}

// Debug sends debug level data onto logging.
func (l *Default) Debug(args ...interface{}) {
	l.checkInnerLogger()
	l.logger.Debug(args...)
}

// Debugf sends debug level message template formatted with specified arguments.
func (l *Default) Debugf(fmt string, args ...interface{}) {
	l.checkInnerLogger()
	l.logger.Debugf(fmt, args)
}

// Info sends info level data onto logging.
func (l *Default) Info(args ...interface{}) {
	l.checkInnerLogger()
	l.logger.Info(args...)
}

// Infof sends info level message template formatted with specified arguments.
func (l *Default) Infof(fmt string, args ...interface{}) {
	l.checkInnerLogger()
	l.logger.Infof(fmt, args)
}

// Warn sends warning level data onto logging.
func (l *Default) Warn(args ...interface{}) {
	l.checkInnerLogger()
	l.logger.Warn(args...)
}

// Warnf sends waning level message template formatted with specified arguments.
func (l *Default) Warnf(fmt string, args ...interface{}) {
	l.checkInnerLogger()
	l.logger.Warnf(fmt, args)
}

// Error sends error level data onto logging.
func (l *Default) Error(args ...interface{}) {
	l.checkInnerLogger()
	l.logger.Error(args...)
}

// Errorf sends error level message template formatted with specified arguments.
func (l *Default) Errorf(fmt string, args ...interface{}) {
	l.checkInnerLogger()
	l.logger.Errorf(fmt, args)
}

// Fatal sends fatal level data onto logging.
func (l *Default) Fatal(args ...interface{}) {
	l.checkInnerLogger()
	l.logger.Fatal(args...)
}

// Fatalf sends fatal level message template formatted with specified arguments.
func (l *Default) Fatalf(fmt string, args ...interface{}) {
	l.checkInnerLogger()
	l.logger.Fatalf(fmt, args)
}

// innerLogger creates internal logger if not exists.
func (l *Default) checkInnerLogger() {
	l.mu.Lock()
	if l.logger == nil {
		l.logger = NewLogger(useConfig.Level)
	}

	l.mu.Unlock()
}
