package metrics

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

// Middleware returns a HTTP handler with the addition of status codes for responses.
func Middleware(httpDuration *prometheus.HistogramVec) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wi := &responseWriterInterceptor{
				StatusCode:     http.StatusOK,
				ResponseWriter: w,
			}

			route := mux.CurrentRoute(r)
			path, _ := route.GetPathTemplate()
			start := time.Now()
			next.ServeHTTP(wi, r)
			httpDuration.WithLabelValues(path, strconv.Itoa(wi.StatusCode)).Observe(time.Since(start).Seconds())
		})
	}
}

// responseWriterInterceptor is a simple wrapper to intercept set data on a
// ResponseWriter.
type responseWriterInterceptor struct {
	http.ResponseWriter
	StatusCode int
}

func (w *responseWriterInterceptor) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterInterceptor) Write(p []byte) (int, error) {
	return w.ResponseWriter.Write(p)
}

func (w *responseWriterInterceptor) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, ErrTypeAssertion
	}

	return h.Hijack()
}

func (w *responseWriterInterceptor) Flush() {
	f, ok := w.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}

	f.Flush()
}

var (
	ErrTypeAssertion = errors.New("type assertion failed http.ResponseWriter not a http.Hijacker")
)
