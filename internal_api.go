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

// registerWorkerRequest is the JSON body for POST /internal/workers.
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

	err := a.client.request(ctx, "POST", "/internal/workers", body, nil, headers)
	if err != nil {
		if sdkErr, ok := err.(*SdkError); ok && sdkErr.Status == 409 {
			return nil
		}
		return fmt.Errorf("register worker: %w", err)
	}
	return nil
}

// workerTokenResponse is the JSON body returned by the token endpoint.
type workerTokenResponse struct {
	Token string `json:"token"`
}

// WorkerToken obtains a JWT for the given worker using the bootstrap secret.
// The returned string is a Bearer token suitable for Authorization headers.
func (a *InternalAPI) WorkerToken(ctx context.Context, workerID, bootstrapSecret string) (string, error) {
	path := fmt.Sprintf("/internal/workers/%s/token", workerID)
	headers := map[string]string{"X-Worker-Bootstrap": bootstrapSecret}

	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", path, nil, &apiResp, headers)
	if err != nil {
		return "", fmt.Errorf("worker token: %w", err)
	}

	var tokenResp workerTokenResponse
	if err := json.Unmarshal(apiResp.Data, &tokenResp); err != nil {
		return "", fmt.Errorf("worker token: failed to decode token response: %w", err)
	}

	return tokenResp.Token, nil
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
