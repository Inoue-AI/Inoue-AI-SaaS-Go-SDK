package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// JobsAPI provides access to the job management endpoints.
type JobsAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// JobResponse represents a job record returned by the API.
type JobResponse struct {
	ID                string                 `json:"id"`
	ModelID           string                 `json:"model_id"`
	JobType           string                 `json:"job_type"`
	Status            string                 `json:"status"`
	InputJSON         map[string]interface{} `json:"input_json"`
	ProgressJSON      map[string]interface{} `json:"progress_json"`
	ErrorJSON         map[string]interface{} `json:"error_json"`
	RequestedByUserID string                 `json:"requested_by_user_id"`
	ClaimedByWorkerID string                 `json:"claimed_by_worker_id"`
	CreatedAt         string                 `json:"created_at"`
	UpdatedAt         string                 `json:"updated_at"`
	StartedAt         string                 `json:"started_at"`
	FinishedAt        string                 `json:"finished_at"`
}

// AssetPublicResponse represents a public asset record returned by the API.
type AssetPublicResponse struct {
	ID          string `json:"id"`
	OwnerUserID string `json:"owner_user_id"`
	OwnerOrgID  string `json:"owner_org_id"`
	AssetType   string `json:"asset_type"`
	StorageKey  string `json:"storage_key"`
	StorageURL  string `json:"storage_url"`
	Filename    string `json:"filename"`
	MimeType    string `json:"mime_type"`
	SizeBytes   int64  `json:"size_bytes"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// JobStatusHistoryResponse represents a single status history entry for a job.
type JobStatusHistoryResponse struct {
	ID        string `json:"id"`
	JobID     string `json:"job_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// PipelineResponse represents a pipeline record returned by the API.
type PipelineResponse struct {
	ID        string                 `json:"id"`
	Status    string                 `json:"status"`
	Jobs      []JobResponse          `json:"jobs"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// List returns a paginated list of jobs.
func (a *JobsAPI) List(ctx context.Context, page, pageSize int) (*Page[JobResponse], error) {
	path := fmt.Sprintf("/v1/jobs?page=%d&page_size=%d", page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list jobs: %w", err)
	}
	var result Page[JobResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list jobs decode: %w", err)
	}
	return &result, nil
}

// Get retrieves a single job by ID.
func (a *JobsAPI) Get(ctx context.Context, jobID string) (*JobResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/jobs/"+jobID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get job: %w", err)
	}
	var result JobResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get job decode: %w", err)
	}
	return &result, nil
}

// Retry retries a failed job.
func (a *JobsAPI) Retry(ctx context.Context, jobID string) (*JobResponse, error) {
	body := map[string]string{"job_id": jobID}
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/jobs/retry", body, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("retry job: %w", err)
	}
	var result JobResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("retry job decode: %w", err)
	}
	return &result, nil
}

// History returns the status history for a job.
func (a *JobsAPI) History(ctx context.Context, jobID string) ([]JobStatusHistoryResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/jobs/history/"+jobID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("job history: %w", err)
	}
	var result []JobStatusHistoryResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("job history decode: %w", err)
	}
	return result, nil
}

// Outputs returns the output assets for a job.
func (a *JobsAPI) Outputs(ctx context.Context, jobID string) ([]AssetPublicResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/jobs/"+jobID+"/outputs", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("job outputs: %w", err)
	}
	var result []AssetPublicResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("job outputs decode: %w", err)
	}
	return result, nil
}

// OutputThumbnails returns thumbnail assets keyed by job ID.
func (a *JobsAPI) OutputThumbnails(ctx context.Context, jobIDs []string) (map[string]AssetPublicResponse, error) {
	body := map[string]interface{}{"job_ids": jobIDs}
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/jobs/outputs/thumbnails", body, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("job output thumbnails: %w", err)
	}
	var result map[string]AssetPublicResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("job output thumbnails decode: %w", err)
	}
	return result, nil
}

// Pipeline retrieves a pipeline by ID.
func (a *JobsAPI) Pipeline(ctx context.Context, pipelineID string) (*PipelineResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/jobs/pipelines/"+pipelineID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get pipeline: %w", err)
	}
	var result PipelineResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get pipeline decode: %w", err)
	}
	return &result, nil
}

// SaveOutputAsset saves a job output asset to the user's asset library.
func (a *JobsAPI) SaveOutputAsset(ctx context.Context, jobID, assetID string) (*AssetPublicResponse, error) {
	var apiResp ApiResponse
	path := fmt.Sprintf("/v1/jobs/%s/outputs/%s/save", jobID, assetID)
	if err := a.client.request(ctx, "POST", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("save job output asset: %w", err)
	}
	var result AssetPublicResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("save job output asset decode: %w", err)
	}
	return &result, nil
}
