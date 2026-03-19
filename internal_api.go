package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// InternalAPI provides access to internal/worker management endpoints.
// These endpoints are used by worker processes to register themselves,
// obtain authentication tokens, and report their status.
type InternalAPI struct {
	client *InoueClient
}

// registerWorkerRequest is the JSON body for POST /internal/workers/register.
type registerWorkerRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// RegisterWorker registers this worker with the backend.
// The operation is idempotent: a 409 Conflict response (already registered) is
// treated as success and returns nil.
func (a *InternalAPI) RegisterWorker(ctx context.Context, workerID, name, adminToken string) error {
	body := registerWorkerRequest{ID: workerID, Name: name}
	headers := map[string]string{"X-Admin-Token": adminToken}

	err := a.client.request(ctx, "POST", "/internal/workers/register", body, nil, headers)
	if err != nil {
		if sdkErr, ok := err.(*SdkError); ok && sdkErr.Status == 409 {
			return nil
		}
		return fmt.Errorf("register worker: %w", err)
	}
	return nil
}

// workerTokenRequest is the JSON body for POST /internal/workers/token.
type workerTokenRequest struct {
	WorkerID string `json:"worker_id"`
}

// workerTokenResponse is the JSON body returned by the token endpoint.
type workerTokenResponse struct {
	WorkerToken string `json:"worker_token"`
}

// WorkerToken obtains a JWT for the given worker using the bootstrap secret.
// The returned string is a Bearer token suitable for Authorization headers.
func (a *InternalAPI) WorkerToken(ctx context.Context, workerID, bootstrapSecret string) (string, error) {
	body := workerTokenRequest{WorkerID: workerID}
	headers := map[string]string{"X-Worker-Bootstrap": bootstrapSecret}

	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/internal/workers/token", body, &apiResp, headers)
	if err != nil {
		return "", fmt.Errorf("worker token: %w", err)
	}

	var tokenResp workerTokenResponse
	if err := json.Unmarshal(apiResp.Data, &tokenResp); err != nil {
		return "", fmt.Errorf("worker token: failed to decode token response: %w", err)
	}

	return tokenResp.WorkerToken, nil
}

// setWorkerStatusRequest is the JSON body for the worker status endpoint.
type setWorkerStatusRequest struct {
	Status string `json:"status"`
}

