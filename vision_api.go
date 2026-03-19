package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// VisionAPI provides access to the vision and image/video generation endpoints.
type VisionAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// VisionJobRequest is a generic request for vision job endpoints.
type VisionJobRequest struct {
	ModelID    string                 `json:"model_id,omitempty"`
	AssetID    string                 `json:"asset_id,omitempty"`
	AssetIDs   []string               `json:"asset_ids,omitempty"`
	Prompt     string                 `json:"prompt,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	OwnerOrgID string                 `json:"owner_org_id,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// VisionJobResponse represents a vision job result.
type VisionJobResponse struct {
	JobID     string                 `json:"job_id"`
	Status    string                 `json:"status"`
	Output    map[string]interface{} `json:"output"`
	CreatedAt string                 `json:"created_at"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// FaceSwap creates a face swap job.
func (a *VisionAPI) FaceSwap(ctx context.Context, req VisionJobRequest) (*VisionJobResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/vision/face-swap", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("face swap: %w", err)
	}
	var result VisionJobResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("face swap decode: %w", err)
	}
	return &result, nil
}

// ImageCrop creates an image crop job.
func (a *VisionAPI) ImageCrop(ctx context.Context, req VisionJobRequest) (*VisionJobResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/vision/image-crop", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("image crop: %w", err)
	}
	var result VisionJobResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("image crop decode: %w", err)
	}
	return &result, nil
}

// Seedream creates a Seedream generation job.
func (a *VisionAPI) Seedream(ctx context.Context, req VisionJobRequest) (*VisionJobResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/vision/seedream", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("seedream: %w", err)
	}
	var result VisionJobResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("seedream decode: %w", err)
	}
	return &result, nil
}

// Topaz creates a Topaz upscale job.
func (a *VisionAPI) Topaz(ctx context.Context, req VisionJobRequest) (*VisionJobResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/vision/topaz", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("topaz: %w", err)
	}
	var result VisionJobResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("topaz decode: %w", err)
	}
	return &result, nil
}

// NanoBanana creates a NanoBanana generation job.
func (a *VisionAPI) NanoBanana(ctx context.Context, req VisionJobRequest) (*VisionJobResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/vision/nanobanana", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("nanobanana: %w", err)
	}
	var result VisionJobResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("nanobanana decode: %w", err)
	}
	return &result, nil
}

// Sora creates a Sora video generation job.
func (a *VisionAPI) Sora(ctx context.Context, req VisionJobRequest) (*VisionJobResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/vision/sora", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("sora: %w", err)
	}
	var result VisionJobResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("sora decode: %w", err)
	}
	return &result, nil
}

// Kling creates a Kling video generation job.
func (a *VisionAPI) Kling(ctx context.Context, req VisionJobRequest) (*VisionJobResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/vision/kling", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("kling: %w", err)
	}
	var result VisionJobResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("kling decode: %w", err)
	}
	return &result, nil
}

// Grok creates a Grok generation job.
func (a *VisionAPI) Grok(ctx context.Context, req VisionJobRequest) (*VisionJobResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/vision/grok", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("grok: %w", err)
	}
	var result VisionJobResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("grok decode: %w", err)
	}
	return &result, nil
}

// Flux creates a Flux generation job.
func (a *VisionAPI) Flux(ctx context.Context, req VisionJobRequest) (*VisionJobResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/vision/flux", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("flux: %w", err)
	}
	var result VisionJobResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("flux decode: %w", err)
	}
	return &result, nil
}
