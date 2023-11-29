package logging

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/imperfectgo/zap-syslog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/imperfectgo/zap-syslog/syslog"
)

// Backend implements logging.Backend using zap.Logger.
type Backend struct {
	_mu     sync.Mutex // protect global _core
	_core   *zap.Logger
	_config Config
}

// Init do logging backend initialisation. Returns error if initialisation failed.
// If core already initialized it just validates config and remembers it settings.
func (backend *Backend) Init(config Config) (err error) {
	var (
		encoder zapcore.Encoder
		syncer  zapcore.WriteSyncer
		_c      zapcore.Core
	)

	if backend._core != nil {
		return fmt.Errorf("%w: already configured", Error)
	}

	if encoder, err = backend.makeEncoder(config); err != nil {
		return err
	}

	if syncer, err = Output(config.Output.String()); err != nil {
		return err
	}

	_c = zapcore.NewCore(encoder, syncer, zapLevel(config.Level))

	backend._mu.Lock()
	backend._core = zap.New(_c).WithOptions(backend.makeOptions(config)...)
	backend._config = config
	backend._mu.Unlock()

	return nil
}

// MustInit do logging subsystem initialisation. Panics if initialisation failed.
func (backend *Backend) MustInit(config Config) {
	if err := backend.Init(config); err != nil {
		panic(err)
	}
}

// NewLogger initializes a new logger instance.
// Optionally takes level name to init.
// If no logging level specified DefaultLevel will set instead.
// Panics if global logging.Init is not called before.
// NOTE: this logger requires Sync() called manually to write any buffered log entryTable before exit.
// To Have automatically synced logger use NewLoggerCtx instead.
func (backend *Backend) NewLogger(levels ...Level) Logger {
	return backend.newLoggerForLevels(levels...)
}

// NewLoggerCtx initializes a new provider instance with automatic flushing any buffered log entryTable on context done.
// Takes logger context and optional logger level to set.
// If no logging level specified DefaultLevel will set instead.
// Panics if global logging.Init is not called before.
func (backend *Backend) NewLoggerCtx(ctx context.Context, levels ...Level) Logger {
	logger := backend.newLoggerForLevels(levels...)

	go func() {
		<-ctx.Done()

		_ = logger.Sync() //nolint:nolintlint,errcheck
	}()

	return logger
}

// NewNamedLogger initializes a new named logger.
// Takes logger name to initialize and optional logging level to set.
// If no logging level specified global configuration searched for specified logger name custom level,
// DefaultLevel will set if neither custom level configured nor specified with argument.
// NOTE: this logger requires Sync() called manually to write any buffered log entryTable before exit.
// To Have automatically synced logger use NewNamedLoggerCtx instead.
// Panics if global logging.Init is not called before.
func (backend *Backend) NewNamedLogger(name string, levels ...Level) (logger Logger) {
	if custom, ok := backend._config.CustomLevels[name]; ok {
		levels = append(levels, custom) // custom is not overlaps argument if provided, but can be first
	}

	return backend.makeNamed(backend.newLoggerForLevels(levels...), name)
}

// NewNamedLoggerCtx initializes a new named provider instance providing flushing buffered log entryTable on context done.
// Takes logger buffering context, logger name and optional logging level to set.
// If no logging level specified DefaultLevel will set instead.
// Panics if global logging.Init is not called before.
func (backend *Backend) NewNamedLoggerCtx(ctx context.Context, name string, levels ...Level) Logger {
	logger := (backend.NewNamedLogger(name, levels...)).(*zapLogger)

	go func() {
		<-ctx.Done()

		_ = logger.Sync() //nolint:errcheck
	}()

	return logger
}

// makeEncoder makes a zapcore.Encoder for zapcore configuration.
func (backend *Backend) makeEncoder(config Config) (zapcore.Encoder, error) {
	encoderConfig := backend.makeConfig(config)

	if config.Output == SysLog {
		syslogEncoderConfig := zapsyslog.SyslogEncoderConfig{
			EncoderConfig: encoderConfig,
			Facility:      syslog.LOG_DEBUG,
			Hostname:      "localhost",
			PID:           os.Getpid(),
			App:           os.Args[0],
		}
		return zapsyslog.NewSyslogEncoder(syslogEncoderConfig), nil
	}

	switch config.Format {
	case FormatText:
		return zapcore.NewConsoleEncoder(encoderConfig), nil
	case FormatJSON:
		return zapcore.NewJSONEncoder(encoderConfig), nil
	default:
		return nil, fmt.Errorf("%w: unknown format: %v", Error, config.Format)
	}
}

func (backend *Backend) makeOptions(config Config) []zap.Option {
	options := make([]zap.Option, 0)
	if config.Output == StdOut { // addEntry stacktrace only for console
		options = append(options, zap.AddStacktrace(zapcore.FatalLevel))
	}

	if config.Level == LevelDebug || config.Level == LevelTrace {
		options = append(options, zap.AddCaller())
	}

	options = append(options, zap.AddCallerSkip(1)) // increase caller frame distance as using per-logger level

	return options
}

func (backend *Backend) makeConfig(config Config) zapcore.EncoderConfig {
	var timeEncoder zapcore.TimeEncoder

	switch {
	case config.Output == SysLog:
		timeEncoder = zapcore.EpochTimeEncoder
	case config.Format == FormatText:
		timeEncoder = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(TimestampFormatConsole))
		}
	default: // assume JSON encode
		timeEncoder = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(TimestampFormatJSON))
		}
	}

	encoderConfig := zapcore.EncoderConfig{ //nolint:exhaustivestruct
		TimeKey:          KeyTimestamp.String(),
		LevelKey:         KeyLevel.String(),
		NameKey:          KeyLogger.String(),
		CallerKey:        KeyCaller.String(),
		FunctionKey:      zapcore.OmitKey,
		MessageKey:       KeyMessage.String(),
		StacktraceKey:    KeyStackTrace.String(),
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.LowercaseLevelEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		EncodeTime:       timeEncoder,
		ConsoleSeparator: " ",
	}

	return encoderConfig
}

// makeNamed makes new logger of original zap type.
func (backend *Backend) newLoggerForLevels(levels ...Level) *zapLogger {
	if backend._core == nil {
		backend.MustInit(*NewConfig())
	}

	level := DefaultLevel
	if len(levels) > 0 {
		level = levels[0]
	}

	return backend.newLogger(level, backend._core)
}

// makeNamed makes new logger of original zap type.
func (backend *Backend) newLogger(level Level, logger *zap.Logger) *zapLogger {
	return &zapLogger{SugaredLogger: logger.Sugar(), level: level}
}

// makeNamed makes named logger of original zap type.
func (backend *Backend) makeNamed(logger *zapLogger, name string) Logger {
	logger.SugaredLogger = logger.SugaredLogger.Named(name)

	return logger
}

// zapLevel maps logging Level to underlying zapcore.Level.
// NOTE: LevelTrace has no direct mapping onto zap logging level and mapped to zapcore.DebugLevel.
func zapLevel(l Level) zapcore.Level {
	switch l {
	case LevelTrace:
		return zapcore.DebugLevel
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelPanic:
		return zapcore.PanicLevel
	case LevelFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.Level(l)
	}
}
