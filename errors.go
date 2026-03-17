package inouesdk

import "fmt"

// SdkError represents an error returned by the Inoue AI API.
// It contains the HTTP status code, an application-level error code,
// a human-readable message, the trace ID for the request, and optional details.
type SdkError struct {
	Code    string
	Message string
	Status  int
	TraceID string
	Details map[string]interface{}
}

// Error returns a formatted string representation of the SDK error,
// including the raw response body from details for debugging.
func (e *SdkError) Error() string {
	var base string
	if e.TraceID != "" {
		base = fmt.Sprintf("inoue-sdk: %s (code=%s, status=%d, trace=%s)", e.Message, e.Code, e.Status, e.TraceID)
	} else {
		base = fmt.Sprintf("inoue-sdk: %s (code=%s, status=%d)", e.Message, e.Code, e.Status)
	}
	if raw, ok := e.Details["raw_body"]; ok {
		return fmt.Sprintf("%s body=%s", base, raw)
	}
	return base
}
