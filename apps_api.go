package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// AppsAPI provides access to the application management endpoints.
type AppsAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// AppResponse represents an application record.
type AppResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Slug        string                 `json:"slug"`
	IconURL     string                 `json:"icon_url"`
	Status      string                 `json:"status"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

// AppAccessResponse represents access status for an application.
type AppAccessResponse struct {
	AppID     string `json:"app_id"`
	HasAccess bool   `json:"has_access"`
	Plan      string `json:"plan"`
	ExpiresAt string `json:"expires_at"`
}

// AppVersionResponse represents an application version.
type AppVersionResponse struct {
	ID        string                 `json:"id"`
	AppID     string                 `json:"app_id"`
	Version   string                 `json:"version"`
	Changelog string                 `json:"changelog"`
	URL       string                 `json:"url"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt string                 `json:"created_at"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// List returns the list of available applications.
func (a *AppsAPI) List(ctx context.Context) ([]AppResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/apps", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list apps: %w", err)
	}
	var result []AppResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list apps decode: %w", err)
	}
	return result, nil
}

// Access returns the access status for an application.
func (a *AppsAPI) Access(ctx context.Context, appID string) (*AppAccessResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/apps/"+appID+"/access", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("app access: %w", err)
	}
	var result AppAccessResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("app access decode: %w", err)
	}
	return &result, nil
}

// LatestVersion returns the latest version of an application.
func (a *AppsAPI) LatestVersion(ctx context.Context, appID string) (*AppVersionResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/apps/"+appID+"/versions/latest", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("app latest version: %w", err)
	}
	var result AppVersionResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("app latest version decode: %w", err)
	}
	return &result, nil
}

// ListVersions returns the list of versions for an application.
func (a *AppsAPI) ListVersions(ctx context.Context, appID string) ([]AppVersionResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/apps/"+appID+"/versions", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list app versions: %w", err)
	}
	var result []AppVersionResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list app versions decode: %w", err)
	}
	return result, nil
}
