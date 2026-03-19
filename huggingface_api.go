package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// HuggingFaceAPI provides access to the HuggingFace integration endpoints.
type HuggingFaceAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// HuggingFaceKeyCreateRequest is the body for creating a HuggingFace key.
type HuggingFaceKeyCreateRequest struct {
	Name       string `json:"name"`
	APIKey     string `json:"api_key"`
	OwnerOrgID string `json:"owner_org_id,omitempty"`
}

// HuggingFaceKeyUpdateRequest is the body for updating a HuggingFace key.
type HuggingFaceKeyUpdateRequest struct {
	Name   string `json:"name,omitempty"`
	APIKey string `json:"api_key,omitempty"`
}

// HuggingFaceShareRequest is the body for sharing a HuggingFace key.
type HuggingFaceShareRequest struct {
	UserID     string `json:"user_id,omitempty"`
	OrgID      string `json:"org_id,omitempty"`
	Permission string `json:"permission,omitempty"`
}

// HuggingFaceProbeRequest is the body for probing a HuggingFace repo.
type HuggingFaceProbeRequest struct {
	RepoID string `json:"repo_id"`
	KeyID  string `json:"key_id,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// HuggingFaceKeyResponse represents a HuggingFace API key.
type HuggingFaceKeyResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	OwnerUserID string `json:"owner_user_id"`
	OwnerOrgID  string `json:"owner_org_id"`
	Masked      string `json:"masked"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// HuggingFaceProbeResponse represents the result of a HuggingFace repo probe.
type HuggingFaceProbeResponse struct {
	RepoID    string                 `json:"repo_id"`
	Exists    bool                   `json:"exists"`
	RepoType  string                 `json:"repo_type"`
	ModelInfo  map[string]interface{} `json:"model_info"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// ListKeys returns the list of HuggingFace API keys.
func (a *HuggingFaceAPI) ListKeys(ctx context.Context) ([]HuggingFaceKeyResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/huggingface/keys", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list huggingface keys: %w", err)
	}
	var result []HuggingFaceKeyResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list huggingface keys decode: %w", err)
	}
	return result, nil
}

// CreateKey creates a new HuggingFace API key.
func (a *HuggingFaceAPI) CreateKey(ctx context.Context, req HuggingFaceKeyCreateRequest) (*HuggingFaceKeyResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/huggingface/keys", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("create huggingface key: %w", err)
	}
	var result HuggingFaceKeyResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create huggingface key decode: %w", err)
	}
	return &result, nil
}

// UpdateKey updates an existing HuggingFace API key.
func (a *HuggingFaceAPI) UpdateKey(ctx context.Context, keyID string, req HuggingFaceKeyUpdateRequest) (*HuggingFaceKeyResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "PATCH", "/v1/huggingface/keys/"+keyID, req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("update huggingface key: %w", err)
	}
	var result HuggingFaceKeyResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("update huggingface key decode: %w", err)
	}
	return &result, nil
}

// DeleteKey deletes a HuggingFace API key.
func (a *HuggingFaceAPI) DeleteKey(ctx context.Context, keyID string) error {
	body := map[string]string{"key_id": keyID}
	if err := a.client.request(ctx, "POST", "/v1/huggingface/keys/delete", body, nil, nil); err != nil {
		return fmt.Errorf("delete huggingface key: %w", err)
	}
	return nil
}

// ShareKey shares a HuggingFace key with a user or organization.
func (a *HuggingFaceAPI) ShareKey(ctx context.Context, keyID string, req HuggingFaceShareRequest) (*ModelShareGrant, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/huggingface/keys/"+keyID+"/share", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("share huggingface key: %w", err)
	}
	var result ModelShareGrant
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("share huggingface key decode: %w", err)
	}
	return &result, nil
}

// RevokeKey revokes a share on a HuggingFace key.
func (a *HuggingFaceAPI) RevokeKey(ctx context.Context, keyID, grantID string) error {
	path := fmt.Sprintf("/v1/huggingface/keys/%s/share/%s", keyID, grantID)
	if err := a.client.request(ctx, "DELETE", path, nil, nil, nil); err != nil {
		return fmt.Errorf("revoke huggingface key share: %w", err)
	}
	return nil
}

// ListShares returns the list of shares for a HuggingFace key.
func (a *HuggingFaceAPI) ListShares(ctx context.Context, keyID string) ([]ModelShareGrant, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/huggingface/keys/"+keyID+"/shares", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list huggingface key shares: %w", err)
	}
	var result []ModelShareGrant
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list huggingface key shares decode: %w", err)
	}
	return result, nil
}

// ProbeRepo probes a HuggingFace repository for model info.
func (a *HuggingFaceAPI) ProbeRepo(ctx context.Context, req HuggingFaceProbeRequest) (*HuggingFaceProbeResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/huggingface/probe", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("probe huggingface repo: %w", err)
	}
	var result HuggingFaceProbeResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("probe huggingface repo decode: %w", err)
	}
	return &result, nil
}
