package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// CollectionsAPI provides access to the collection management endpoints.
type CollectionsAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// CollectionCreateRequest is the body for creating a collection.
type CollectionCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	OwnerOrgID  string `json:"owner_org_id,omitempty"`
}

// CollectionUpdateRequest is the body for updating a collection.
type CollectionUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// CollectionItemAddRequest is the body for adding an item to a collection.
type CollectionItemAddRequest struct {
	AssetID  string `json:"asset_id"`
	Position int    `json:"position,omitempty"`
}

// CollectionItemUpdateRequest is the body for updating an item in a collection.
type CollectionItemUpdateRequest struct {
	Position int `json:"position"`
}

// CollectionShareRequest is the body for sharing a collection.
type CollectionShareRequest struct {
	UserID     string `json:"user_id,omitempty"`
	OrgID      string `json:"org_id,omitempty"`
	Permission string `json:"permission,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// CollectionResponse represents a collection record.
type CollectionResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerUserID string `json:"owner_user_id"`
	OwnerOrgID  string `json:"owner_org_id"`
	ItemCount   int    `json:"item_count"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CollectionItemResponse represents an item within a collection.
type CollectionItemResponse struct {
	ID           string              `json:"id"`
	CollectionID string              `json:"collection_id"`
	AssetID      string              `json:"asset_id"`
	Position     int                 `json:"position"`
	Asset        *AssetPublicResponse `json:"asset,omitempty"`
	CreatedAt    string              `json:"created_at"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// Create creates a new collection.
func (a *CollectionsAPI) Create(ctx context.Context, req CollectionCreateRequest) (*CollectionResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/collections", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("create collection: %w", err)
	}
	var result CollectionResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create collection decode: %w", err)
	}
	return &result, nil
}

// List returns a paginated list of collections.
func (a *CollectionsAPI) List(ctx context.Context, page, pageSize int) (*Page[CollectionResponse], error) {
	path := fmt.Sprintf("/v1/collections?page=%d&page_size=%d", page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list collections: %w", err)
	}
	var result Page[CollectionResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list collections decode: %w", err)
	}
	return &result, nil
}

// Get retrieves a single collection by ID.
func (a *CollectionsAPI) Get(ctx context.Context, collectionID string) (*CollectionResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/collections/"+collectionID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get collection: %w", err)
	}
	var result CollectionResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get collection decode: %w", err)
	}
	return &result, nil
}

// Update updates an existing collection.
func (a *CollectionsAPI) Update(ctx context.Context, collectionID string, req CollectionUpdateRequest) (*CollectionResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "PATCH", "/v1/collections/"+collectionID, req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("update collection: %w", err)
	}
	var result CollectionResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("update collection decode: %w", err)
	}
	return &result, nil
}

// Delete deletes a collection by ID.
func (a *CollectionsAPI) Delete(ctx context.Context, collectionID string) error {
	body := map[string]string{"collection_id": collectionID}
	if err := a.client.request(ctx, "POST", "/v1/collections/delete", body, nil, nil); err != nil {
		return fmt.Errorf("delete collection: %w", err)
	}
	return nil
}

// ListItems returns a paginated list of items in a collection.
func (a *CollectionsAPI) ListItems(ctx context.Context, collectionID string, page, pageSize int) (*Page[CollectionItemResponse], error) {
	path := fmt.Sprintf("/v1/collections/%s/items?page=%d&page_size=%d", collectionID, page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list collection items: %w", err)
	}
	var result Page[CollectionItemResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list collection items decode: %w", err)
	}
	return &result, nil
}

// AddItem adds an item to a collection.
func (a *CollectionsAPI) AddItem(ctx context.Context, collectionID string, req CollectionItemAddRequest) (*CollectionItemResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/collections/"+collectionID+"/items", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("add collection item: %w", err)
	}
	var result CollectionItemResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("add collection item decode: %w", err)
	}
	return &result, nil
}

// UpdateItem updates an item's position in a collection.
func (a *CollectionsAPI) UpdateItem(ctx context.Context, collectionID, itemID string, req CollectionItemUpdateRequest) (*CollectionItemResponse, error) {
	var apiResp ApiResponse
	path := fmt.Sprintf("/v1/collections/%s/items/%s", collectionID, itemID)
	if err := a.client.request(ctx, "PATCH", path, req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("update collection item: %w", err)
	}
	var result CollectionItemResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("update collection item decode: %w", err)
	}
	return &result, nil
}

// RemoveItem removes an item from a collection.
func (a *CollectionsAPI) RemoveItem(ctx context.Context, collectionID, itemID string) error {
	path := fmt.Sprintf("/v1/collections/%s/items/%s", collectionID, itemID)
	if err := a.client.request(ctx, "DELETE", path, nil, nil, nil); err != nil {
		return fmt.Errorf("remove collection item: %w", err)
	}
	return nil
}

// Link links a collection to a model.
func (a *CollectionsAPI) Link(ctx context.Context, collectionID, modelID string) error {
	body := map[string]string{
		"collection_id": collectionID,
		"model_id":      modelID,
	}
	if err := a.client.request(ctx, "POST", "/v1/collections/link", body, nil, nil); err != nil {
		return fmt.Errorf("link collection: %w", err)
	}
	return nil
}

// Unlink removes the link between a collection and a model.
func (a *CollectionsAPI) Unlink(ctx context.Context, collectionID, modelID string) error {
	body := map[string]string{
		"collection_id": collectionID,
		"model_id":      modelID,
	}
	if err := a.client.request(ctx, "POST", "/v1/collections/unlink", body, nil, nil); err != nil {
		return fmt.Errorf("unlink collection: %w", err)
	}
	return nil
}

// Share shares a collection with a user or organization.
func (a *CollectionsAPI) Share(ctx context.Context, collectionID string, req CollectionShareRequest) (*ModelShareGrant, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/collections/"+collectionID+"/share", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("share collection: %w", err)
	}
	var result ModelShareGrant
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("share collection decode: %w", err)
	}
	return &result, nil
}

// Revoke revokes a share grant on a collection.
func (a *CollectionsAPI) Revoke(ctx context.Context, collectionID, grantID string) error {
	path := fmt.Sprintf("/v1/collections/%s/share/%s", collectionID, grantID)
	if err := a.client.request(ctx, "DELETE", path, nil, nil, nil); err != nil {
		return fmt.Errorf("revoke collection share: %w", err)
	}
	return nil
}

// Shares returns the list of share grants for a collection.
func (a *CollectionsAPI) Shares(ctx context.Context, collectionID string) ([]ModelShareGrant, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/collections/"+collectionID+"/shares", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list collection shares: %w", err)
	}
	var result []ModelShareGrant
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list collection shares decode: %w", err)
	}
	return result, nil
}
