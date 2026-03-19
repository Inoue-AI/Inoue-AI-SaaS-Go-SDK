package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// PoolingAPI provides access to the model pooling and invite endpoints.
type PoolingAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// PoolProposeRequest is the body for proposing a pool contribution.
type PoolProposeRequest struct {
	ModelID    string  `json:"model_id"`
	Credits    float64 `json:"credits"`
	Message    string  `json:"message,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// PoolResponse represents a model pool.
type PoolResponse struct {
	ID            string  `json:"id"`
	ModelID       string  `json:"model_id"`
	TotalCredits  float64 `json:"total_credits"`
	Contributors  int     `json:"contributors"`
	Status        string  `json:"status"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

// PoolInviteResponse represents a pool invite.
type PoolInviteResponse struct {
	ID        string  `json:"id"`
	PoolID    string  `json:"pool_id"`
	ModelID   string  `json:"model_id"`
	UserID    string  `json:"user_id"`
	Credits   float64 `json:"credits"`
	Status    string  `json:"status"`
	Message   string  `json:"message"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// ModelPool retrieves the pool for a specific model.
func (a *PoolingAPI) ModelPool(ctx context.Context, modelID string) (*PoolResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/pooling/models/"+modelID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get model pool: %w", err)
	}
	var result PoolResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get model pool decode: %w", err)
	}
	return &result, nil
}

// Invites returns the list of pool invites for the current user.
func (a *PoolingAPI) Invites(ctx context.Context) ([]PoolInviteResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/pooling/invites", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list pool invites: %w", err)
	}
	var result []PoolInviteResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list pool invites decode: %w", err)
	}
	return result, nil
}

// Propose proposes a contribution to a model pool.
func (a *PoolingAPI) Propose(ctx context.Context, req PoolProposeRequest) (*PoolInviteResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/pooling/propose", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("propose pool contribution: %w", err)
	}
	var result PoolInviteResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("propose pool contribution decode: %w", err)
	}
	return &result, nil
}

// Accept accepts a pool invite.
func (a *PoolingAPI) Accept(ctx context.Context, inviteID string) (*PoolInviteResponse, error) {
	body := map[string]string{"invite_id": inviteID}
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/pooling/accept", body, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("accept pool invite: %w", err)
	}
	var result PoolInviteResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("accept pool invite decode: %w", err)
	}
	return &result, nil
}

// Decline declines a pool invite.
func (a *PoolingAPI) Decline(ctx context.Context, inviteID string) (*PoolInviteResponse, error) {
	body := map[string]string{"invite_id": inviteID}
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/pooling/decline", body, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("decline pool invite: %w", err)
	}
	var result PoolInviteResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("decline pool invite decode: %w", err)
	}
	return &result, nil
}

// Revoke revokes a pool invite.
func (a *PoolingAPI) Revoke(ctx context.Context, inviteID string) error {
	body := map[string]string{"invite_id": inviteID}
	if err := a.client.request(ctx, "POST", "/v1/pooling/revoke", body, nil, nil); err != nil {
		return fmt.Errorf("revoke pool invite: %w", err)
	}
	return nil
}
