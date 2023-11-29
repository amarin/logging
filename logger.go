package logging

import (
	"go.uber.org/zap"
)

// Logger wraps zap.SugaredLogger to hide zap requirements.
type zapLogger struct {
	*zap.SugaredLogger
	level Level
}

// Level returns current logger level.
func (logger zapLogger) Level() Level {
	return logger.level
}

// WithLevel returns a copy of logger with requested logging.Level set.
func (logger zapLogger) WithLevel(level Level) Logger {
	logger.level = level

	if !logger.IsEnabledForLevel(level) {
		logger.SugaredLogger = logger.SugaredLogger.Desugar().WithOptions(zap.IncreaseLevel(zapLevel(level))).Sugar()
	}

	return logger
}

// IsEnabledForLevel detects if internal logging level suitable to produce messages with specified logging.Level.
// Used to filter messages in Trace, Debug, Info, Warn, Error
// and formatting method companions Tracef, Debugf, Infof, Warnf, Errorf.
func (logger zapLogger) IsEnabledForLevel(level Level) bool {
	return logger.level.IsEnabledForLevel(level)
}

// Trace sends trace data onto logging.
func (logger zapLogger) Trace(args ...interface{}) {
	if logger.IsEnabledForLevel(LevelTrace) {
		logger.SugaredLogger.Debug(args...)
	}
}

// Tracef sends message template and filling arguments onto logging.
func (logger zapLogger) Tracef(fmt string, args ...interface{}) {
	if logger.IsEnabledForLevel(LevelTrace) {
		logger.SugaredLogger.Debugf(fmt, args...)
	}
}

// Debug sends debug message onto logging.
func (logger zapLogger) Debug(args ...interface{}) {
	if logger.IsEnabledForLevel(LevelDebug) {
		logger.SugaredLogger.Debug(args...)
	}
}

// Debugf sends message template and filling arguments onto logging.
func (logger zapLogger) Debugf(fmt string, args ...interface{}) {
	if logger.IsEnabledForLevel(LevelDebug) {
		logger.SugaredLogger.Debugf(fmt, args...)
	}
}

// Info sends trace data onto logging.
func (logger zapLogger) Info(args ...interface{}) {
	if logger.IsEnabledForLevel(LevelInfo) {
		logger.SugaredLogger.Info(args...)
	}
}

// Infof sends message template and filling arguments onto logging.
func (logger zapLogger) Infof(fmt string, args ...interface{}) {
	if logger.IsEnabledForLevel(LevelInfo) {
		logger.SugaredLogger.Infof(fmt, args...)
	}
}

// Warn sends trace data onto logging.
func (logger zapLogger) Warn(args ...interface{}) {
	if logger.IsEnabledForLevel(LevelWarn) {
		logger.SugaredLogger.Warn(args...)
	}
}

// Warnf sends message template and filling arguments onto logging.
func (logger zapLogger) Warnf(fmt string, args ...interface{}) {
	if logger.IsEnabledForLevel(LevelWarn) {
		logger.SugaredLogger.Warnf(fmt, args...)
	}
}

// Error sends error data onto logging.
func (logger zapLogger) Error(args ...interface{}) {
	logger.SugaredLogger.Error(args...)
}

// Errorf sends message template and filling arguments onto logging.
func (logger zapLogger) Errorf(fmt string, args ...interface{}) {
	logger.SugaredLogger.Errorf(fmt, args...)
}

// Fatal sends error data onto logging and calls os.exit(1).
func (logger zapLogger) Fatal(args ...interface{}) {
	logger.SugaredLogger.Fatal(args...)
}

// Fatalf sends message template and filling arguments onto logging and calls os.exit(1).
func (logger zapLogger) Fatalf(fmt string, args ...interface{}) {
	logger.SugaredLogger.Fatalf(fmt, args...)
}

// WithKeys provides a new logger instance having specified key-value set.
func (logger zapLogger) WithKeys(fields Keys) Logger {
	zapFields := make([]any, len(fields))
	idx := 0
	for k, v := range fields {
		zapFields[idx] = zap.Any(k.String(), v)
		idx++
	}

	return &zapLogger{
		SugaredLogger: logger.SugaredLogger.With(zapFields...),
		level:         logger.level,
	}
}

// WithKey provides a new logger instance having specified key-value set.
func (logger zapLogger) WithKey(key string, value any) Logger {
	zapField := zap.Any(key, value)

	return &zapLogger{
		SugaredLogger: logger.SugaredLogger.With(zapField),
		level:         logger.level,
	}
}

// WithError provides a new logger instance having specified error key.
func (logger zapLogger) WithError(err error) Logger {
	zapField := zap.Error(err)

	return &zapLogger{
		SugaredLogger: logger.SugaredLogger.WithOptions(zap.AddCaller()).With(zapField),
		level:         logger.level,
	}
}
