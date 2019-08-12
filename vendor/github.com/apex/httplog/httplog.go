package httplog

import (
	"net/http"
	"time"

	"github.com/apex/log"
)

// New middleware wrapping `h`.
func New(h http.Handler) *Logger {
	return &Logger{Handler: h}
}

// Logger middleware wrapping Handler.
type Logger struct {
	http.Handler
}

// wrapper to capture status.
type wrapper struct {
	http.ResponseWriter
	http.Flusher
	http.CloseNotifier

	written int
	status  int
}

// WriteHeader wrapper to capture status code.
func (w *wrapper) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// Write wrapper to capture response size.
func (w *wrapper) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.written += n
	return n, err
}

// Flush implementation.
func (w *wrapper) Flush() {
	if w.Flusher != nil {
		w.Flusher.Flush()
	}
}

// ServeHTTP implementation.
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	res := &wrapper{
		ResponseWriter: w,
		written:        0,
		status:         200,
	}

	if f, ok := w.(http.Flusher); ok {
		res.Flusher = f
	}

	if c, ok := w.(http.CloseNotifier); ok {
		res.CloseNotifier = c
	}

	ctx := log.WithFields(log.Fields{
		"url":        r.RequestURI,
		"method":     r.Method,
		"remoteAddr": r.RemoteAddr,
	})

	ctx.Info("request")
	l.Handler.ServeHTTP(res, r)

	ctx = ctx.WithFields(log.Fields{
		"status":   res.status,
		"size":     res.written,
		"duration": time.Since(start),
	})

	switch {
	case res.status >= 500:
		ctx.Error("response")
	case res.status >= 400:
		ctx.Warn("response")
	default:
		ctx.Info("response")
	}
}
