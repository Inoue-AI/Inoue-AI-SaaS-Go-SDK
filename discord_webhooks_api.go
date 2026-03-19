package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// DiscordWebhooksAPI provides access to the Discord webhook integration endpoints.
type DiscordWebhooksAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// DiscordWebhookCreateRequest is the body for creating a Discord webhook.
type DiscordWebhookCreateRequest struct {
	Name       string   `json:"name"`
	URL        string   `json:"url"`
	ModelID    string   `json:"model_id,omitempty"`
	Events     []string `json:"events,omitempty"`
	OwnerOrgID string   `json:"owner_org_id,omitempty"`
}

// DiscordWebhookUpdateRequest is the body for updating a Discord webhook.
type DiscordWebhookUpdateRequest struct {
	Name   string   `json:"name,omitempty"`
	URL    string   `json:"url,omitempty"`
	Events []string `json:"events,omitempty"`
	Active *bool    `json:"active,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// DiscordWebhookResponse represents a Discord webhook record.
type DiscordWebhookResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	ModelID     string   `json:"model_id"`
	Events      []string `json:"events"`
	Active      bool     `json:"active"`
	OwnerUserID string   `json:"owner_user_id"`
	OwnerOrgID  string   `json:"owner_org_id"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// DiscordWebhookTestResult represents the result of a webhook test.
type DiscordWebhookTestResult struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// List returns the list of Discord webhooks.
func (a *DiscordWebhooksAPI) List(ctx context.Context) ([]DiscordWebhookResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/discord-webhooks", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list discord webhooks: %w", err)
	}
	var result []DiscordWebhookResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list discord webhooks decode: %w", err)
	}
	return result, nil
}

// Create creates a new Discord webhook.
func (a *DiscordWebhooksAPI) Create(ctx context.Context, req DiscordWebhookCreateRequest) (*DiscordWebhookResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/discord-webhooks", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("create discord webhook: %w", err)
	}
	var result DiscordWebhookResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create discord webhook decode: %w", err)
	}
	return &result, nil
}

// Update updates an existing Discord webhook.
func (a *DiscordWebhooksAPI) Update(ctx context.Context, webhookID string, req DiscordWebhookUpdateRequest) (*DiscordWebhookResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "PATCH", "/v1/discord-webhooks/"+webhookID, req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("update discord webhook: %w", err)
	}
	var result DiscordWebhookResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("update discord webhook decode: %w", err)
	}
	return &result, nil
}

// Delete deletes a Discord webhook by ID.
func (a *DiscordWebhooksAPI) Delete(ctx context.Context, webhookID string) error {
	body := map[string]string{"webhook_id": webhookID}
	if err := a.client.request(ctx, "POST", "/v1/discord-webhooks/delete", body, nil, nil); err != nil {
		return fmt.Errorf("delete discord webhook: %w", err)
	}
	return nil
}

// Test sends a test message to a Discord webhook.
func (a *DiscordWebhooksAPI) Test(ctx context.Context, webhookID string) (*DiscordWebhookTestResult, error) {
	body := map[string]string{"webhook_id": webhookID}
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/discord-webhooks/test", body, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("test discord webhook: %w", err)
	}
	var result DiscordWebhookTestResult
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("test discord webhook decode: %w", err)
	}
	return &result, nil
}
