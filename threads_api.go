package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// ThreadsAPI provides access to the Threads (Meta) integration endpoints.
type ThreadsAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// ThreadsConnectRequest is the body for initiating a Threads connection.
type ThreadsConnectRequest struct {
	ModelID     string `json:"model_id"`
	RedirectURI string `json:"redirect_uri,omitempty"`
}

// ThreadsCallbackRequest is the body for the OAuth callback.
type ThreadsCallbackRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

// ThreadsPostRequest is the body for creating a Threads post.
type ThreadsPostRequest struct {
	ModelID    string                 `json:"model_id"`
	Text       string                 `json:"text,omitempty"`
	AssetIDs   []string               `json:"asset_ids,omitempty"`
	ReplyToID  string                 `json:"reply_to_id,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// ThreadsConnectResponse represents a Threads connection response.
type ThreadsConnectResponse struct {
	AuthURL string `json:"auth_url"`
	State   string `json:"state"`
}

// ThreadsAccountResponse represents a connected Threads account.
type ThreadsAccountResponse struct {
	ID          string `json:"id"`
	ModelID     string `json:"model_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Connected   bool   `json:"connected"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ThreadsPostResponse represents a Threads post result.
type ThreadsPostResponse struct {
	ID        string                 `json:"id"`
	PostID    string                 `json:"post_id"`
	Status    string                 `json:"status"`
	Text      string                 `json:"text"`
	URL       string                 `json:"url"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt string                 `json:"created_at"`
}

// ThreadsMediaResponse represents a Threads media item.
type ThreadsMediaResponse struct {
	ID        string                 `json:"id"`
	MediaID   string                 `json:"media_id"`
	Type      string                 `json:"type"`
	URL       string                 `json:"url"`
	Text      string                 `json:"text"`
	Stats     map[string]interface{} `json:"stats"`
	CreatedAt string                 `json:"created_at"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// Connect initiates a Threads OAuth connection.
func (a *ThreadsAPI) Connect(ctx context.Context, req ThreadsConnectRequest) (*ThreadsConnectResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/threads/connect", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("threads connect: %w", err)
	}
	var result ThreadsConnectResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("threads connect decode: %w", err)
	}
	return &result, nil
}

// Callback handles the Threads OAuth callback.
func (a *ThreadsAPI) Callback(ctx context.Context, req ThreadsCallbackRequest) (*ThreadsAccountResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/threads/callback", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("threads callback: %w", err)
	}
	var result ThreadsAccountResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("threads callback decode: %w", err)
	}
	return &result, nil
}

// Accounts returns the list of connected Threads accounts.
func (a *ThreadsAPI) Accounts(ctx context.Context) ([]ThreadsAccountResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/threads/accounts", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list threads accounts: %w", err)
	}
	var result []ThreadsAccountResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list threads accounts decode: %w", err)
	}
	return result, nil
}

// Account retrieves a specific Threads account by model ID.
func (a *ThreadsAPI) Account(ctx context.Context, modelID string) (*ThreadsAccountResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/threads/accounts/"+modelID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get threads account: %w", err)
	}
	var result ThreadsAccountResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get threads account decode: %w", err)
	}
	return &result, nil
}

// Disconnect disconnects a Threads account.
func (a *ThreadsAPI) Disconnect(ctx context.Context, modelID string) error {
	body := map[string]string{"model_id": modelID}
	if err := a.client.request(ctx, "POST", "/v1/threads/disconnect", body, nil, nil); err != nil {
		return fmt.Errorf("threads disconnect: %w", err)
	}
	return nil
}

// Post creates a Threads post.
func (a *ThreadsAPI) Post(ctx context.Context, req ThreadsPostRequest) (*ThreadsPostResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/threads/post", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("threads post: %w", err)
	}
	var result ThreadsPostResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("threads post decode: %w", err)
	}
	return &result, nil
}

// Media returns the list of Threads media for a model.
func (a *ThreadsAPI) Media(ctx context.Context, modelID string, page, pageSize int) (*Page[ThreadsMediaResponse], error) {
	path := fmt.Sprintf("/v1/threads/media/%s?page=%d&page_size=%d", modelID, page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list threads media: %w", err)
	}
	var result Page[ThreadsMediaResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list threads media decode: %w", err)
	}
	return &result, nil
}

// RefreshToken refreshes the Threads access token for a model.
func (a *ThreadsAPI) RefreshToken(ctx context.Context, modelID string) (*ThreadsAccountResponse, error) {
	body := map[string]string{"model_id": modelID}
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/threads/refresh", body, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("threads refresh token: %w", err)
	}
	var result ThreadsAccountResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("threads refresh token decode: %w", err)
	}
	return &result, nil
}

// Stats returns analytics stats for a Threads account.
func (a *ThreadsAPI) Stats(ctx context.Context, modelID string) (map[string]interface{}, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/threads/stats/"+modelID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("threads stats: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("threads stats decode: %w", err)
	}
	return result, nil
}
