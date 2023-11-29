package logging_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/logging"
)

func TestMustInit_Panic(t *testing.T) {
	for _, tt := range []struct {
		name        string
		config      []logging.Option
		shouldPanic bool
	}{
		{
			"ok with default params",
			nil,
			false,
		},
		{
			"panics with invalid format",
			[]logging.Option{logging.WithFormat("someFormat")},
			true,
		},
		{
			"panics with invalid level",
			[]logging.Option{logging.WithLevel(logging.Level(99))},
			true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				require.Panicsf(t, func() { logging.MustInit(tt.config...) }, "config: %v", tt.config)
			} else {
				require.NotPanicsf(t, func() { logging.MustInit(tt.config...) }, "config: %v", tt.config)
			}
		})
	}
}
