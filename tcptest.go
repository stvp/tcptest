package tcptest

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// A Server is a TCP server listening on a system-chosen port on the local
// loopback interface, for use in end-to-end TCP tests.
type Server struct {
	listener net.Listener
	lines    []string
	wg       sync.WaitGroup
	sync.Mutex
}

// NewServer opens and returns a new Server. The caller should call Close when
// finished to shut it down.
func NewServer() (server *Server, err error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return
	}

	server = &Server{
		listener: listener,
		lines:    []string{},
		wg:       sync.WaitGroup{},
	}
	go server.run()

	return
}

// Address returns the host:port string for this Server.
func (s *Server) Address() string {
	return s.listener.Addr().String()
}

// WaitForLines blocks until an expected number of lines have been received by
// the server. If the timeout expires, an error is returned.
func (s *Server) WaitForLines(count int, timeout time.Duration) error {
	deadline := time.After(timeout)

	for {
		select {
		case <-deadline:
			return fmt.Errorf("tcptest: expected %d lines but only received %d", count, len(s.lines))
		default:
			if len(s.lines) >= count {
				return nil
			}
		}
	}
}

// Lines returns all received lines.
func (s *Server) Lines() []string {
	lines := make([]string, len(s.lines))
	for i, line := range s.lines {
		lines[i] = line
	}
	return lines
}

// Received returns true if the given string has been received.
func (s *Server) Received(expect string) bool {
	for _, line := range s.lines {
		if strings.Contains(line, expect) {
			return true
		}
	}
	return false
}

// Received returns true if the given full line (minus terminating newline) has
// been received.
func (s *Server) ReceivedLine(expect string) bool {
	for _, got := range s.lines {
		if got == expect {
			return true
		}
	}
	return false
}

// Close waits for all client connections to close themselves and stops the
// Server.
func (s *Server) Close() {
	s.wg.Wait()
	s.listener.Close()
}

func (s *Server) run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		s.wg.Add(1)
		go func() {
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				s.Lock()
				s.lines = append(s.lines, scanner.Text())
				s.Unlock()
			}
			s.wg.Done()
		}()
	}
}