// SetWorkerStatus updates this worker's status (e.g. "active", "disabled").
// The request uses the client's current Bearer token for authentication.
func (a *InternalAPI) SetWorkerStatus(ctx context.Context, status string) error {
	body := setWorkerStatusRequest{Status: status}

	err := a.client.request(ctx, "POST", "/internal/workers/me/status", body, nil, nil)
	if err != nil {
		return fmt.Errorf("set worker status: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Transcription Request Lifecycle
// ---------------------------------------------------------------------------
// These methods mirror the download lifecycle pattern in the Python SDK:
//   get_download, claim_download, reclaim_download, download_heartbeat,
//   download_progress, complete_download, fail_download
// All routes live under /internal/transcription-requests/{id}/*.

// TranscriptionRequest represents a transcription request record from the backend.
type TranscriptionRequest struct {
	ID               string                 `json:"id"`
	CaptionProjectID string                 `json:"caption_project_id"`
	AssetID          string                 `json:"asset_id"`
	Status           string                 `json:"status"`
	LanguageHint     string                 `json:"language_hint,omitempty"`
	Provider         string                 `json:"provider"`
	ClaimedByWorker  string                 `json:"claimed_by_worker_id,omitempty"`
	ClaimedAt        string                 `json:"claimed_at,omitempty"`
	StartedAt        string                 `json:"started_at,omitempty"`
	FinishedAt       string                 `json:"finished_at,omitempty"`
	ProgressJSON     map[string]interface{} `json:"progress_json,omitempty"`
	ErrorJSON        map[string]interface{} `json:"error_json,omitempty"`
	CreatedAt        string                 `json:"created_at"`
	UpdatedAt        string                 `json:"updated_at"`
}

// TranscriptionClaimResponse is the response body from the claim/reclaim endpoints.
type TranscriptionClaimResponse struct {
	ID              string `json:"id"`
	Status          string `json:"status"`
	ClaimedByWorker string `json:"claimed_by_worker_id"`
	ClaimedAt       string `json:"claimed_at"`
}

// TranscriptionProgressRequest is the JSON body for the progress endpoint.
type TranscriptionProgressRequest struct {
	ProgressJSON map[string]interface{} `json:"progress_json"`
}

// TranscriptionCompleteRequest is the JSON body for the complete endpoint.
type TranscriptionCompleteRequest struct {
	WordsJSON        []TranscriptionWord    `json:"words_json"`
	ChunksJSON       []map[string]interface{} `json:"chunks_json"`
	LanguageDetected string                 `json:"language_detected"`
	DurationSeconds  float64                `json:"duration_seconds"`
}

// TranscriptionWord represents a single word with timing from the transcription.
type TranscriptionWord struct {
	Text       string  `json:"text"`
	StartTime  float64 `json:"start_time"`
	EndTime    float64 `json:"end_time"`
	Confidence float64 `json:"confidence"`
	Speaker    string  `json:"speaker,omitempty"`
}

// TranscriptionFailRequest is the JSON body for the fail endpoint.
type TranscriptionFailRequest struct {
	ErrorJSON map[string]interface{} `json:"error_json"`
}

// GetTranscription fetches a transcription request by ID.
func (a *InternalAPI) GetTranscription(ctx context.Context, requestID string) (*TranscriptionRequest, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", "/internal/transcription-requests/"+requestID, nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("get transcription: %w", err)
	}

	var req TranscriptionRequest
	if err := json.Unmarshal(apiResp.Data, &req); err != nil {
		return nil, fmt.Errorf("get transcription: failed to decode response: %w", err)
	}
	return &req, nil
}

// ClaimTranscription atomically claims a queued transcription request for this worker.
func (a *InternalAPI) ClaimTranscription(ctx context.Context, requestID string) (*TranscriptionClaimResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/internal/transcription-requests/"+requestID+"/claim", nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("claim transcription: %w", err)
	}

	var resp TranscriptionClaimResponse
	if err := json.Unmarshal(apiResp.Data, &resp); err != nil {
		return nil, fmt.Errorf("claim transcription: failed to decode response: %w", err)
	}
	return &resp, nil
}

// ReclaimTranscription reclaims a running transcription request whose previous worker died.
func (a *InternalAPI) ReclaimTranscription(ctx context.Context, requestID string) (*TranscriptionClaimResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/internal/transcription-requests/"+requestID+"/reclaim", nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("reclaim transcription: %w", err)
	}

	var resp TranscriptionClaimResponse
	if err := json.Unmarshal(apiResp.Data, &resp); err != nil {
		return nil, fmt.Errorf("reclaim transcription: failed to decode response: %w", err)
	}
	return &resp, nil
}

// TranscriptionHeartbeat sends a liveness ping for a running transcription request.
func (a *InternalAPI) TranscriptionHeartbeat(ctx context.Context, requestID string) error {
	err := a.client.request(ctx, "POST", "/internal/transcription-requests/"+requestID+"/heartbeat", nil, nil, nil)
	if err != nil {
		return fmt.Errorf("transcription heartbeat: %w", err)
	}
	return nil
}

// TranscriptionProgress reports progress for a running transcription request.
func (a *InternalAPI) TranscriptionProgress(ctx context.Context, requestID string, req TranscriptionProgressRequest) error {
	err := a.client.request(ctx, "POST", "/internal/transcription-requests/"+requestID+"/progress", req, nil, nil)
	if err != nil {
		return fmt.Errorf("transcription progress: %w", err)
	}
	return nil
}

// CompleteTranscription marks a transcription request as succeeded with the transcription result.
func (a *InternalAPI) CompleteTranscription(ctx context.Context, requestID string, req TranscriptionCompleteRequest) error {
	err := a.client.request(ctx, "POST", "/internal/transcription-requests/"+requestID+"/complete", req, nil, nil)
	if err != nil {
		return fmt.Errorf("complete transcription: %w", err)
	}
	return nil
}

// FailTranscription marks a transcription request as failed with error details.
func (a *InternalAPI) FailTranscription(ctx context.Context, requestID string, req TranscriptionFailRequest) error {
	err := a.client.request(ctx, "POST", "/internal/transcription-requests/"+requestID+"/fail", req, nil, nil)
	if err != nil {
		return fmt.Errorf("fail transcription: %w", err)
	}
	return nil
}
