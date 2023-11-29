package logging

import (
	"context"
	"fmt"
	"sync"
)

var (
	mu           = new(sync.Mutex) // protect global usingBackend
	usingBackend *Backend
	initDone     bool
	useConfig    = NewConfig() // useConfig stores globally used config.
)

// SetBackend sets logging backend to use.
// Could be omitted when using Init of MustInit with backend provided.
func setBackend(newBackend *Backend) {
	mu.Lock()
	initDone = false
	usingBackend = newBackend
	mu.Unlock()
}

// Init do logging subsystem initialization. Returns error if init failed.
// If already initialized it just validates config and remembers it settings.
// Optionally takes single backend instance instead using SetBackend manually.
// Panics if backend is not set with SetBackend.
func Init(opts ...Option) (err error) {
	var backend *Backend

	config := NewConfig()
	config.Apply(opts...)

	if err = config.Validate(); err != nil {
		return err
	}

	backend = new(Backend)
	setBackend(backend)

	if err = backend.Init(*config); err != nil {
		return fmt.Errorf("%w: backend: %v", Error, err)
	}

	mu.Lock()
	initDone = true
	useConfig = config
	mu.Unlock()

	return nil
}

// MustInit do logging subsystem initialisation. Panics if init failed.
// Simply wraps Init doing panic on error.
// Takes configuration Option's to overwrite defaults.
func MustInit(opts ...Option) {
	if err := Init(opts...); err != nil {
		panic(err)
	}
}

// NewLogger initializes a new logger instance.
// Optionally takes level name to initDone.
// If no logging level specified DefaultLevel will set instead.
// Panics if global logging.Init is not called before.
// NOTE: this logger requires Sync() called manually to write any buffered log entryTable before exit.
// To Have automatically synced logger use NewLoggerCtx instead.
// Panics if global logging.Init is not called before or backend is not set with SetBackend.
func NewLogger(levels ...Level) Logger {
	mu.Lock()
	backend := usingBackend
	initialized := initDone
	mu.Unlock()

	switch {
	case backend == nil:
		panic(fmt.Errorf("%w: set backend first", Error))
	case !initialized:
		panic(fmt.Errorf("%w: init first", Error))
	default:
		return backend.NewLogger(levels...)
	}
}

// NewLoggerCtx initializes a new logger instance doing flushing any buffered log entries on context done.
// Takes logger context and optional logger level to set.
// If no logging level specified DefaultLevel will set instead.
// Automatically adds logging keys from context using any installed with WithContextExtractors.
// Panics if global logging.Init is not called before or backend is not set with SetBackend.
func NewLoggerCtx(ctx context.Context, levels ...Level) Logger {
	mu.Lock()
	backend := usingBackend
	initialized := initDone
	cfg := useConfig
	mu.Unlock()

	switch {
	case backend == nil:
		panic(fmt.Errorf("%w: set backend first", Error))
	case !initialized:
		panic(fmt.Errorf("%w: init first", Error))
	default:
		return backend.NewLoggerCtx(ctx, levels...).WithKeys(cfg.contextKeys(ctx))
	}
}

// NewNamedLogger initializes a new named logger.
// Takes logger name to initialize and optional logging level to set.
// If no logging level specified global configuration searched for specified logger name custom level,
// DefaultLevel will set if neither custom level configured nor specified with argument.
// NOTE: this logger may require usingBackend Sync() called manually to write any buffered log entryTable before exit.
// To Have automatically synced logger use NewNamedLoggerCtx instead.
// Panics if global logging.Init is not called before or backend is not set with SetBackend.
func NewNamedLogger(name string, levels ...Level) Logger {

	mu.Lock()
	backend := usingBackend
	initialized := initDone
	config := useConfig
	mu.Unlock()

	switch {
	case backend == nil:
		panic(fmt.Errorf("%w: set backend first", Error))
	case !initialized:
		panic(fmt.Errorf("%w: init first", Error))
	default:
		return backend.NewNamedLogger(name, config.levelForNamed(name, levels...))
	}
}

// NewNamedLoggerCtx initializes a new named provider instance providing flushing buffered log entryTable on context done.
// Takes logger buffering context, logger name and optional logging level to set.
// If no logging level specified DefaultLevel will set instead.
// Automatically adds logging keys from context using any installed with WithContextExtractors.
// Panics if global logging.Init is not called before or backend is not set with SetBackend.
func NewNamedLoggerCtx(ctx context.Context, name string, levels ...Level) Logger {
	mu.Lock()
	backend := usingBackend
	initialized := initDone
	cfg := useConfig
	mu.Unlock()

	switch {
	case backend == nil:
		panic(fmt.Errorf("%w: set backend first", Error))
	case !initialized:
		panic(fmt.Errorf("%w: init first", Error))
	default:
		return backend.NewNamedLoggerCtx(ctx, name, cfg.levelForNamed(name, levels...)).WithKeys(cfg.contextKeys(ctx))
	}
}
