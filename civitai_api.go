package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// CivitaiAPI provides access to the CivitAI integration endpoints.
type CivitaiAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// CivitaiKeyCreateRequest is the body for creating a CivitAI key.
type CivitaiKeyCreateRequest struct {
	Name       string `json:"name"`
	APIKey     string `json:"api_key"`
	OwnerOrgID string `json:"owner_org_id,omitempty"`
}

// CivitaiKeyUpdateRequest is the body for updating a CivitAI key.
type CivitaiKeyUpdateRequest struct {
	Name   string `json:"name,omitempty"`
	APIKey string `json:"api_key,omitempty"`
}

// CivitaiShareRequest is the body for sharing a CivitAI key.
type CivitaiShareRequest struct {
	UserID     string `json:"user_id,omitempty"`
	OrgID      string `json:"org_id,omitempty"`
	Permission string `json:"permission,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// CivitaiKeyResponse represents a CivitAI API key.
type CivitaiKeyResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	OwnerUserID string `json:"owner_user_id"`
	OwnerOrgID  string `json:"owner_org_id"`
	Masked      string `json:"masked"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// ListKeys returns the list of CivitAI API keys.
func (a *CivitaiAPI) ListKeys(ctx context.Context) ([]CivitaiKeyResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/civitai/keys", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list civitai keys: %w", err)
	}
	var result []CivitaiKeyResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list civitai keys decode: %w", err)
	}
	return result, nil
}

// CreateKey creates a new CivitAI API key.
func (a *CivitaiAPI) CreateKey(ctx context.Context, req CivitaiKeyCreateRequest) (*CivitaiKeyResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/civitai/keys", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("create civitai key: %w", err)
	}
	var result CivitaiKeyResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create civitai key decode: %w", err)
	}
	return &result, nil
}

// UpdateKey updates an existing CivitAI API key.
func (a *CivitaiAPI) UpdateKey(ctx context.Context, keyID string, req CivitaiKeyUpdateRequest) (*CivitaiKeyResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "PATCH", "/v1/civitai/keys/"+keyID, req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("update civitai key: %w", err)
	}
	var result CivitaiKeyResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("update civitai key decode: %w", err)
	}
	return &result, nil
}

// DeleteKey deletes a CivitAI API key.
func (a *CivitaiAPI) DeleteKey(ctx context.Context, keyID string) error {
	body := map[string]string{"key_id": keyID}
	if err := a.client.request(ctx, "POST", "/v1/civitai/keys/delete", body, nil, nil); err != nil {
		return fmt.Errorf("delete civitai key: %w", err)
	}
	return nil
}

// ShareKey shares a CivitAI key with a user or organization.
func (a *CivitaiAPI) ShareKey(ctx context.Context, keyID string, req CivitaiShareRequest) (*ModelShareGrant, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/civitai/keys/"+keyID+"/share", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("share civitai key: %w", err)
	}
	var result ModelShareGrant
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("share civitai key decode: %w", err)
	}
	return &result, nil
}

// RevokeKey revokes a share on a CivitAI key.
func (a *CivitaiAPI) RevokeKey(ctx context.Context, keyID, grantID string) error {
	path := fmt.Sprintf("/v1/civitai/keys/%s/share/%s", keyID, grantID)
	if err := a.client.request(ctx, "DELETE", path, nil, nil, nil); err != nil {
		return fmt.Errorf("revoke civitai key share: %w", err)
	}
	return nil
}

// ListShares returns the list of shares for a CivitAI key.
func (a *CivitaiAPI) ListShares(ctx context.Context, keyID string) ([]ModelShareGrant, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/civitai/keys/"+keyID+"/shares", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list civitai key shares: %w", err)
	}
	var result []ModelShareGrant
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list civitai key shares decode: %w", err)
	}
	return result, nil
}
