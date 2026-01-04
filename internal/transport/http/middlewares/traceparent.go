package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"tbtt-heatmaps-service/pkg/tracing"
)

func TraceparentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		incoming := c.GetHeader("traceparent")

		traceparent := tracing.GetOrCreateTrace(incoming)

		ctx := context.WithValue(
			c.Request.Context(),
			tracing.TraceParentKey{},
			traceparent,
		)
		c.Request = c.Request.WithContext(ctx)

		c.Header("traceparent", traceparent)
		c.Next()
	}
}
