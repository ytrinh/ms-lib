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
    Mux     *http.ServeMux
    Closers []io.Closer
}
// GRPCTransportOptions hold options for server
type HTTPTransportOptions struct {
	Addr string
	Mux  *http.ServeMux
}
// NewServer creates a new http transport
func NewHTTPTransport(opts *HTTPTransportOptions) (*HTTPTransport, error) {
	logger := log.WithField("module", "HTTPTransport")
	
	if opts == nil {
		return nil, errors.New("missing HTTPTransportOptions")
	}

	if opts.Mux == nil {
		opts.Mux = http.NewServeMux()
	}

    return &HTTPTransport{
        Logger: logger,
        Addr:   opts.Addr,
        Mux:    opts.Mux,
    }, nil
}

// ServeMux
func (t *HTTPTransport) ServeMux() *http.ServeMux {
    return t.Mux
}

// Add
func (t *HTTPTransport) Add(c io.Closer) {
    t.Closers = append(t.Closers, c)
}

// Run
func (t *HTTPTransport) Run() error {
    l := t.Logger.WithField("function", "Run")

    l.WithField("addr", t.Addr).Info("listening")

    return http.ListenAndServe(t.Addr, t.Mux)
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
