package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	httpErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_errors_total",
			Help: "Total number of HTTP error responses",
		},
		[]string{"method", "path", "status"},
	)
)

func RegisterHTTPMetrics() {
	prometheus.MustRegister(
		httpRequestCount,
		httpRequestDuration,
		httpErrorsTotal,
	)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()

		routePattern := chi.RouteContext(r.Context()).RoutePattern()

		httpRequestDuration.
			WithLabelValues(r.Method, routePattern).
			Observe(duration)

		httpRequestCount.
			WithLabelValues(
				r.Method,
				routePattern,
				strconv.Itoa(rw.status),
			).
			Inc()

		if rw.status >= 400 {
			httpErrorsTotal.
				WithLabelValues(
					r.Method,
					routePattern,
					strconv.Itoa(rw.status),
				).
				Inc()
		}
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
