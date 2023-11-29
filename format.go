package logging

import "fmt"

const (
	// FormatText defines constant value for format logging lines as text.
	FormatText Format = "text"

	// FormatJSON defines constant value for format logging lines as JSON.
	FormatJSON Format = "json"

	// DefaultFormat defines default formatting if omitted.
	DefaultFormat = FormatText
)

// Format defines output format. It should be either FormatText or FormatJSON
type Format string

// String returns string representation of Format. Implements fmt.Stringer.
func (f Format) String() string {
	return string(f)
}

// Validate returns error if format value is not valid.
func (f Format) Validate() error {
	switch f {
	case FormatText:
		return nil
	case FormatJSON:
		return nil
	default:
		return fmt.Errorf("%w: unexpected format `%v`, want `%s` or `%s", Error, f, FormatText, FormatJSON)
	}
}

// WithFormat adds specified logging output format to configuration.
func WithFormat(format Format) Option {
	return func(config *Config) {
		config.Format = format
	}
}
