package middlewares

import (
	"github.com/alexcesaro/statsd"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"sync/atomic"
	"tbtt-heatmaps-service/pkg/logging"
	"tbtt-heatmaps-service/pkg/tracing"
	"time"
)

func AccessLogMiddleware(base *zap.Logger, stats *statsd.Client) gin.HandlerFunc {
	var inflight int64

	return func(c *gin.Context) {
		start := time.Now()

		atomic.AddInt64(&inflight, 1)
		stats.Gauge("http.inflight", inflight)

		c.Next()

		log := logging.FromContext(
			c.Request.Context(),
			base,
		)

		statusCode := c.Writer.Status()

		var loggerFn func(msg string, fields ...zap.Field)
		switch {
		case statusCode >= 500:
			loggerFn = log.Error
		case statusCode >= 400:
			loggerFn = log.Warn
		default:
			loggerFn = log.Info
		}

		tp, _ := c.Request.Context().
			Value(tracing.TraceParentKey{}).(string)

		atomic.AddInt64(&inflight, -1)
		stats.Gauge("http.inflight", inflight)

		duration := time.Since(start)
		path := c.Request.URL.Path

		message := logging.RequestMessage(
			c.Request.Method,
			path,
			tp,
			"request completed, time="+duration.String(),
		)

		fields := make([]zap.Field, 0, 6)
		fields = append(fields,
			zap.Int("status", statusCode),
			zap.Duration("duration", duration),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("client_ip", c.ClientIP()),
			zap.String("trace_id", tp),
		)

		stats.Increment("http.requests.total")
		stats.Gauge("debug.duration_ms", float64(duration.Milliseconds()))
		stats.Gauge("debug.duration_micros", float64(duration.Microseconds()))

		stats.Increment("http.responses.total")
		stats.Increment(
			"http.responses.status." + strconv.Itoa(statusCode),
		)

		if statusCode >= 400 {
			stats.Increment("http.errors.total")
			fields = append(fields, zap.String("user_agent", c.Request.UserAgent()))
		}

		loggerFn(message, fields...)
	}
}
