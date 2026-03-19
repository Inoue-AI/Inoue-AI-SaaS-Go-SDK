package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// PostsAPI provides access to the social media post endpoints.
type PostsAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// PostCreateRequest is the body for creating a post.
type PostCreateRequest struct {
	ModelID    string                 `json:"model_id"`
	PlatformID string                `json:"platform_id"`
	Caption    string                 `json:"caption,omitempty"`
	AssetIDs   []string               `json:"asset_ids,omitempty"`
	ScheduleAt string                 `json:"schedule_at,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// PostPolicySetRequest is the body for setting a post policy.
type PostPolicySetRequest struct {
	ModelID       string `json:"model_id"`
	RequireReview bool   `json:"require_review"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// PostResponse represents a post record.
type PostResponse struct {
	ID         string                 `json:"id"`
	ModelID    string                 `json:"model_id"`
	PlatformID string                `json:"platform_id"`
	Caption    string                 `json:"caption"`
	Status     string                 `json:"status"`
	AssetIDs   []string               `json:"asset_ids"`
	Metadata   map[string]interface{} `json:"metadata"`
	CreatedAt  string                 `json:"created_at"`
	UpdatedAt  string                 `json:"updated_at"`
}

// PostPolicyResponse represents a post policy configuration.
type PostPolicyResponse struct {
	ModelID       string `json:"model_id"`
	RequireReview bool   `json:"require_review"`
	UpdatedAt     string `json:"updated_at"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// GetPolicy retrieves the post policy for a model.
func (a *PostsAPI) GetPolicy(ctx context.Context, modelID string) (*PostPolicyResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/posts/policy/"+modelID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get post policy: %w", err)
	}
	var result PostPolicyResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get post policy decode: %w", err)
	}
	return &result, nil
}

// SetPolicy sets the post policy for a model.
func (a *PostsAPI) SetPolicy(ctx context.Context, req PostPolicySetRequest) (*PostPolicyResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/posts/policy", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("set post policy: %w", err)
	}
	var result PostPolicyResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("set post policy decode: %w", err)
	}
	return &result, nil
}

// Create creates a new post.
func (a *PostsAPI) Create(ctx context.Context, req PostCreateRequest) (*PostResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/posts", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}
	var result PostResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create post decode: %w", err)
	}
	return &result, nil
}

// Approve approves a post for publication.
func (a *PostsAPI) Approve(ctx context.Context, postID string) (*PostResponse, error) {
	body := map[string]string{"post_id": postID}
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/posts/approve", body, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("approve post: %w", err)
	}
	var result PostResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("approve post decode: %w", err)
	}
	return &result, nil
}

// Publish publishes a post immediately.
func (a *PostsAPI) Publish(ctx context.Context, postID string) (*PostResponse, error) {
	body := map[string]string{"post_id": postID}
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/posts/publish", body, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("publish post: %w", err)
	}
	var result PostResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("publish post decode: %w", err)
	}
	return &result, nil
}
