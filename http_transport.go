package ms

import (
	"errors"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// HTTPTransport
type HTTPTransport struct {
	Logger  *log.Entry
	Addr    string
	Handler *http.Handler
	Closers []io.Closer
}

// GRPCTransportOptions hold options for server
type HTTPTransportOptions struct {
	Addr    string
	Handler *http.Handler
}

// NewServer creates a new http transport
func NewHTTPTransport(opts *HTTPTransportOptions) (*HTTPTransport, error) {
	logger := log.WithField("module", "HTTPTransport")

	if opts == nil {
		return nil, errors.New("missing HTTPTransportOptions")
	}

	if opts.Handler == nil {
		return nil, errors.New("missing handler")
	}

	return &HTTPTransport{
		Logger:  logger,
		Addr:    opts.Addr,
		Handler: opts.Handler,
	}, nil
}

// ServeMux
func (t *HTTPTransport) ServeHandler() *http.Handler {
	return t.Handler
}

// Add
func (t *HTTPTransport) Add(c io.Closer) {
	t.Closers = append(t.Closers, c)
}

// Run
func (t *HTTPTransport) Run() error {
	l := t.Logger.WithField("function", "Run")

	l.WithField("addr", t.Addr).Info("listening")

	return http.ListenAndServe(t.Addr, *t.Handler)
}

// Close
func (t *HTTPTransport) Close() error {
	l := t.Logger.WithField("function", "Close")

	for _, v := range t.Closers {
		if err := v.Close(); err != nil {
			l.WithField("err", err).Error()
		}
	}

	return nil
}
