package logging

import (
	"fmt"
)

const (
	// StdOut is a standard way to disable file logging. Use target: stdout to log to console stdout.
	StdOut Target = "stdout"

	// StdErr is a standard way to disable file logging. Use target: stderr to log to console std error.
	StdErr Target = "stderr"

	// SysLog is a predefined output constant to route logging messages into syslog.
	SysLog Target = "syslog"

	// DefaultOutput defines default output if omitted.
	DefaultOutput = StdOut
)

// Target identifies logging output target.
type Target string

// String returns string representation of Target.
func (output Target) String() string {
	return string(output)
}

// Validate returns error if Target settings is not valid.
func (output Target) Validate() error {
	switch output {
	case "":
		return fmt.Errorf("%w: output empty, want '%v', '%v', '%v' or file path", Error, StdOut, StdErr, SysLog)
	case StdOut, StdErr, SysLog:
		return nil
	}
	// expected to be path, it may be there should be path check?
	return nil
}

func WithTarget(target Target) Option {
	return func(config *Config) {
		config.Output = target
	}
}
