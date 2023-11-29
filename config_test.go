package logging_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	. "github.com/amarin/logging"
)

//nolint:funlen
func TestConfig_UnmarshalYAML(t *testing.T) {
	emptyCustomLevels := make(map[string]Level)
	for _, tt := range []struct { //nolint:paralleltest
		name        string
		bytes       string
		level       Level
		out         Target
		format      Format
		NamedLevels map[string]Level
		wantErr     bool
	}{
		{
			"debug",
			`level: debug`,
			LevelDebug, "", "", emptyCustomLevels,
			false,
		},
		{
			"info",
			`level: info`,
			LevelInfo, "", "", emptyCustomLevels,
			false,
		},
		{
			"warn",
			`level: warn`,
			LevelWarn, "", "", emptyCustomLevels,
			false,
		},
		{
			"warning",
			`level: warning`,
			LevelWarn, "", "", emptyCustomLevels,
			false,
		},
		{
			"error",
			`level: error`,
			LevelError, "", "", emptyCustomLevels,
			false,
		},
		{
			"stdout",
			`output: stdout`,
			0, StdOut, "", emptyCustomLevels,
			false,
		},
		{
			"syslog",
			`output: syslog`,
			0, SysLog, "", emptyCustomLevels,
			false,
		},
		{
			"json",
			`format: json`,
			0, "", FormatJSON, emptyCustomLevels,
			false,
		},
		{
			"text",
			`format: text`,
			0, "", FormatText, emptyCustomLevels,
			false,
		},
		{
			"invalid_yaml_error",
			`not a yaml`,
			DefaultLevel, DefaultOutput, DefaultFormat, emptyCustomLevels,
			true,
		},
		{
			"unknown_level_error",
			`level: pain`,
			DefaultLevel, DefaultOutput, DefaultFormat, emptyCustomLevels,
			true,
		},
		{
			"valid_custom_logger_levels",
			`
level: info
format: text
output: stdout
customLevels: 
  traceLogger: trace
  debugLogger: debug
  infoLogger: info
  warnLogger: warn
  errorLogger: error
`,
			DefaultLevel, DefaultOutput, DefaultFormat,
			map[string]Level{
				"traceLogger": LevelTrace,
				"debugLogger": LevelDebug,
				"infoLogger":  LevelInfo,
				"warnLogger":  LevelWarn,
				"errorLogger": LevelError,
			},
			false,
		},
		{
			"valid_custom_logger_levels_short",
			`
level: warn
format: text
output: stdout
customLevels: 
  traceLogger: t
  debugLogger: d
  infoLogger: i
  warnLogger: w
  errorLogger: e
`,
			LevelWarn, DefaultOutput, DefaultFormat,
			map[string]Level{
				"traceLogger": LevelTrace,
				"debugLogger": LevelDebug,
				"infoLogger":  LevelInfo,
				"warnLogger":  LevelWarn,
				"errorLogger": LevelError,
			},
			false,
		},
		{
			"invalid_custom_logger_levels",
			`
customLevels: 
  traceLogger: trace us
  debugLogger: or debug us
  infoLogger: inform us
  warnLogger: warn us
  errorLogger: judge on error us
`,
			DefaultLevel, DefaultOutput, DefaultFormat,
			map[string]Level{
				"traceLogger": LevelTrace,
				"debugLogger": LevelDebug,
				"infoLogger":  LevelInfo,
				"warnLogger":  LevelWarn,
				"errorLogger": LevelError,
			},
			true,
		},
	} {
		tt := tt // pin tt
		t.Run(tt.name, func(t *testing.T) {
			tt := tt // pin tt
			config := new(Config)
			err := yaml.Unmarshal([]byte(tt.bytes), config)
			require.Equalf(t, err != nil, tt.wantErr, "UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
			if err != nil {
				return
			}

			require.Equal(t, tt.level, config.Level)
			require.Equal(t, tt.format, config.Format)
			require.Equal(t, tt.out, config.Output)

			for k, v := range tt.NamedLevels {
				require.Contains(t, config.CustomLevels, k)
				require.Equalf(t,
					config.CustomLevels[k], v,
					"expected custom level %v=%v not %v",
					k, config.CustomLevels[k], v)
			}

			for k, v := range config.CustomLevels {
				require.Contains(t, tt.NamedLevels, k)
				require.Equal(t, tt.NamedLevels[k], v)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	emptyCustomLevels := make(map[string]Level)
	for _, tt := range []struct { //nolint:paralleltest
		name        string
		level       Level
		out         Target
		format      Format
		NamedLevels map[string]Level
		wantErr     bool
	}{
		{
			"unknown_level_error",
			-11, DefaultOutput, DefaultFormat, emptyCustomLevels,
			true,
		},
		{
			"empty_path_error",
			DefaultLevel, "", DefaultFormat, emptyCustomLevels,
			true,
		},
		{
			"unknown_format_error",
			DefaultLevel, DefaultOutput, "some format", emptyCustomLevels,
			true,
		},
	} {
		tt := tt // pin tt
		t.Run(tt.name, func(t *testing.T) {
			tt := tt // pin tt
			config := new(Config)
			config.Level = tt.level
			config.Output = tt.out
			config.Format = tt.format
			err := config.Validate()
			require.Equalf(t, err != nil, tt.wantErr, "Validate error = %v, wantErr %v", err, tt.wantErr)
			if err != nil {
				return
			}
		})
	}
}

func TestConfig_MarshalYAML(t *testing.T) {
	config := new(Config)
	config.Output = StdOut
	config.Format = FormatText

	for _, tt := range []struct {
		level  Level
		expect string
	}{
		{LevelTrace, "trace"},
		{LevelDebug, "debug"},
		{LevelInfo, "info"},
		{LevelWarn, "warn"},
		{LevelError, "error"},
		{LevelPanic, "panic"},
		{LevelFatal, "fatal"},
	} {
		t.Run("level_"+tt.level.String(), func(t *testing.T) {
			config.Level = tt.level
			got, err := yaml.Marshal(config)
			require.NoError(t, err)
			require.Truef(t, strings.Contains(string(got), "level: "+tt.expect), "yaml:\n%v", string(got))
		})
	}
}

func TestConfig_String(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   string
	}{
		{"debug_text_stdout",
			Config{Level: LevelDebug, Format: FormatText, Output: StdOut},
			"level='debug',format='text',output='stdout',customLevels={}"},
		{"error_json_syslog",
			Config{Level: LevelError, Format: FormatJSON, Output: SysLog},
			"level='error',format='json',output='syslog',customLevels={}"},
		{"default_json_syslog_plus_custom",
			Config{
				Level:        DefaultLevel,
				Format:       FormatText,
				Output:       "/var/log/my.log",
				CustomLevels: map[string]Level{"myLogger": LevelTrace}},
			"level='info',format='text',output='/var/log/my.log',customLevels={myLogger='trace'}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.config.String())
		})
	}
}
