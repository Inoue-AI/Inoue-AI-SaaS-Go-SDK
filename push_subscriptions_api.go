package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// PushSubscriptionsAPI provides access to the push notification subscription endpoints.
type PushSubscriptionsAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// PushSubscriptionCreateRequest is the body for creating a push subscription.
type PushSubscriptionCreateRequest struct {
	Endpoint string                 `json:"endpoint"`
	Keys     map[string]string      `json:"keys"`
	UserAgent string                `json:"user_agent,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// PushSubscriptionResponse represents a push subscription record.
type PushSubscriptionResponse struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Endpoint  string            `json:"endpoint"`
	Keys      map[string]string `json:"keys"`
	UserAgent string            `json:"user_agent"`
	CreatedAt string            `json:"created_at"`
}

// VapidPublicKeyResponse represents the VAPID public key.
type VapidPublicKeyResponse struct {
	PublicKey string `json:"public_key"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// VapidPublicKey returns the server's VAPID public key for push notifications.
func (a *PushSubscriptionsAPI) VapidPublicKey(ctx context.Context) (*VapidPublicKeyResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/push-subscriptions/vapid-public-key", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("vapid public key: %w", err)
	}
	var result VapidPublicKeyResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("vapid public key decode: %w", err)
	}
	return &result, nil
}

// Create creates a new push subscription.
func (a *PushSubscriptionsAPI) Create(ctx context.Context, req PushSubscriptionCreateRequest) (*PushSubscriptionResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/push-subscriptions", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("create push subscription: %w", err)
	}
	var result PushSubscriptionResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create push subscription decode: %w", err)
	}
	return &result, nil
}

// List returns the list of push subscriptions.
func (a *PushSubscriptionsAPI) List(ctx context.Context) ([]PushSubscriptionResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/push-subscriptions", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list push subscriptions: %w", err)
	}
	var result []PushSubscriptionResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list push subscriptions decode: %w", err)
	}
	return result, nil
}

// Delete deletes a push subscription by ID.
func (a *PushSubscriptionsAPI) Delete(ctx context.Context, subscriptionID string) error {
	body := map[string]string{"subscription_id": subscriptionID}
	if err := a.client.request(ctx, "POST", "/v1/push-subscriptions/delete", body, nil, nil); err != nil {
		return fmt.Errorf("delete push subscription: %w", err)
	}
	return nil
}
