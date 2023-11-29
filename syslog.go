package logging

import (
	"fmt"
	"net"
	"sync"
)

// SyslogWriter describes connection sink for syslog.
type SyslogWriter struct {
	network string
	raddr   string
	conn    net.Conn
	mu      sync.Mutex
}

// NewSyslogWriter returns a new conn sink for syslog.
func NewSyslogWriter(network, raddr string) (s *SyslogWriter, err error) {
	s = &SyslogWriter{
		network: network,
		raddr:   raddr,
		mu:      sync.Mutex{},
	}

	if err = s.connect(); err != nil {
		return nil, fmt.Errorf("%w: connect syslog socket: %v", Error, err)
	}

	return s, nil
}

// connect makes a connection to the syslog server.
// If syslog server already connected new connection will be created.
func (s *SyslogWriter) connect() (err error) {
	var conn net.Conn

	s.mu.Lock()
	conn = s.conn
	s.mu.Unlock()

	if conn != nil {
		s.Disconnect()
	}

	if conn, err = net.Dial(s.network, s.raddr); err != nil {
		return fmt.Errorf("%w: connect syslog: %v", Error, err)
	}

	s.mu.Lock()
	s.conn = conn
	s.mu.Unlock()

	return nil
}

// Write writes to syslog. If connection is not ready it tries to connect first.
// Returns error if write failed or connection can not be established.
func (s *SyslogWriter) Write(p []byte) (n int, err error) {
	if s.conn == nil {
		if err = s.connect(); err != nil {
			return 0, err
		}
	}

	return s.conn.Write(p)
}

// Sync implements zapcore.WriteSyncer interface.
func (s *SyslogWriter) Sync() error {
	return nil
}

// Disconnect disconnects from socket
func (s *SyslogWriter) Disconnect() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn == nil {
		return
	}

	_ = s.conn.Close() // ignore err from close, it makes sense to continue anyway
	s.conn = nil
}
