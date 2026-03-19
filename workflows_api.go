package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// WorkflowsAPI provides access to the workflow management endpoints.
type WorkflowsAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// WorkflowCreateRequest is the body for creating a workflow.
type WorkflowCreateRequest struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description,omitempty"`
	Steps       []map[string]interface{} `json:"steps"`
	ModelID     string                   `json:"model_id,omitempty"`
	OwnerOrgID  string                   `json:"owner_org_id,omitempty"`
}

// WorkflowUpdateRequest is the body for updating a workflow.
type WorkflowUpdateRequest struct {
	Name        string                   `json:"name,omitempty"`
	Description string                   `json:"description,omitempty"`
	Steps       []map[string]interface{} `json:"steps,omitempty"`
}

// WorkflowRunRequest is the body for running a workflow.
type WorkflowRunRequest struct {
	WorkflowID string                 `json:"workflow_id"`
	Inputs     map[string]interface{} `json:"inputs,omitempty"`
}

// WorkflowBatchRunRequest is the body for batch-running a workflow.
type WorkflowBatchRunRequest struct {
	WorkflowID string                   `json:"workflow_id"`
	BatchInputs []map[string]interface{} `json:"batch_inputs"`
}

// WorkflowStepInputsUpdateRequest is the body for updating step inputs.
type WorkflowStepInputsUpdateRequest struct {
	Inputs map[string]interface{} `json:"inputs"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// WorkflowResponse represents a workflow record.
type WorkflowResponse struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Steps       []map[string]interface{} `json:"steps"`
	ModelID     string                   `json:"model_id"`
	OwnerUserID string                   `json:"owner_user_id"`
	OwnerOrgID  string                   `json:"owner_org_id"`
	Status      string                   `json:"status"`
	CreatedAt   string                   `json:"created_at"`
	UpdatedAt   string                   `json:"updated_at"`
}

// WorkflowRunResponse represents a workflow run record.
type WorkflowRunResponse struct {
	ID         string                   `json:"id"`
	WorkflowID string                   `json:"workflow_id"`
	Status     string                   `json:"status"`
	Inputs     map[string]interface{}   `json:"inputs"`
	Outputs    map[string]interface{}   `json:"outputs"`
	Steps      []map[string]interface{} `json:"steps"`
	Error      map[string]interface{}   `json:"error"`
	StartedAt  string                   `json:"started_at"`
	FinishedAt string                   `json:"finished_at"`
	CreatedAt  string                   `json:"created_at"`
	UpdatedAt  string                   `json:"updated_at"`
}

// WorkflowBatchRunResponse represents a batch workflow run record.
type WorkflowBatchRunResponse struct {
	ID         string                `json:"id"`
	WorkflowID string                `json:"workflow_id"`
	Status     string                `json:"status"`
	Runs       []WorkflowRunResponse `json:"runs"`
	CreatedAt  string                `json:"created_at"`
	UpdatedAt  string                `json:"updated_at"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// List returns a paginated list of workflows.
func (a *WorkflowsAPI) List(ctx context.Context, page, pageSize int) (*Page[WorkflowResponse], error) {
	path := fmt.Sprintf("/v1/workflows?page=%d&page_size=%d", page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list workflows: %w", err)
	}
	var result Page[WorkflowResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list workflows decode: %w", err)
	}
	return &result, nil
}

// Get retrieves a single workflow by ID.
func (a *WorkflowsAPI) Get(ctx context.Context, workflowID string) (*WorkflowResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/workflows/"+workflowID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get workflow: %w", err)
	}
	var result WorkflowResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get workflow decode: %w", err)
	}
	return &result, nil
}

// Create creates a new workflow.
func (a *WorkflowsAPI) Create(ctx context.Context, req WorkflowCreateRequest) (*WorkflowResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/workflows", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("create workflow: %w", err)
	}
	var result WorkflowResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create workflow decode: %w", err)
	}
	return &result, nil
}

// Update updates an existing workflow.
func (a *WorkflowsAPI) Update(ctx context.Context, workflowID string, req WorkflowUpdateRequest) (*WorkflowResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "PATCH", "/v1/workflows/"+workflowID, req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("update workflow: %w", err)
	}
	var result WorkflowResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("update workflow decode: %w", err)
	}
	return &result, nil
}

