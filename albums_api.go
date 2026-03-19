package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// AlbumsAPI provides access to the album management endpoints.
type AlbumsAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// AlbumCreateRequest is the body for creating an album.
type AlbumCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ModelID     string `json:"model_id,omitempty"`
	OwnerOrgID  string `json:"owner_org_id,omitempty"`
}

// AlbumUpdateRequest is the body for updating an album.
type AlbumUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// AlbumItemAddRequest is the body for adding an item to an album.
type AlbumItemAddRequest struct {
	AssetID  string `json:"asset_id"`
	Position int    `json:"position,omitempty"`
}

// AlbumBulkAddItemsRequest is the body for bulk-adding items to an album.
type AlbumBulkAddItemsRequest struct {
	AssetIDs []string `json:"asset_ids"`
}

// AlbumItemUpdateRequest is the body for updating an item in an album.
type AlbumItemUpdateRequest struct {
	Position int `json:"position"`
}

// AlbumShareRequest is the body for sharing an album.
type AlbumShareRequest struct {
	UserID     string `json:"user_id,omitempty"`
	OrgID      string `json:"org_id,omitempty"`
	Permission string `json:"permission,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// AlbumResponse represents an album record.
type AlbumResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ModelID     string `json:"model_id"`
	OwnerUserID string `json:"owner_user_id"`
	OwnerOrgID  string `json:"owner_org_id"`
	ItemCount   int    `json:"item_count"`
	CoverURL    string `json:"cover_url"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// AlbumItemResponse represents an item within an album.
type AlbumItemResponse struct {
	ID        string              `json:"id"`
	AlbumID   string              `json:"album_id"`
	AssetID   string              `json:"asset_id"`
	Position  int                 `json:"position"`
	Asset     *AssetPublicResponse `json:"asset,omitempty"`
	CreatedAt string              `json:"created_at"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// List returns a paginated list of albums.
func (a *AlbumsAPI) List(ctx context.Context, page, pageSize int) (*Page[AlbumResponse], error) {
	path := fmt.Sprintf("/v1/albums?page=%d&page_size=%d", page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list albums: %w", err)
	}
	var result Page[AlbumResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list albums decode: %w", err)
	}
	return &result, nil
}

// Create creates a new album.
func (a *AlbumsAPI) Create(ctx context.Context, req AlbumCreateRequest) (*AlbumResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/albums", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("create album: %w", err)
	}
	var result AlbumResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create album decode: %w", err)
	}
	return &result, nil
}

// Get retrieves a single album by ID.
func (a *AlbumsAPI) Get(ctx context.Context, albumID string) (*AlbumResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/albums/"+albumID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get album: %w", err)
	}
	var result AlbumResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get album decode: %w", err)
	}
	return &result, nil
}

// Update updates an existing album.
func (a *AlbumsAPI) Update(ctx context.Context, albumID string, req AlbumUpdateRequest) (*AlbumResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "PATCH", "/v1/albums/"+albumID, req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("update album: %w", err)
	}
	var result AlbumResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("update album decode: %w", err)
	}
	return &result, nil
}

// Delete deletes an album by ID.
func (a *AlbumsAPI) Delete(ctx context.Context, albumID string) error {
	body := map[string]string{"album_id": albumID}
	if err := a.client.request(ctx, "POST", "/v1/albums/delete", body, nil, nil); err != nil {
		return fmt.Errorf("delete album: %w", err)
	}
	return nil
}

// ListItems returns a paginated list of items in an album.
func (a *AlbumsAPI) ListItems(ctx context.Context, albumID string, page, pageSize int) (*Page[AlbumItemResponse], error) {
	path := fmt.Sprintf("/v1/albums/%s/items?page=%d&page_size=%d", albumID, page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list album items: %w", err)
	}
	var result Page[AlbumItemResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list album items decode: %w", err)
	}
	return &result, nil
}

// AddItem adds an item to an album.
func (a *AlbumsAPI) AddItem(ctx context.Context, albumID string, req AlbumItemAddRequest) (*AlbumItemResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/albums/"+albumID+"/items", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("add album item: %w", err)
	}
	var result AlbumItemResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("add album item decode: %w", err)
	}
	return &result, nil
}

// BulkAddItems adds multiple items to an album at once.
func (a *AlbumsAPI) BulkAddItems(ctx context.Context, albumID string, req AlbumBulkAddItemsRequest) ([]AlbumItemResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/albums/"+albumID+"/items/bulk", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("bulk add album items: %w", err)
	}
	var result []AlbumItemResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("bulk add album items decode: %w", err)
	}
	return result, nil
}

// UpdateItem updates an item's position in an album.
func (a *AlbumsAPI) UpdateItem(ctx context.Context, albumID, itemID string, req AlbumItemUpdateRequest) (*AlbumItemResponse, error) {
	var apiResp ApiResponse
	path := fmt.Sprintf("/v1/albums/%s/items/%s", albumID, itemID)
	if err := a.client.request(ctx, "PATCH", path, req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("update album item: %w", err)
	}
	var result AlbumItemResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("update album item decode: %w", err)
	}
	return &result, nil
}

// RemoveItem removes an item from an album.
func (a *AlbumsAPI) RemoveItem(ctx context.Context, albumID, itemID string) error {
	path := fmt.Sprintf("/v1/albums/%s/items/%s", albumID, itemID)
	if err := a.client.request(ctx, "DELETE", path, nil, nil, nil); err != nil {
		return fmt.Errorf("remove album item: %w", err)
	}
	return nil
}

// Link links an album to a model.
func (a *AlbumsAPI) Link(ctx context.Context, albumID, modelID string) error {
	body := map[string]string{
		"album_id": albumID,
		"model_id": modelID,
	}
	if err := a.client.request(ctx, "POST", "/v1/albums/link", body, nil, nil); err != nil {
		return fmt.Errorf("link album: %w", err)
	}
	return nil
}

// Unlink removes the link between an album and a model.
func (a *AlbumsAPI) Unlink(ctx context.Context, albumID, modelID string) error {
	body := map[string]string{
		"album_id": albumID,
		"model_id": modelID,
	}
	if err := a.client.request(ctx, "POST", "/v1/albums/unlink", body, nil, nil); err != nil {
		return fmt.Errorf("unlink album: %w", err)
	}
	return nil
}

// Share shares an album with a user or organization.
func (a *AlbumsAPI) Share(ctx context.Context, albumID string, req AlbumShareRequest) (*ModelShareGrant, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/albums/"+albumID+"/share", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("share album: %w", err)
	}
	var result ModelShareGrant
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("share album decode: %w", err)
	}
	return &result, nil
}

// Revoke revokes a share grant on an album.
func (a *AlbumsAPI) Revoke(ctx context.Context, albumID, grantID string) error {
	path := fmt.Sprintf("/v1/albums/%s/share/%s", albumID, grantID)
	if err := a.client.request(ctx, "DELETE", path, nil, nil, nil); err != nil {
		return fmt.Errorf("revoke album share: %w", err)
	}
	return nil
}

// Shares returns the list of share grants for an album.
func (a *AlbumsAPI) Shares(ctx context.Context, albumID string) ([]ModelShareGrant, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/albums/"+albumID+"/shares", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list album shares: %w", err)
	}
	var result []ModelShareGrant
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list album shares decode: %w", err)
	}
	return result, nil
}
