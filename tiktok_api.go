package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// TikTokAPI provides access to the TikTok integration endpoints.
type TikTokAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// TikTokConnectRequest is the body for initiating a TikTok connection.
type TikTokConnectRequest struct {
	ModelID    string `json:"model_id"`
	RedirectURI string `json:"redirect_uri,omitempty"`
}

// TikTokCallbackRequest is the body for the OAuth callback.
type TikTokCallbackRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

// TikTokPostRequest is the body for creating a TikTok post.
type TikTokPostRequest struct {
	ModelID     string                 `json:"model_id"`
	Caption     string                 `json:"caption,omitempty"`
	AssetIDs    []string               `json:"asset_ids,omitempty"`
	Privacy     string                 `json:"privacy,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// TikTokConnectResponse represents a TikTok connection response.
type TikTokConnectResponse struct {
	AuthURL string `json:"auth_url"`
	State   string `json:"state"`
}

// TikTokAccountResponse represents a connected TikTok account.
type TikTokAccountResponse struct {
	ID          string `json:"id"`
	ModelID     string `json:"model_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Followers   int    `json:"followers"`
	Connected   bool   `json:"connected"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// TikTokPostResponse represents a TikTok post result.
type TikTokPostResponse struct {
	ID        string                 `json:"id"`
	PostID    string                 `json:"post_id"`
	Status    string                 `json:"status"`
	Caption   string                 `json:"caption"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt string                 `json:"created_at"`
}

// TikTokVideoResponse represents a TikTok video.
type TikTokVideoResponse struct {
	ID        string                 `json:"id"`
	VideoID   string                 `json:"video_id"`
	Caption   string                 `json:"caption"`
	URL       string                 `json:"url"`
	Stats     map[string]interface{} `json:"stats"`
	CreatedAt string                 `json:"created_at"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// Connect initiates a TikTok OAuth connection.
func (a *TikTokAPI) Connect(ctx context.Context, req TikTokConnectRequest) (*TikTokConnectResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/tiktok/connect", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("tiktok connect: %w", err)
	}
	var result TikTokConnectResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("tiktok connect decode: %w", err)
	}
	return &result, nil
}

// Callback handles the TikTok OAuth callback.
func (a *TikTokAPI) Callback(ctx context.Context, req TikTokCallbackRequest) (*TikTokAccountResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/tiktok/callback", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("tiktok callback: %w", err)
	}
	var result TikTokAccountResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("tiktok callback decode: %w", err)
	}
	return &result, nil
}

// Accounts returns the list of connected TikTok accounts.
func (a *TikTokAPI) Accounts(ctx context.Context) ([]TikTokAccountResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/tiktok/accounts", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list tiktok accounts: %w", err)
	}
	var result []TikTokAccountResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list tiktok accounts decode: %w", err)
	}
	return result, nil
}

// Account retrieves a specific TikTok account by model ID.
func (a *TikTokAPI) Account(ctx context.Context, modelID string) (*TikTokAccountResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/tiktok/accounts/"+modelID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get tiktok account: %w", err)
	}
	var result TikTokAccountResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get tiktok account decode: %w", err)
	}
	return &result, nil
}

// Disconnect disconnects a TikTok account.
func (a *TikTokAPI) Disconnect(ctx context.Context, modelID string) error {
	body := map[string]string{"model_id": modelID}
	if err := a.client.request(ctx, "POST", "/v1/tiktok/disconnect", body, nil, nil); err != nil {
		return fmt.Errorf("tiktok disconnect: %w", err)
	}
	return nil
}

// Post creates a TikTok post.
func (a *TikTokAPI) Post(ctx context.Context, req TikTokPostRequest) (*TikTokPostResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/tiktok/post", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("tiktok post: %w", err)
	}
	var result TikTokPostResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("tiktok post decode: %w", err)
	}
	return &result, nil
}

// Videos returns the list of TikTok videos for a model.
func (a *TikTokAPI) Videos(ctx context.Context, modelID string, page, pageSize int) (*Page[TikTokVideoResponse], error) {
	path := fmt.Sprintf("/v1/tiktok/videos/%s?page=%d&page_size=%d", modelID, page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list tiktok videos: %w", err)
	}
	var result Page[TikTokVideoResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list tiktok videos decode: %w", err)
	}
	return &result, nil
}

// RefreshToken refreshes the TikTok access token for a model.
func (a *TikTokAPI) RefreshToken(ctx context.Context, modelID string) (*TikTokAccountResponse, error) {
	body := map[string]string{"model_id": modelID}
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/tiktok/refresh", body, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("tiktok refresh token: %w", err)
	}
	var result TikTokAccountResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("tiktok refresh token decode: %w", err)
	}
	return &result, nil
}

// Stats returns analytics stats for a TikTok account.
func (a *TikTokAPI) Stats(ctx context.Context, modelID string) (map[string]interface{}, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/tiktok/stats/"+modelID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("tiktok stats: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("tiktok stats decode: %w", err)
	}
	return result, nil
}
