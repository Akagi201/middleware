// Package logrus is a simple net/http middleware for logging
package logrus

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/Sirupsen/logrus"
)

// Middleware is a middleware handler that logs the request as it goes in and the response as it goes out.
type Middleware struct {
	// Logger is the log.Logger instance used to log messages with the Logger middleware
	Logger *logrus.Logger
	// Name is the name of the application as recorded in latency metrics
	Name string

	Before func(*logrus.Entry, *http.Request, string) *logrus.Entry
	After  func(*logrus.Entry, ResponseWrapper, time.Duration, string) *logrus.Entry

	// Exclude URLs from logging
	excludeURLs []string

	logStarting bool
	h           http.Handler
}

// ResponseWrapper wrapper to capture status.
type ResponseWrapper struct {
	http.ResponseWriter
	written int
	status  int
}

// WriteHeader capture status.
func (w *ResponseWrapper) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// Write capture written bytes.
func (w *ResponseWrapper) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.written += n
	return n, err
}

// New creates a new default middleware func.
func New() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		log := logrus.New()
		log.Level = logrus.InfoLevel
		log.Formatter = &logrus.TextFormatter{}

		return &Middleware{
			Logger:      log,
			Name:        "web",
			logStarting: true,
		}
	}
}

// NewLogger creates a new middleware func which writes to a given logrus logger.
func NewLogger(logger *logrus.Logger, name string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &Middleware{
			Logger:      logger,
			Name:        name,
			logStarting: true,
		}
	}
}

// SetLogStarting accepts a bool to control the logging of "started handling
// request" prior to passing to the next middleware
func (m *Middleware) SetLogStarting(v bool) {
	m.logStarting = v
}

// ExcludeURL adds a new URL u to be ignored during logging. The URL u is parsed, hence the returned error
func (m *Middleware) ExcludeURL(u string) error {
	if _, err := url.Parse(u); err != nil {
		return err
	}
	m.excludeURLs = append(m.excludeURLs, u)
	return nil
}

// ExcludedURLs returns the list of excluded URLs for this middleware
func (m *Middleware) ExcludedURLs() []string {
	return m.excludeURLs
}

// ServeHTTP calls the "real" handler and logs using the logger
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Try to get the real IP
	remoteAddr := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		remoteAddr = realIP
	}

	entry := logrus.NewEntry(m.Logger)

	if reqID := r.Header.Get("X-Request-Id"); reqID != "" {
		entry = entry.WithField("request_id", reqID)
	}

	entry = m.Before(entry, r, remoteAddr)

	if m.logStarting {
		entry.Info("started handling request")
	}

	rw := &ResponseWrapper{w, 0, 200}

	m.h.ServeHTTP(rw, r)

	latency := time.Since(start)

	m.After(entry, w, latency, m.Name).Info("completed handling request")
}

// BeforeFunc is the func type used to modify or replace the *logrus.Entry prior
// to calling the next func in the middleware chain
type BeforeFunc func(*logrus.Entry, *http.Request, string) *logrus.Entry

// AfterFunc is the func type used to modify or replace the *logrus.Entry after
// calling the next func in the middleware chain
type AfterFunc func(*logrus.Entry, http.ResponseWriter, time.Duration, string) *logrus.Entry

// DefaultBefore is the default func assigned to *Middleware.Before
func DefaultBefore(entry *logrus.Entry, req *http.Request, remoteAddr string) *logrus.Entry {
	return entry.WithFields(logrus.Fields{
		"request": req.RequestURI,
		"method":  req.Method,
		"remote":  remoteAddr,
	})
}

// DefaultAfter is the default func assigned to *Middleware.After
func DefaultAfter(entry *logrus.Entry, rw ResponseWrapper, latency time.Duration, name string) *logrus.Entry {
	return entry.WithFields(logrus.Fields{
		"status":      rw.status,
		"text_status": http.StatusText(rw.status),
		"writen":      humanize.Bytes(uint64(rw.written)),
		"took":        latency,
		fmt.Sprintf("measure#%s.latency", name): latency.Nanoseconds(),
	})
}
