package inouesdk

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// InoueClient is the top-level client for the Inoue AI SaaS API.
// Create one with NewClient and then use the sub-API accessors (Internal, Schedule)
// to call specific endpoints.
type InoueClient struct {
	baseURL     string
	accessToken string
	httpClient  *http.Client
	Internal    *InternalAPI
	Schedule    *ScheduleAPI
}

// Option configures an InoueClient during construction.
type Option func(*InoueClient)

// WithAccessToken sets the initial Bearer token on the client.
func WithAccessToken(token string) Option {
	return func(c *InoueClient) {
		c.accessToken = token
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *InoueClient) {
		c.httpClient.Timeout = d
	}
}

// WithHTTPClient replaces the underlying http.Client entirely.
func WithHTTPClient(client *http.Client) Option {
	return func(c *InoueClient) {
		c.httpClient = client
	}
}

// NewClient creates an InoueClient pointed at the given base URL.
// Supply Options to configure authentication, timeouts, or a custom http.Client.
func NewClient(baseURL string, opts ...Option) *InoueClient {
	c := &InoueClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	c.Internal = &InternalAPI{client: c}
	c.Schedule = &ScheduleAPI{client: c}
	return c
}

// SetAccessToken updates the Bearer token used for subsequent requests.
func (c *InoueClient) SetAccessToken(token string) {
	c.accessToken = token
}

// request performs an HTTP request and decodes the JSON response into dest.
// It attaches the Authorization, Content-Type, and X-Trace-Id headers automatically.
// If the server returns a status >= 400 the method returns a *SdkError.
func (c *InoueClient) request(ctx context.Context, method, path string, body interface{}, dest interface{}, extraHeaders map[string]string) error {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("inoue-sdk: failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("inoue-sdk: failed to create request: %w", err)
	}

	traceID := generateTraceID()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Trace-Id", traceID)
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}
	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("inoue-sdk: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("inoue-sdk: failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return parseErrorResponse(respBody, resp.StatusCode, traceID, method, path)
	}

	if dest != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, dest); err != nil {
			return fmt.Errorf("inoue-sdk: failed to decode response: %w", err)
		}
	}

	return nil
}

// parseErrorResponse attempts to extract a structured error from the response body.
// If the body cannot be parsed it falls back to a generic SdkError with the raw body.
// The details map always includes the raw response for debugging.
func parseErrorResponse(body []byte, status int, traceID, method, path string) *SdkError {
	details := map[string]interface{}{
		"raw_body": string(body),
		"method":   method,
		"path":     path,
	}

	var apiResp ApiResponse
	if err := json.Unmarshal(body, &apiResp); err == nil && apiResp.Error != nil {
		if apiResp.Error.Details != nil {
			for k, v := range apiResp.Error.Details {
				details[k] = v
			}
		}
		return &SdkError{
			Code:    apiResp.Error.Code,
			Message: apiResp.Error.Message,
			Status:  status,
			TraceID: traceID,
			Details: details,
		}
	}

	return &SdkError{
		Code:    "unknown_error",
		Message: string(body),
		Status:  status,
		TraceID: traceID,
		Details: details,
	}
}

// generateTraceID produces a random UUID v4 string for request tracing.
func generateTraceID() string {
	var uuid [16]byte
	_, _ = rand.Read(uuid[:])
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
}
