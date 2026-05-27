package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var HttpRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests",
	},
	[]string{"path"},
)

func init() {
	prometheus.MustRegister(HttpRequests)
}

func TrackRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HttpRequests.WithLabelValues(r.URL.Path).Inc()
		next.ServeHTTP(w, r)
	})
}