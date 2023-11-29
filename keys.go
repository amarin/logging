package logging

const (
	// KeyTimestamp defines logging key for timestamp registering.
	KeyTimestamp Key = "ts"
	// KeyLevel defines logging key for message logging level registering.
	KeyLevel Key = "level"
	// KeyMessage defines logging key for registering message itself.
	KeyMessage Key = "msg"
	// KeyLogger defines logging key for provider, module or subsystem name registering.
	KeyLogger Key = "logger"
	// KeyCaller defines logging key for function caller registering.
	KeyCaller Key = "caller"
	// KeyStackTrace defines logging key for provider name registering.
	KeyStackTrace Key = "stacktrace"
)

// Key defines logging key name string interface.
type Key string

// String returns string representation of Key.
func (k Key) String() string {
	return string(k)
}

// Keys simply wraps map[Key]interface{} to use in Logger.WithKeys.
// It seems that logging.Keys is shorter than map[Key]interface{}.
type Keys map[Key]any
