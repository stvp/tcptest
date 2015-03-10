package tcptest

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

type Server struct {
	listener net.Listener
	lines    []string
	wg       sync.WaitGroup
	sync.Mutex
}

// NewServer opens and returns a new Server. The caller should
// call Close when finished to shut it down.
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

func (s *Server) Received(expect string) bool {
	for _, got := range s.lines {
		if got == expect {
			return true
		}
	}
	return false
}

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
