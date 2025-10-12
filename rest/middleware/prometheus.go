package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_activity_tracking_api_requests_total",
			Help: "Total number of requests processed by the user-activity-tracking-api.",
		},
		[]string{"path", "status"},
	)

	ErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_activity_tracking_api_requests_errors_total",
			Help: "Total number of error requests processed by the user-activity-tracking-api.",
		},
		[]string{"path", "status"},
	)
)

func init() {
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(ErrorCount)
}

func TrackMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if path == "/metrics" {
			return
		}

		c.Next()
		status := c.Writer.Status()
		RequestCount.WithLabelValues(path, http.StatusText(status)).Inc()

		if status >= 400 {
			ErrorCount.WithLabelValues(path, http.StatusText(status)).Inc()
		}
	}
}
