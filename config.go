package logging

/* Config implements service logging configuration structure data.

It can be loaded from YAML configuration file.
Configuration structure in YAML should looks like:

	level: [debug|info|warn|error]
	format: [text|json]
	output: [stdout|/path/to/file]

Default is output to console (stdout) and use text format on info level.

*/
import (
	"fmt"
	"strings"
)

// CurrentConfig returns pointer to stored config.
// It is definitely nil until Init() executed with valid config.
// If not nil, its CustomLevels used in NewNamedLogger and NewNamedLoggerCtx constructors.
// NOTE: returned config CustomLevels is not safe to use for both reading and writing parallel in different goroutines.
func CurrentConfig() *Config {
	mu.Lock()
	d := useConfig
	mu.Unlock()

	return d
}

// Config type implements parameters storage for logging initialization.
// NOTE: once created and passed to Init configuration is not expected to be changed.
// Use NewConfig constructor to make configuration instance as some internals should be initialized properly.
type Config struct {
	// Level provides logging level for application logging.
	// Default LevelInfo.
	Level Level `yaml:"level"`
	// Format defines logging format. Either JSON or TEXT.
	// Default JSON.
	Format Format `yaml:"format,omitempty"`
	// Output defines logging path.
	// Can be absolute path to log into file or substituted with one of predefined StdOut or SysLog constants.
	// Default StdOut.
	Output Target `yaml:"output,omitempty"`
	// CustomLevels allows separate level definitions for named loggers using their names.
	// Each specified level expected to be less verbose than global level defined in Config.Level attribute.
	CustomLevels map[string]Level `yaml:"customLevels,omitempty"`

	// contextExtractors registers extract context-provided data as fields
	contextExtractors map[Key]ContextExtractorFun
}

// NewConfig creates new logging configuration with defaults set.
func NewConfig() *Config {
	return &Config{
		Level:             DefaultLevel,
		Output:            DefaultOutput,
		Format:            DefaultFormat,
		CustomLevels:      make(map[string]Level),
		contextExtractors: make(map[Key]ContextExtractorFun),
	}
}

// Validate checks logging config valid. Returns error if any misconfiguration detected.
func (config Config) Validate() error {
	switch {
	case config.Output == "":
		return fmt.Errorf("%w: output empty, expected '%v', '%v' or absolute file path", Error, StdOut, SysLog)
	}

	if err := config.Format.Validate(); err != nil {
		return err
	}

	switch config.Level {
	case LevelTrace, LevelDebug, LevelInfo, LevelWarn, LevelError, LevelPanic, LevelFatal:
	default:
		return fmt.Errorf("%w: unknown level: `%v`", Error, config.Level)
	}

	for k, lvl := range config.CustomLevels {
		switch lvl {
		case LevelTrace, LevelDebug, LevelInfo, LevelWarn, LevelError, LevelPanic, LevelFatal:
		default:
			return fmt.Errorf("%w: unknown level %v: `%v`", Error, k, lvl)
		}
	}

	return nil
}

// String returns string representation of logging config.
func (config Config) String() string {

	wrapValue := func(i string) string { return "'" + i + "'" }
	add := func(n string, v string) string { return n + "=" + wrapValue(v) }

	if config.CustomLevels == nil { // can be nil if created with new(Config)
		config.CustomLevels = make(map[string]Level)
	}

	customLevels := make([]string, len(config.CustomLevels))
	idx := 0
	for name, lvl := range config.CustomLevels {
		customLevels[idx] = add(name, lvl.String())
		idx++
	}

	return strings.Join([]string{
		add("level", config.Level.String()),
		add("format", config.Format.String()),
		add("output", config.Output.String()),
		"customLevels={" + strings.Join(customLevels, ",") + "}",
	}, ",")
}

// levelForNamed returns configured level for named logger.
// Returns first specified layer if no custom layer set for logger or config global layer.
func (config Config) levelForNamed(name string, levels ...Level) Level {
	var (
		customLevel Level
		ok          bool
	)

	customLevel, ok = config.CustomLevels[name]
	switch {
	case !ok && len(levels) == 0:
		customLevel = config.Level
	case !ok:
		customLevel = levels[0]
	}

	return customLevel
}

type Option func(*Config)

// Apply applies specified configuration options.
func (config *Config) Apply(opts ...Option) {
	for _, o := range opts {
		o(config)
	}
}
