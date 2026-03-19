package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// AssetsAPI provides access to the asset management endpoints.
type AssetsAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// AssetCreateRequest is the body for POST /v1/assets.
type AssetCreateRequest struct {
	AssetType  string                 `json:"asset_type"`
	StorageKey string                 `json:"storage_key,omitempty"`
	Filename   string                 `json:"filename,omitempty"`
	MimeType   string                 `json:"mime_type,omitempty"`
	SizeBytes  int64                  `json:"size_bytes,omitempty"`
	ModelID    string                 `json:"model_id,omitempty"`
	OwnerOrgID string                 `json:"owner_org_id,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// AssetUpdateRequest is the body for PATCH /v1/assets/{id}.
type AssetUpdateRequest struct {
	Filename string                 `json:"filename,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Tags     []string               `json:"tags,omitempty"`
}

// AssetLinkRequest is the body for POST /v1/assets/link.
type AssetLinkRequest struct {
	AssetID string `json:"asset_id"`
	ModelID string `json:"model_id"`
}

// AssetLinkResponse represents an asset-model link.
type AssetLinkResponse struct {
	ID        string `json:"id"`
	AssetID   string `json:"asset_id"`
	ModelID   string `json:"model_id"`
	CreatedAt string `json:"created_at"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// List returns a paginated list of assets.
func (a *AssetsAPI) List(ctx context.Context, page, pageSize int) (*Page[AssetPublicResponse], error) {
	path := fmt.Sprintf("/v1/assets?page=%d&page_size=%d", page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list assets: %w", err)
	}
	var result Page[AssetPublicResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list assets decode: %w", err)
	}
	return &result, nil
}

// Create creates a new asset record.
func (a *AssetsAPI) Create(ctx context.Context, req AssetCreateRequest) (*AssetPublicResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/assets", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("create asset: %w", err)
	}
	var result AssetPublicResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create asset decode: %w", err)
	}
	return &result, nil
}

// Get retrieves a single asset by ID.
func (a *AssetsAPI) Get(ctx context.Context, assetID string) (*AssetPublicResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/assets/"+assetID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get asset: %w", err)
	}
	var result AssetPublicResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get asset decode: %w", err)
	}
	return &result, nil
}

// Update updates an existing asset.
func (a *AssetsAPI) Update(ctx context.Context, assetID string, req AssetUpdateRequest) (*AssetPublicResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "PATCH", "/v1/assets/"+assetID, req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("update asset: %w", err)
	}
	var result AssetPublicResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("update asset decode: %w", err)
	}
	return &result, nil
}

// Delete deletes an asset by ID.
func (a *AssetsAPI) Delete(ctx context.Context, assetID string) error {
	body := map[string]string{"asset_id": assetID}
	if err := a.client.request(ctx, "POST", "/v1/assets/delete", body, nil, nil); err != nil {
		return fmt.Errorf("delete asset: %w", err)
	}
	return nil
}

// Link links an asset to a model.
func (a *AssetsAPI) Link(ctx context.Context, req AssetLinkRequest) (*AssetLinkResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/assets/link", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("link asset: %w", err)
	}
	var result AssetLinkResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("link asset decode: %w", err)
	}
	return &result, nil
}

// Unlink removes the link between an asset and a model.
func (a *AssetsAPI) Unlink(ctx context.Context, assetID, modelID string) error {
	body := map[string]string{
		"asset_id": assetID,
		"model_id": modelID,
	}
	if err := a.client.request(ctx, "POST", "/v1/assets/unlink", body, nil, nil); err != nil {
		return fmt.Errorf("unlink asset: %w", err)
	}
	return nil
}

// Archive archives an asset.
func (a *AssetsAPI) Archive(ctx context.Context, assetID string) (*AssetPublicResponse, error) {
	body := map[string]string{"asset_id": assetID}
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/assets/archive", body, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("archive asset: %w", err)
	}
	var result AssetPublicResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("archive asset decode: %w", err)
	}
	return &result, nil
}

// Permadelete permanently deletes an asset.
func (a *AssetsAPI) Permadelete(ctx context.Context, assetID string) error {
	body := map[string]string{"asset_id": assetID}
	if err := a.client.request(ctx, "POST", "/v1/assets/permadelete", body, nil, nil); err != nil {
		return fmt.Errorf("permadelete asset: %w", err)
	}
	return nil
}

// ModelAssets returns a paginated list of assets for a model.
func (a *AssetsAPI) ModelAssets(ctx context.Context, modelID string, page, pageSize int) (*Page[AssetPublicResponse], error) {
	path := fmt.Sprintf("/v1/assets/models/%s?page=%d&page_size=%d", modelID, page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list model assets: %w", err)
	}
	var result Page[AssetPublicResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list model assets decode: %w", err)
	}
	return &result, nil
}

// UploadURL requests a pre-signed upload URL for an asset.
func (a *AssetsAPI) UploadURL(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/assets/upload", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("asset upload url: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("asset upload url decode: %w", err)
	}
	return result, nil
}

// GetURL returns the access URL for an asset.
func (a *AssetsAPI) GetURL(ctx context.Context, assetID string) (map[string]interface{}, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/assets/"+assetID+"/url", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get asset url: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get asset url decode: %w", err)
	}
	return result, nil
}

// Download returns the download URL for an asset.
func (a *AssetsAPI) Download(ctx context.Context, assetID string) (map[string]interface{}, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/assets/"+assetID+"/download", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("download asset: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("download asset decode: %w", err)
	}
	return result, nil
}

// BulkUpdate performs a bulk update on multiple assets.
func (a *AssetsAPI) BulkUpdate(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/assets/bulk-update", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("bulk update assets: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("bulk update assets decode: %w", err)
	}
	return result, nil
}
