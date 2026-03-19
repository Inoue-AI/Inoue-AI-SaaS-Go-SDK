package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// DownloadsAPI provides access to the content download endpoints.
type DownloadsAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// ContentDownloadCreateRequest is the body for POST /v1/downloads.
type ContentDownloadCreateRequest struct {
	URL        string `json:"url"`
	ModelID    string `json:"model_id,omitempty"`
	OwnerOrgID string `json:"owner_org_id,omitempty"`
}

// ContentDownloadBatchCreateRequest is the body for POST /v1/downloads/batch.
type ContentDownloadBatchCreateRequest struct {
	URLs       []string `json:"urls"`
	ModelID    string   `json:"model_id,omitempty"`
	OwnerOrgID string   `json:"owner_org_id,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// ContentDownloadResponse represents a content download record.
type ContentDownloadResponse struct {
	ID              string                 `json:"id"`
	URL             string                 `json:"url"`
	Platform        string                 `json:"platform"`
	Status          string                 `json:"status"`
	ModelID         string                 `json:"model_id"`
	OwnerUserID     string                 `json:"owner_user_id"`
	OwnerOrgID      string                 `json:"owner_org_id"`
	ProgressJSON    map[string]interface{} `json:"progress_json"`
	ErrorJSON       map[string]interface{} `json:"error_json"`
	ClaimedByWorker string                 `json:"claimed_by_worker_id"`
	CreatedAt       string                 `json:"created_at"`
	UpdatedAt       string                 `json:"updated_at"`
	FinishedAt      string                 `json:"finished_at"`
}

// ContentDownloadOutputResponse represents a single output from a download.
type ContentDownloadOutputResponse struct {
	ID         string `json:"id"`
	DownloadID string `json:"download_id"`
	AssetID    string `json:"asset_id"`
	AssetType  string `json:"asset_type"`
	StorageKey string `json:"storage_key"`
	StorageURL string `json:"storage_url"`
	Filename   string `json:"filename"`
	MimeType   string `json:"mime_type"`
	SizeBytes  int64  `json:"size_bytes"`
	CreatedAt  string `json:"created_at"`
}

// ContentDownloadBatchCreateResult is the response from a batch download creation.
type ContentDownloadBatchCreateResult struct {
	Downloads []ContentDownloadResponse `json:"downloads"`
	Errors    []map[string]interface{}  `json:"errors"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// Create creates a new content download request.
func (a *DownloadsAPI) Create(ctx context.Context, req ContentDownloadCreateRequest) (*ContentDownloadResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/downloads", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("create download: %w", err)
	}
	var result ContentDownloadResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create download decode: %w", err)
	}
	return &result, nil
}

// BatchCreate creates multiple content download requests at once.
func (a *DownloadsAPI) BatchCreate(ctx context.Context, req ContentDownloadBatchCreateRequest) (*ContentDownloadBatchCreateResult, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/downloads/batch", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("batch create downloads: %w", err)
	}
	var result ContentDownloadBatchCreateResult
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("batch create downloads decode: %w", err)
	}
	return &result, nil
}

// List returns a paginated list of content downloads.
func (a *DownloadsAPI) List(ctx context.Context, page, pageSize int) (*Page[ContentDownloadResponse], error) {
	path := fmt.Sprintf("/v1/downloads?page=%d&page_size=%d", page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list downloads: %w", err)
	}
	var result Page[ContentDownloadResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list downloads decode: %w", err)
	}
	return &result, nil
}

// Get retrieves a single content download by ID.
func (a *DownloadsAPI) Get(ctx context.Context, downloadID string) (*ContentDownloadResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/downloads/"+downloadID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get download: %w", err)
	}
	var result ContentDownloadResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get download decode: %w", err)
	}
	return &result, nil
}

// Outputs returns the outputs for a content download.
func (a *DownloadsAPI) Outputs(ctx context.Context, downloadID string) ([]ContentDownloadOutputResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/downloads/"+downloadID+"/outputs", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("download outputs: %w", err)
	}
	var result []ContentDownloadOutputResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("download outputs decode: %w", err)
	}
	return result, nil
}

// Cancel cancels a pending content download.
func (a *DownloadsAPI) Cancel(ctx context.Context, downloadID string) (*ContentDownloadResponse, error) {
	body := map[string]string{"download_id": downloadID}
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/downloads/cancel", body, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("cancel download: %w", err)
	}
	var result ContentDownloadResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("cancel download decode: %w", err)
	}
	return &result, nil
}

// Retry retries a failed content download.
func (a *DownloadsAPI) Retry(ctx context.Context, downloadID string) (*ContentDownloadResponse, error) {
	body := map[string]string{"download_id": downloadID}
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/downloads/retry", body, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("retry download: %w", err)
	}
	var result ContentDownloadResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("retry download decode: %w", err)
	}
	return &result, nil
}

// Delete deletes a content download by ID.
func (a *DownloadsAPI) Delete(ctx context.Context, downloadID string) error {
	body := map[string]string{"download_id": downloadID}
	if err := a.client.request(ctx, "POST", "/v1/downloads/delete", body, nil, nil); err != nil {
		return fmt.Errorf("delete download: %w", err)
	}
	return nil
}
