package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// LorasAPI provides access to the LoRA management endpoints.
type LorasAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// LoraCreateRequest is the body for creating a LoRA.
type LoraCreateRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	ModelID     string                 `json:"model_id,omitempty"`
	AssetID     string                 `json:"asset_id,omitempty"`
	OwnerOrgID  string                 `json:"owner_org_id,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// LoraUpdateRequest is the body for updating a LoRA.
type LoraUpdateRequest struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// LoraResponse represents a LoRA record.
type LoraResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	ModelID     string                 `json:"model_id"`
	AssetID     string                 `json:"asset_id"`
	OwnerUserID string                 `json:"owner_user_id"`
	OwnerOrgID  string                 `json:"owner_org_id"`
	Status      string                 `json:"status"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// List returns a paginated list of LoRAs.
func (a *LorasAPI) List(ctx context.Context, page, pageSize int) (*Page[LoraResponse], error) {
	path := fmt.Sprintf("/v1/loras?page=%d&page_size=%d", page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list loras: %w", err)
	}
	var result Page[LoraResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list loras decode: %w", err)
	}
	return &result, nil
}

// Create creates a new LoRA.
func (a *LorasAPI) Create(ctx context.Context, req LoraCreateRequest) (*LoraResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/loras", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("create lora: %w", err)
	}
	var result LoraResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create lora decode: %w", err)
	}
	return &result, nil
}

// Update updates an existing LoRA.
func (a *LorasAPI) Update(ctx context.Context, loraID string, req LoraUpdateRequest) (*LoraResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "PATCH", "/v1/loras/"+loraID, req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("update lora: %w", err)
	}
	var result LoraResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("update lora decode: %w", err)
	}
	return &result, nil
}

// Delete deletes a LoRA by ID.
func (a *LorasAPI) Delete(ctx context.Context, loraID string) error {
	body := map[string]string{"lora_id": loraID}
	if err := a.client.request(ctx, "POST", "/v1/loras/delete", body, nil, nil); err != nil {
		return fmt.Errorf("delete lora: %w", err)
	}
	return nil
}
