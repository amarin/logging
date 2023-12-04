package logging

import (
	"context"
)

// ContextExtractorFunc is a function returning field name
type ContextExtractorFunc func(ctx context.Context) (Key, any)

// WithContextExtractors adds context value to logging field extraction functions to logging config.
// Predefined Key's could be added with Key.Extractor method.
func WithContextExtractors(extractors ...ContextExtractorFunc) Option {
	return func(config *Config) {
		ctx := context.Background()
		for _, extractor := range extractors {
			k, _ := extractor(ctx)
			config.contextExtractors[k] = extractor
		}
	}
}

// contextKeys extracts logging keys from context using set with WithContextExtractors.
func (config *Config) contextKeys(ctx context.Context) Keys {
	res := make(Keys, len(config.contextExtractors))
	for _, e := range config.contextExtractors {
		if k, v := e(ctx); v != nil {
			res[k] = v
		}
	}

	return res
}

// loggableContextKey used internally to represent logging keys in context.
type loggableContextKey string

// loggingKey returns context key type for logging Key.
func (k Key) loggingKey() loggableContextKey {
	return loggableContextKey(k.String())
}

// SetToCtx sets Key value into context.
func (k Key) SetToCtx(ctx context.Context, value any) context.Context {
	return context.WithValue(ctx, k.loggingKey(), value)
}

// getFromCtx returns logging key value from context.
func (k Key) getFromCtx(ctx context.Context) any {
	return ctx.Value(k.loggingKey())
}

// Extractor returns ContextExtractorFunc for Key.
func (k Key) Extractor() ContextExtractorFunc {
	return func(ctx context.Context) (Key, any) {
		return k, k.getFromCtx(ctx)
	}
}

// KeysCtx returns logging keys from context using all extractors wet with WithContextExtractors.
func KeysCtx(ctx context.Context) Keys {
	mu.Lock()
	cfg := useConfig
	mu.Unlock()

	return cfg.contextKeys(ctx)
}
