package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// WebhooksAPI provides access to the webhook endpoints.
type WebhooksAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// WebhookIngestResult represents the result of a webhook ingestion.
type WebhookIngestResult struct {
	Received bool                   `json:"received"`
	EventID  string                 `json:"event_id"`
	Details  map[string]interface{} `json:"details"`
}

// WebhookHealthResponse represents webhook system health.
type WebhookHealthResponse struct {
	Status    string                 `json:"status"`
	Providers map[string]interface{} `json:"providers"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// IngestMock sends a mock webhook event for testing.
func (a *WebhooksAPI) IngestMock(ctx context.Context, payload map[string]interface{}) (*WebhookIngestResult, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/webhooks/ingest/mock", payload, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("ingest mock webhook: %w", err)
	}
	var result WebhookIngestResult
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("ingest mock webhook decode: %w", err)
	}
	return &result, nil
}

// Health returns the webhook system health.
func (a *WebhooksAPI) Health(ctx context.Context) (*WebhookHealthResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/webhooks/health", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("webhook health: %w", err)
	}
	var result WebhookHealthResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("webhook health decode: %w", err)
	}
	return &result, nil
}
