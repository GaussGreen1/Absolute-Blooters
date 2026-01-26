package middleware

import (
	"blooters/internal/metrics"
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK, &bytes.Buffer{}}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		start := time.Now()

		var requestBody bytes.Buffer
		if r.Body != nil {
			tee := io.TeeReader(r.Body, &requestBody)
			bodyBytes, _ := io.ReadAll(tee)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		logrus.WithFields(logrus.Fields{
			"request_id":       requestID,
			"method":           r.Method,
			"path":             r.URL.Path,
			"status":           rw.statusCode,
			"latency_ms":       duration.Milliseconds(),
			"request_preview":  truncate(requestBody.String(), 8000),
			"response_preview": truncate(rw.body.String(), 200),
		}).Info("completed request")

		// Record metrics
		metrics.HTTPRequestDuration.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(rw.statusCode)).Observe(duration.Seconds())
		metrics.HTTPRequestCount.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(rw.statusCode)).Inc()
	})
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
