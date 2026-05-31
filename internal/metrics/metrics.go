package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/google/uuid"
    "time"
	"sync"
)

var (
    HttpRequests = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"path"},
    )

    UniqueVisitors = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "unique_visitors_total",
            Help: "Total number of unique visitors to the site",
        },
    )

	ActiveUsers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_users",
			Help: "Number of currently active users",
		},
	)
	VisitDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "visit_duration_seconds",
			Help:    "Duration of user visits in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
	sessions = make(map[string]time.Time)
    mu     sync.Mutex
)

func init() {
	prometheus.MustRegister(HttpRequests)
	prometheus.MustRegister(UniqueVisitors)
	prometheus.MustRegister(ActiveUsers)
    prometheus.MustRegister(VisitDuration)
}

func TrackRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HttpRequests.WithLabelValues(r.URL.Path).Inc()
		next.ServeHTTP(w, r)
	})
}

func TrackUniqueVisitors(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr
        if _, exists := visitorIPs.Load(ip); exists {
            next.ServeHTTP(w, r)
            return
        }
        UniqueVisitors.Inc()
        visitorIPs.Store(ip, struct{}{})

        http.SetCookie(w, &http.Cookie{
            Name:    "session_id",
            Value:    uuid.New().String(),
            Expires:   time.Now().Add(24 * time.Hour),
            HttpOnly: true,
        })
        next.ServeHTTP(w, r)
    })
}

var (
    visitorIPs sync.Map
)

func TrackActiveUsers(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        sessionCookie, err := r.Cookie("session_id")
        if err != nil || sessionCookie.Value == "" {
            next.ServeHTTP(w, r)
            return
        }

        mu.Lock()
        sessions[sessionCookie.Value] = time.Now().Add(24 * time.Hour)
        ActiveUsers.Inc()
        mu.Unlock()
        defer ActiveUsers.Dec()

        next.ServeHTTP(w, r)
    })
}

func TrackVisitDuration(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()


        rec := &responseRecorder{
            ResponseWriter: w,
            start:          start,
            path:           r.URL.Path,
        }

        next.ServeHTTP(rec, r)
    })
}

type responseRecorder struct {
    http.ResponseWriter
    start time.Time
    path  string
}

func (rec *responseRecorder) WriteHeader(code int) {
    rec.ResponseWriter.WriteHeader(code)
}

func (rec *responseRecorder) Write(data []byte) (int, error) {
    duration := time.Since(rec.start).Seconds()
    VisitDuration.WithLabelValues(rec.path).Observe(duration)
    return rec.ResponseWriter.Write(data)
}