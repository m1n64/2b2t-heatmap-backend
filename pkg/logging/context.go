package logging

import (
	"context"
	"go.uber.org/zap"
	"tbtt-heatmaps-service/pkg/tracing"
)

func FromContext(ctx context.Context, base *zap.Logger) *zap.Logger {
	v := ctx.Value(tracing.TraceParentKey{})
	if v == nil {
		return base
	}

	traceparent, ok := v.(string)
	if !ok {
		return base
	}

	return base.With(
		zap.String("trace_id", traceparent),
	)
}
