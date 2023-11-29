package logging_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amarin/logging"
)

func TestLevel_UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		level   logging.Level
		wantErr bool
	}{
		{"trace", logging.LevelTrace, false},
		{"trc", logging.LevelTrace, false},
		{"t", logging.LevelTrace, false},
		{"debug", logging.LevelDebug, false},
		{"dbg", logging.LevelDebug, false},
		{"d", logging.LevelDebug, false},
		{"info", logging.LevelInfo, false},
		{"inf", logging.LevelInfo, false},
		{"i", logging.LevelInfo, false},
		{"warning", logging.LevelWarn, false},
		{"warn", logging.LevelWarn, false},
		{"w", logging.LevelWarn, false},
		{"error", logging.LevelError, false},
		{"err", logging.LevelError, false},
		{"e", logging.LevelError, false},
		{"panic", logging.LevelPanic, false},
		{"p", logging.LevelPanic, false},
		{"fatal", logging.LevelFatal, false},
		{"f", logging.LevelFatal, false},
		{"any", logging.LevelFatal, true},
	}
	for _, tt := range tests {
		tt := tt
		var testNameSuffix string
		if tt.wantErr {
			testNameSuffix = "is_not_known"
		} else {
			testNameSuffix = "is_" + strings.ToLower(tt.level.String())
		}

		t.Run(tt.name+"_"+testNameSuffix, func(t *testing.T) {
			testLevel := new(logging.Level)
			err := testLevel.UnmarshalText([]byte(tt.name))
			require.Equal(t, err != nil, tt.wantErr)
			if err != nil {
				return
			}
			require.Equal(t, tt.level, *testLevel)
		})
	}
}

func TestLevel_String(t *testing.T) {
	tests := []struct {
		name string
		l    logging.Level
	}{
		{"trace", logging.LevelTrace},
		{"debug", logging.LevelDebug},
		{"info", logging.LevelInfo},
		{"warn", logging.LevelWarn},
		{"error", logging.LevelError},
		{"panic", logging.LevelPanic},
		{"fatal", logging.LevelFatal},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.name, tt.l.String())
		})
	}
}

func TestLevel_MarshalText(t *testing.T) {
	tests := []struct {
		name           string
		l              logging.Level
		wantTextString string
		wantErr        bool
	}{
		{logging.LevelTrace.String(), logging.LevelTrace, "trace", false},
		{logging.LevelDebug.String(), logging.LevelDebug, "debug", false},
		{logging.LevelInfo.String(), logging.LevelInfo, "info", false},
		{logging.LevelWarn.String(), logging.LevelWarn, "warn", false},
		{logging.LevelError.String(), logging.LevelError, "error", false},
		{logging.LevelPanic.String(), logging.LevelPanic, "panic", false},
		{logging.LevelFatal.String(), logging.LevelFatal, "fatal", false},
		{"unknown_level", 13, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotText, err := tt.l.MarshalText()
			require.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				return
			}
			require.Equal(t, tt.wantTextString, string(gotText))
		})
	}
}

func TestWithLevel(t *testing.T) {
	tests := []struct {
		level logging.Level
	}{
		{level: logging.LevelTrace},
		{level: logging.LevelDebug},
		{level: logging.LevelInfo},
		{level: logging.LevelWarn},
		{level: logging.LevelError},
		{level: logging.LevelPanic},
		{level: logging.LevelFatal},
	}
	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			cfg := logging.NewConfig()
			cfg.Apply(logging.WithLevel(tt.level))
			require.Equal(t, tt.level, cfg.Level)
		})
	}
}
