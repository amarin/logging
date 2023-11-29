package logging

/* Package provides a simple output writer wrapping to prevent simultaneous writing into single writer from
multiple coroutines or different backends.

To use wrapper backend should use logging writer registry which provides wrapped writers.
To create wrapped writer backend should check if expected wrapper is not registered yet (f.e. by another wrapper)
and put desired writer to registry in order next writer receive attempt could receive wrapper writer.

See LockedWriter documentation.
*/
import (
	"fmt"
	"io"
	"os"
	"sync"
)

// The Syncer interface indicates implementation could flush buffered data to the underlying writer.
type Syncer interface {
	// Sync sends any buffered data to the underlying io.Writer.
	Sync() error
}

// A WriteSyncer is an io.Writer that can also flush any buffered data. Note
// that *os.File (and thus, os.Stderr and os.Stdout) implement WriteSyncer.
type WriteSyncer interface {
	io.Writer
	Syncer
}

var (
	// ErrClosed indicates output writer is closed.
	ErrClosed = fmt.Errorf("%w: stream closed", Error)
	writers   = newWritersRegistry()
)

type syncerWrapper struct {
	io.Writer
}

func (w syncerWrapper) Sync() error {
	return nil
}

// AddSync converts an io.Writer to a WriteSyncer. It attempts to be
// intelligent: if the concrete type of the io.Writer implements WriteSyncer,
// we'll use the existing Sync method. If it doesn't, we'll add a no-op Sync.
func wrapWriter(w io.Writer) WriteSyncer {
	switch w := w.(type) {
	case WriteSyncer:
		return w
	default:
		return &syncerWrapper{Writer: w}
	}
}

// LockedWriter implements simple write-locked operations and implements io.Writer, io.Closer and Flushed interfaces.
type LockedWriter struct {
	stream WriteSyncer
	mu     sync.Mutex
}

// Lock wraps a WriteSyncer in a mutex to make it safe for concurrent use.
func wrapSyncer(writeSyncer WriteSyncer) (wrapped *LockedWriter) {
	var ok bool

	if wrapped, ok = writeSyncer.(*LockedWriter); ok {
		// no need to layer on another lock
		return wrapped
	}
	return &LockedWriter{stream: writeSyncer, mu: sync.Mutex{}}
}

// wrap wraps specified writer into WriteSyncer when required and returns LockedWriter instance.
func wrap(writer io.Writer) (wrapped *LockedWriter) {
	return wrapSyncer(wrapWriter(writer))
}

// Write writes len(p) bytes from p to the underlying data stream.
// It returns the number of bytes written from p (0 <= n <= len(p))
// and any error encountered that caused to stop early.
func (s *LockedWriter) Write(p []byte) (n int, err error) {
	s.mu.Lock()

	if s.stream == nil {
		return 0, ErrClosed
	}

	n, err = s.stream.Write(p)
	s.mu.Unlock()

	return n, err
}

// Close closes underlying writer if implements io.Closer and sets it to nil.
// Any consequent Write calls will return ErrClosed.
// If already closed returns no error.
func (s *LockedWriter) Close() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stream == nil {
		return nil // already closed
	}

	if closer, ok := s.stream.(io.Closer); ok {
		err = closer.Close()
	}

	s.stream = nil

	return err
}

// Sync flushes underlying WriteSyncer.
func (s *LockedWriter) Sync() error {
	return s.stream.Sync()
}

// Output returns io.Writer to use in Backend instances or error if output open/create failed.
func Output(output string) (writer *LockedWriter, err error) {
	return writers.registeredOutput(output)
}

// WriterName defines simple string type to prevent coding errors dealing with WritersRegistry.
type WriterName string

// WritersRegistry implements simple map-like storage of stream writers.
type WritersRegistry struct {
	mu      sync.Mutex                   // protect underlying map
	writers map[WriterName]*LockedWriter // writes mapping itself
}

// newWritersRegistry makes a new WritersRegistry
func newWritersRegistry() *WritersRegistry {
	return &WritersRegistry{
		mu:      sync.Mutex{},
		writers: make(map[WriterName]*LockedWriter),
	}
}

// put tries to put named writer into registry. Returns error if such name already taken.
func (registry *WritersRegistry) put(name WriterName, writer io.Writer) (err error) {
	var ok bool

	registry.mu.Lock()

	defer registry.mu.Unlock()

	_, ok = registry.writers[name]

	switch {
	case writer == nil:
		return fmt.Errorf("%w: nil writer", Error)
	case len(name) == 0:
		return fmt.Errorf("%w: empty writer name", Error)
	case ok:
		return fmt.Errorf("%w: writer %s already registered", Error, name)
	}

	switch v := writer.(type) {
	case *LockedWriter:
		registry.writers[name] = v
	default:
		registry.writers[name] = wrapSyncer(wrapWriter(v))
	}

	return nil
}

// get returns a registered writer from registry or error if writer not found.
func (registry *WritersRegistry) get(name WriterName) (writer *LockedWriter, err error) {
	var ok bool

	registry.mu.Lock()

	defer registry.mu.Unlock()

	if writer, ok = registry.writers[name]; !ok {
		return nil, fmt.Errorf("%w: writer %s unknown", Error, name)
	}

	return writer, nil
}

// del deletes a registered writer from registry. If no such writer registered does nothing.
func (registry *WritersRegistry) del(name WriterName) {
	registry.mu.Lock()

	if _, ok := registry.writers[name]; ok {
		delete(registry.writers, name)
	}

	registry.mu.Unlock()
}

// registeredOutput returns a locked writer for specified output or error.
func (registry *WritersRegistry) registeredOutput(output string) (writer *LockedWriter, err error) {
	var (
		ok bool
		fh io.Writer
	)

	registry.mu.Lock()

	defer registry.mu.Unlock()

	if writer, ok = registry.writers[WriterName(output)]; ok { // already exists
		return writer, nil
	}

	switch Target(output) {
	case StdOut:
		writer = wrap(os.Stdout)
	case StdErr:
		writer = wrap(os.Stderr)
	case SysLog:
		if fh, err = NewSyslogWriter("udp", "localhost:514"); err != nil {
			return nil, fmt.Errorf("%w: set syslog output: %v", Error, err)
		}
		writer = wrap(fh)
	default:
		if fh, err = os.Create(output); err != nil {
			return nil, err
		}
		writer = wrap(fh)
	}

	registry.writers[WriterName(output)] = writer

	return writer, nil
}
