package logging_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/logging"
)

var (
	Trace = logging.LevelTrace
	Debug = logging.LevelDebug
	Info  = logging.LevelInfo
	Warn  = logging.LevelWarn
	Err   = logging.LevelError
)

func TestNewLogger(t *testing.T) { //nolint:paralleltest
	backend := new(logging.Backend)
	for _, tt := range []struct { //nolint:paralleltest
		name     string
		arg      *logging.Level
		expected logging.Level
	}{
		{"default_logger", nil, logging.DefaultLevel},
		{"trace_logger", &Trace, Trace},
		{"debug_logger", &Debug, Debug},
		{"info_logger", &Info, Info},
		{"warn_logger", &Warn, Warn},
		{"error_logger", &Err, Err},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			args := make([]logging.Level, 0)

			if tt.arg != nil {
				args = append(args, *tt.arg)
				require.Len(t, args, 1)
			}

			got := backend.NewLogger(args...)
			require.Equal(t, tt.expected, got.Level())
		})
	}
}

func TestNewLoggerCtx(t *testing.T) { //nolint:paralleltest
	backend := new(logging.Backend)
	for _, tt := range []struct { //nolint:paralleltest
		name     string
		arg      *logging.Level
		expected logging.Level
	}{
		{"default_logger", nil, logging.DefaultLevel},
		{"trace_logger", &Trace, Trace},
		{"debug_logger", &Debug, Debug},
		{"info_logger", &Info, Info},
		{"warn_logger", &Warn, Warn},
		{"error_logger", &Err, Err},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			args := make([]logging.Level, 0)

			if tt.arg != nil {
				args = append(args, *tt.arg)
				require.Len(t, args, 1)
			}

			got := backend.NewLoggerCtx(context.TODO(), args...)
			require.Equal(t, tt.expected, got.Level())
		})
	}
}

func TestLogger_Levels(t *testing.T) { //nolint:paralleltest
	*(logging.CurrentConfig()) = *logging.NewConfig() // set default config
	backend := new(logging.Backend)
	config := logging.CurrentConfig()

	tests := []struct { //nolint:paralleltest
		loggerName    string
		argumentLevel *logging.Level
		customLevel   *logging.Level
		expectedLevel logging.Level
	}{
		{"arg_trace", &Trace, nil, Trace},
		{"arg_debug", &Debug, nil, Debug},
		{"arg_info", &Info, nil, Info},
		{"arg_warn", &Warn, nil, Warn},
		{"arg_error", &Err, nil, Err},
		{"arg_trace_overlaps_custom_debug", &Trace, &Debug, Trace},
		{"arg_trace_overlaps_custom_info", &Trace, &Info, Trace},
		{"arg_trace_overlaps_custom_warn", &Trace, &Warn, Trace},
		{"arg_trace_overlaps_custom_error", &Trace, &Err, Trace},
		{"arg_debug_overlaps_custom_trace", &Debug, &Trace, Debug},
		{"arg_debug_overlaps_custom_info", &Debug, &Info, Debug},
		{"arg_debug_overlaps_custom_warn", &Debug, &Warn, Debug},
		{"arg_debug_overlaps_custom_error", &Debug, &Err, Debug},
		{"arg_info_overlaps_custom_trace", &Info, &Trace, Info},
		{"arg_info_overlaps_custom_debug", &Info, &Debug, Info},
		{"arg_info_overlaps_custom_warn", &Info, &Warn, Info},
		{"arg_info_overlaps_custom_error", &Info, &Err, Info},
		{"arg_warn_overlaps_custom_trace", &Warn, &Trace, Warn},
		{"arg_warn_overlaps_custom_debug", &Warn, &Debug, Warn},
		{"arg_warn_overlaps_custom_info", &Warn, &Info, Warn},
		{"arg_warn_overlaps_custom_error", &Warn, &Err, Warn},
		{"arg_error_overlaps_custom_trace", &Err, &Trace, Err},
		{"arg_error_overlaps_custom_debug", &Err, &Debug, Err},
		{"arg_error_overlaps_custom_info", &Err, &Info, Err},
		{"arg_error_overlaps_custom_warn", &Err, &Warn, Err},
		{"default_level", nil, nil, logging.DefaultLevel},
	}
	for _, tt := range tests {
		if tt.customLevel != nil {
			config.CustomLevels[tt.loggerName] = *tt.customLevel
			require.Equal(t, *tt.customLevel, config.CustomLevels[tt.loggerName])
		}
	}
	// attach config
	require.NotPanics(t, func() { backend.MustInit(*config) })
	// run tests over initialized core
	for _, tt := range tests {
		tt := tt
		t.Run(tt.loggerName, func(t *testing.T) {
			var logger logging.Logger
			if tt.argumentLevel != nil {
				logger = backend.NewNamedLogger(tt.loggerName, *tt.argumentLevel)
			} else {
				logger = backend.NewNamedLogger(tt.loggerName)
			}

			require.Equal(t, tt.expectedLevel, logger.Level())
		})
	}
}

func TestLogger_WithKeys(t *testing.T) {
	backend := new(logging.Backend)
	logger := backend.NewLogger()
	logger.Info("test")
	logger.WithKeys(map[logging.Key]any{"key": "value"}).Info("test with map")
	logger.WithKeys(logging.Keys{"k1": "v1"}).Info("test with Keys")
}
