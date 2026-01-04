package tracing

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

type TraceParentKey struct{}

func randomHex(bytes int) string {
	b := make([]byte, bytes)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func NewTrace() (traceID, spanID, traceparent string) {
	traceID = randomHex(16)
	spanID = randomHex(8)
	traceparent = "00-" + traceID + "-" + spanID + "-01"
	return traceID, spanID, traceparent
}

func GetOrCreateTrace(header string) (traceparent string) {
	if header != "" {
		return header
	}

	_, _, traceparent = NewTrace()

	return traceparent
}

func TraceparentFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(TraceParentKey{})
	if v == nil {
		return "", false
	}
	tp, ok := v.(string)
	return tp, ok
}
