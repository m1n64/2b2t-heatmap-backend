package logging

func RequestMessage(method, path, traceID, msg string) string {
	return "[" + method + "] " +
		"[" + path + "] " +
		"[" + traceID + "] " +
		msg
}
