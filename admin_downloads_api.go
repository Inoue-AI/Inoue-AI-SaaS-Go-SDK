package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// AdminDownloadsAPI provides access to the admin download provider settings endpoints.
type AdminDownloadsAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// ProviderSettingUpdateRequest is the body for updating a provider setting.
type ProviderSettingUpdateRequest struct {
	ProviderID string                 `json:"provider_id"`
	Enabled    *bool                  `json:"enabled,omitempty"`
	Priority   *int                   `json:"priority,omitempty"`
	Config     map[string]interface{} `json:"config,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// ProviderSettingResponse represents a download provider setting.
type ProviderSettingResponse struct {
	ID         string                 `json:"id"`
	ProviderID string                 `json:"provider_id"`
	Name       string                 `json:"name"`
	Enabled    bool                   `json:"enabled"`
	Priority   int                    `json:"priority"`
	Config     map[string]interface{} `json:"config"`
	CreatedAt  string                 `json:"created_at"`
	UpdatedAt  string                 `json:"updated_at"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// ListProviderSettings returns the list of download provider settings.
func (a *AdminDownloadsAPI) ListProviderSettings(ctx context.Context) ([]ProviderSettingResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/admin/downloads/providers", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list provider settings: %w", err)
	}
	var result []ProviderSettingResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list provider settings decode: %w", err)
	}
	return result, nil
}

// UpdateProviderSetting updates a download provider setting.
func (a *AdminDownloadsAPI) UpdateProviderSetting(ctx context.Context, req ProviderSettingUpdateRequest) (*ProviderSettingResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/admin/downloads/providers/update", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("update provider setting: %w", err)
	}
	var result ProviderSettingResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("update provider setting decode: %w", err)
	}
	return &result, nil
}
