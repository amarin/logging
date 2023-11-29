package logging

import (
	"errors"
	"time"
)

const (
	// internal package name to represent it in errors.
	moduleName = "logging"

	// TimestampFormatConsole defines console timestamp format.
	TimestampFormatConsole = "2006-01-02 15:04:05.000000"

	// TimestampFormatJSON defines JSON timestamp format.
	TimestampFormatJSON = time.RFC3339Nano
)

// Error identifies logging module errors.
var Error = errors.New(moduleName)