// Delete deletes a workflow by ID.
func (a *WorkflowsAPI) Delete(ctx context.Context, workflowID string) error {
	body := map[string]string{"workflow_id": workflowID}
	if err := a.client.request(ctx, "POST", "/v1/workflows/delete", body, nil, nil); err != nil {
		return fmt.Errorf("delete workflow: %w", err)
	}
	return nil
}

// Duplicate duplicates a workflow.
func (a *WorkflowsAPI) Duplicate(ctx context.Context, workflowID string) (*WorkflowResponse, error) {
	body := map[string]string{"workflow_id": workflowID}
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/workflows/duplicate", body, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("duplicate workflow: %w", err)
	}
	var result WorkflowResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("duplicate workflow decode: %w", err)
	}
	return &result, nil
}

// Run starts a workflow run.
func (a *WorkflowsAPI) Run(ctx context.Context, req WorkflowRunRequest) (*WorkflowRunResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/workflows/run", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("run workflow: %w", err)
	}
	var result WorkflowRunResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("run workflow decode: %w", err)
	}
	return &result, nil
}

// BatchRun starts a batch workflow run.
func (a *WorkflowsAPI) BatchRun(ctx context.Context, req WorkflowBatchRunRequest) (*WorkflowBatchRunResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/workflows/batch-run", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("batch run workflow: %w", err)
	}
	var result WorkflowBatchRunResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("batch run workflow decode: %w", err)
	}
	return &result, nil
}

// GetRun retrieves a workflow run by ID.
func (a *WorkflowsAPI) GetRun(ctx context.Context, runID string) (*WorkflowRunResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/workflows/runs/"+runID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get workflow run: %w", err)
	}
	var result WorkflowRunResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get workflow run decode: %w", err)
	}
	return &result, nil
}

// CancelRun cancels a running workflow run.
func (a *WorkflowsAPI) CancelRun(ctx context.Context, runID string) (*WorkflowRunResponse, error) {
	body := map[string]string{"run_id": runID}
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/workflows/runs/cancel", body, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("cancel workflow run: %w", err)
	}
	var result WorkflowRunResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("cancel workflow run decode: %w", err)
	}
	return &result, nil
}

// ListRuns returns a paginated list of workflow runs.
func (a *WorkflowsAPI) ListRuns(ctx context.Context, workflowID string, page, pageSize int) (*Page[WorkflowRunResponse], error) {
	path := fmt.Sprintf("/v1/workflows/%s/runs?page=%d&page_size=%d", workflowID, page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list workflow runs: %w", err)
	}
	var result Page[WorkflowRunResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list workflow runs decode: %w", err)
	}
	return &result, nil
}

// GetBatchRun retrieves a batch workflow run by ID.
func (a *WorkflowsAPI) GetBatchRun(ctx context.Context, batchRunID string) (*WorkflowBatchRunResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/workflows/batch-runs/"+batchRunID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get batch workflow run: %w", err)
	}
	var result WorkflowBatchRunResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get batch workflow run decode: %w", err)
	}
	return &result, nil
}

// ApproveStep approves a pending workflow step.
func (a *WorkflowsAPI) ApproveStep(ctx context.Context, runID, stepID string) (*WorkflowRunResponse, error) {
	body := map[string]string{"run_id": runID, "step_id": stepID}
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/workflows/runs/steps/approve", body, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("approve workflow step: %w", err)
	}
	var result WorkflowRunResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("approve workflow step decode: %w", err)
	}
	return &result, nil
}

// RejectStep rejects a pending workflow step.
func (a *WorkflowsAPI) RejectStep(ctx context.Context, runID, stepID string) (*WorkflowRunResponse, error) {
	body := map[string]string{"run_id": runID, "step_id": stepID}
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/workflows/runs/steps/reject", body, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("reject workflow step: %w", err)
	}
	var result WorkflowRunResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("reject workflow step decode: %w", err)
	}
	return &result, nil
}

// UpdateStepInputs updates the inputs for a workflow step.
func (a *WorkflowsAPI) UpdateStepInputs(ctx context.Context, runID, stepID string, req WorkflowStepInputsUpdateRequest) (*WorkflowRunResponse, error) {
	var apiResp ApiResponse
	path := fmt.Sprintf("/v1/workflows/runs/%s/steps/%s/inputs", runID, stepID)
	if err := a.client.request(ctx, "PUT", path, req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("update workflow step inputs: %w", err)
	}
	var result WorkflowRunResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("update workflow step inputs decode: %w", err)
	}
	return &result, nil
}
