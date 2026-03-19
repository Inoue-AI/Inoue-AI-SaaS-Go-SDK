package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// SystemAPI provides access to the system health and configuration endpoints.
type SystemAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// HealthResponse represents the system health status.
type HealthResponse struct {
	Status    string                 `json:"status"`
	Version   string                 `json:"version"`
	Uptime    string                 `json:"uptime"`
	Services  map[string]interface{} `json:"services"`
}

// VersionResponse represents the system version info.
type VersionResponse struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
}

// EnumsResponse represents the available system enum values.
type EnumsResponse struct {
	JobTypes     []string               `json:"job_types"`
	JobStatuses  []string               `json:"job_statuses"`
	AssetTypes   []string               `json:"asset_types"`
	Platforms    []string               `json:"platforms"`
	Extra        map[string]interface{} `json:"extra"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// Health returns the system health status.
func (a *SystemAPI) Health(ctx context.Context) (*HealthResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/system/health", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("system health: %w", err)
	}
	var result HealthResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("system health decode: %w", err)
	}
	return &result, nil
}

// Version returns the system version information.
func (a *SystemAPI) Version(ctx context.Context) (*VersionResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/system/version", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("system version: %w", err)
	}
	var result VersionResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("system version decode: %w", err)
	}
	return &result, nil
}

// Enums returns the available system enum values.
func (a *SystemAPI) Enums(ctx context.Context) (*EnumsResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/system/enums", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("system enums: %w", err)
	}
	var result EnumsResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("system enums decode: %w", err)
	}
	return &result, nil
}
