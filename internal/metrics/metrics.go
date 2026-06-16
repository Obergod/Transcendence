package metrics

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
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

	visitorIPs sync.Map
	sessions   = make(map[string]time.Time)
	mu         sync.Mutex
)

func init() {
	prometheus.MustRegister(HttpRequests)
	prometheus.MustRegister(UniqueVisitors)
	prometheus.MustRegister(ActiveUsers)
	prometheus.MustRegister(VisitDuration)
}

func TrackRequests() gin.HandlerFunc {
	return func(c *gin.Context) {
		HttpRequests.WithLabelValues(c.FullPath()).Inc()
		c.Next()
	}
}

func TrackUniqueVisitors() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if _, exists := visitorIPs.Load(ip); exists {
			c.Next()
			return
		}

		UniqueVisitors.Inc()
		visitorIPs.Store(ip, struct{}{})

		c.SetCookie(
			"session_id",
			uuid.New().String(),
			86400,
			"/",
			"",
			false,
			true,
		)

		c.Next()
	}
}

func TrackActiveUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil || sessionID == "" {
			c.Next()
			return
		}

		mu.Lock()
		sessions[sessionID] = time.Now().Add(24 * time.Hour)
		mu.Unlock()

		ActiveUsers.Inc()
		defer ActiveUsers.Dec()

		c.Next()
	}
}

func TrackVisitDuration() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()

		c.Next()

		duration := time.Since(start).Seconds()
		VisitDuration.WithLabelValues(path).Observe(duration)
	}
}