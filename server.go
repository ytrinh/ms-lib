package ms

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// Runner interface
type Runner interface {
	Run() error
}

// Server property
type Server struct {
	logger  *log.Entry
	runners []Runner
	options *ServerOptions
}

// ServerOptions hold options for server
type ServerOptions struct {
	CloseTimeoutSeconds int64
}

// ServerOptionDefault
var ServerOptionDefault = &ServerOptions{
	CloseTimeoutSeconds: 10,
}

// NewServer creates a new Server
func NewServer(opts *ServerOptions) (*Server, error) {
	logger := log.WithField("module", "Server")

	if opts == nil {
		opts = ServerOptionDefault
	}

	return &Server{
		logger:  logger,
		runners: nil,
		options: opts,
	}, nil
}

// Add a runner
func (s *Server) Add(runner Runner) {
	if runner == nil {
		return
	}

	s.runners = append(s.runners, runner)
}

// Run start the server
func (s *Server) Run() error {
	//l := s.logger.WithField("function", "Run")
	if s.runners == nil || len(s.runners) == 0 {
		return errors.New("no runners")
	}

	// interrupt handler.
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// run all defined runners
	for _, v := range s.runners {
		go func(r Runner) {
			errc <- r.Run()
		}(v)
	}

	// wait for event
	var err = <-errc

	// cleanup
	s.Close()

	return err
}

// Close manages server and runners shutdown by calling any
// runners with Close() interface and wait return or timeout.
func (s *Server) Close() error {
	l := s.logger.WithField("function", "Close")

	// close all runners
	var wg sync.WaitGroup
	for _, v := range s.runners {
		// see if any of the runners is also a closer
		c, ok := v.(io.Closer)
		if !ok {
			continue
		}

		wg.Add(1)
		go func(c io.Closer) {
			defer wg.Done()
			if err := c.Close(); err != nil {
				l.WithField("err", err).Error()
			}
		}(c)
	}

	// wait until all closers are done then notify channel
	errc := make(chan bool)
	go func() {
		wg.Wait()
		errc <- true
	}()

	sec := s.options.CloseTimeoutSeconds
	if sec < 10 { // wait at least 10 sec
		sec = 10
	}

	// wait untill all closers are done or timeout
	select {
	case <-errc:
		break
	case <-time.After(time.Duration(sec) * time.Second):
		l.Warn("timeout reached")
	}

	return nil
}
