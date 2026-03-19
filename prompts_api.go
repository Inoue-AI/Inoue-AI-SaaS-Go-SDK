package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// PromptsAPI provides access to the prompt template, version, and run endpoints.
type PromptsAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// PromptTemplateCreateRequest is the body for POST /v1/prompts.
type PromptTemplateCreateRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Category    string                 `json:"category,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Variables   []map[string]interface{} `json:"variables,omitempty"`
	Body        string                 `json:"body,omitempty"`
	ModelID     string                 `json:"model_id,omitempty"`
	OwnerOrgID  string                 `json:"owner_org_id,omitempty"`
}

// PromptTemplateUpdateRequest is the body for PATCH /v1/prompts/template/{id}.
type PromptTemplateUpdateRequest struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Category    string                 `json:"category,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Variables   []map[string]interface{} `json:"variables,omitempty"`
	Body        string                 `json:"body,omitempty"`
}

// PromptTemplateDuplicateRequest is the body for POST /v1/prompts/template/duplicate.
type PromptTemplateDuplicateRequest struct {
	TemplateID string `json:"template_id"`
	Name       string `json:"name,omitempty"`
}

// PromptTemplateExamplesUpdateRequest is the body for PUT /v1/prompts/template/{id}/examples.
type PromptTemplateExamplesUpdateRequest struct {
	Examples []map[string]interface{} `json:"examples"`
}

// PromptTemplateModelLinkRequest is the body for POST /v1/prompts/link.
type PromptTemplateModelLinkRequest struct {
	TemplateID string `json:"template_id"`
	ModelID    string `json:"model_id"`
}

// PromptVersionCreateRequest is the body for POST /v1/prompts/version.
type PromptVersionCreateRequest struct {
	TemplateID string                 `json:"template_id"`
	Body       string                 `json:"body"`
	Variables  []map[string]interface{} `json:"variables,omitempty"`
	Notes      string                 `json:"notes,omitempty"`
}

// PromptRunCreateRequest is the body for POST /v1/prompts/run.
type PromptRunCreateRequest struct {
	TemplateID string                 `json:"template_id,omitempty"`
	VersionID  string                 `json:"version_id,omitempty"`
	ModelID    string                 `json:"model_id,omitempty"`
	Inputs     map[string]interface{} `json:"inputs,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// PromptTemplateResponse represents a prompt template returned by the API.
type PromptTemplateResponse struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Category    string                   `json:"category"`
	Tags        []string                 `json:"tags"`
	Body        string                   `json:"body"`
	Variables   []map[string]interface{} `json:"variables"`
	Examples    []map[string]interface{} `json:"examples"`
	Locked      bool                     `json:"locked"`
	OwnerUserID string                   `json:"owner_user_id"`
	OwnerOrgID  string                   `json:"owner_org_id"`
	CreatedAt   string                   `json:"created_at"`
	UpdatedAt   string                   `json:"updated_at"`
}

// PromptTemplateModelLinkResponse represents a link between a template and a model.
type PromptTemplateModelLinkResponse struct {
	ID         string `json:"id"`
	TemplateID string `json:"template_id"`
	ModelID    string `json:"model_id"`
	CreatedAt  string `json:"created_at"`
}

// PromptVersionResponse represents a prompt version returned by the API.
type PromptVersionResponse struct {
	ID         string                   `json:"id"`
	TemplateID string                   `json:"template_id"`
	Version    int                      `json:"version"`
	Body       string                   `json:"body"`
	Variables  []map[string]interface{} `json:"variables"`
	Notes      string                   `json:"notes"`
	CreatedAt  string                   `json:"created_at"`
}

// PromptRunResponse represents a prompt run record returned by the API.
type PromptRunResponse struct {
	ID         string                 `json:"id"`
	TemplateID string                 `json:"template_id"`
	VersionID  string                 `json:"version_id"`
	ModelID    string                 `json:"model_id"`
	Inputs     map[string]interface{} `json:"inputs"`
	Output     string                 `json:"output"`
	Status     string                 `json:"status"`
	CreatedAt  string                 `json:"created_at"`
}

// PromptRunBundleResponse represents the full result of running a prompt.
type PromptRunBundleResponse struct {
	Run      PromptRunResponse      `json:"run"`
	Output   string                 `json:"output"`
	Metadata map[string]interface{} `json:"metadata"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// ListTemplates returns a paginated list of prompt templates.
func (a *PromptsAPI) ListTemplates(ctx context.Context, page, pageSize int) (*Page[PromptTemplateResponse], error) {
	path := fmt.Sprintf("/v1/prompts/templates?page=%d&page_size=%d", page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list prompt templates: %w", err)
	}
	var result Page[PromptTemplateResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list prompt templates decode: %w", err)
	}
	return &result, nil
}

// GetTemplate retrieves a single prompt template by ID.
func (a *PromptsAPI) GetTemplate(ctx context.Context, templateID string) (*PromptTemplateResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/prompts/template/"+templateID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get prompt template: %w", err)
	}
	var result PromptTemplateResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get prompt template decode: %w", err)
	}
	return &result, nil
}

// CreateTemplate creates a new prompt template.
func (a *PromptsAPI) CreateTemplate(ctx context.Context, req PromptTemplateCreateRequest) (*PromptTemplateResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/prompts", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("create prompt template: %w", err)
	}
	var result PromptTemplateResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create prompt template decode: %w", err)
	}
	return &result, nil
}

// UpdateTemplate updates an existing prompt template.
func (a *PromptsAPI) UpdateTemplate(ctx context.Context, req PromptTemplateUpdateRequest) (*PromptTemplateResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "PATCH", "/v1/prompts/template/"+req.ID, req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("update prompt template: %w", err)
	}
	var result PromptTemplateResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("update prompt template decode: %w", err)
	}
	return &result, nil
}

// DeleteTemplate deletes a prompt template by ID.
func (a *PromptsAPI) DeleteTemplate(ctx context.Context, templateID string) error {
	if err := a.client.request(ctx, "DELETE", "/v1/prompts/template/"+templateID, nil, nil, nil); err != nil {
		return fmt.Errorf("delete prompt template: %w", err)
	}
	return nil
}

// DuplicateTemplate duplicates an existing prompt template.
func (a *PromptsAPI) DuplicateTemplate(ctx context.Context, req PromptTemplateDuplicateRequest) (*PromptTemplateResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/prompts/template/duplicate", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("duplicate prompt template: %w", err)
	}
	var result PromptTemplateResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("duplicate prompt template decode: %w", err)
	}
	return &result, nil
}

// LockTemplate locks a prompt template to prevent edits.
func (a *PromptsAPI) LockTemplate(ctx context.Context, templateID string) (*PromptTemplateResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/prompts/template/"+templateID+"/lock", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("lock prompt template: %w", err)
	}
	var result PromptTemplateResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("lock prompt template decode: %w", err)
	}
	return &result, nil
}

// UnlockTemplate unlocks a prompt template to allow edits.
func (a *PromptsAPI) UnlockTemplate(ctx context.Context, templateID string) (*PromptTemplateResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/prompts/template/"+templateID+"/unlock", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("unlock prompt template: %w", err)
	}
	var result PromptTemplateResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("unlock prompt template decode: %w", err)
	}
	return &result, nil
}

// SetTemplateExamples replaces all examples on a prompt template.
func (a *PromptsAPI) SetTemplateExamples(ctx context.Context, templateID string, req PromptTemplateExamplesUpdateRequest) (*PromptTemplateResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "PUT", "/v1/prompts/template/"+templateID+"/examples", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("set prompt template examples: %w", err)
	}
	var result PromptTemplateResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("set prompt template examples decode: %w", err)
	}
	return &result, nil
}

// ListVersions returns a paginated list of versions for a prompt template.
func (a *PromptsAPI) ListVersions(ctx context.Context, templateID string, page, pageSize int) (*Page[PromptVersionResponse], error) {
	path := fmt.Sprintf("/v1/prompts/template/%s/versions?page=%d&page_size=%d", templateID, page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list prompt versions: %w", err)
	}
	var result Page[PromptVersionResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list prompt versions decode: %w", err)
	}
	return &result, nil
}

// GetVersion retrieves a single prompt version by ID.
func (a *PromptsAPI) GetVersion(ctx context.Context, versionID string) (*PromptVersionResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/prompts/version/"+versionID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get prompt version: %w", err)
	}
	var result PromptVersionResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get prompt version decode: %w", err)
	}
	return &result, nil
}

// CreateVersion creates a new version for a prompt template.
func (a *PromptsAPI) CreateVersion(ctx context.Context, req PromptVersionCreateRequest) (*PromptVersionResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/prompts/version", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("create prompt version: %w", err)
	}
	var result PromptVersionResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create prompt version decode: %w", err)
	}
	return &result, nil
}

// LinkTemplate links a prompt template to a model.
func (a *PromptsAPI) LinkTemplate(ctx context.Context, req PromptTemplateModelLinkRequest) (*PromptTemplateModelLinkResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/prompts/link", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("link prompt template: %w", err)
	}
	var result PromptTemplateModelLinkResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("link prompt template decode: %w", err)
	}
	return &result, nil
}

// UnlinkTemplate unlinks a prompt template from a model.
func (a *PromptsAPI) UnlinkTemplate(ctx context.Context, templateID, modelID string) error {
	body := map[string]string{
		"template_id": templateID,
		"model_id":    modelID,
	}
	if err := a.client.request(ctx, "POST", "/v1/prompts/unlink", body, nil, nil); err != nil {
		return fmt.Errorf("unlink prompt template: %w", err)
	}
	return nil
}

// Run executes a prompt and returns the run result.
func (a *PromptsAPI) Run(ctx context.Context, req PromptRunCreateRequest) (*PromptRunBundleResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/prompts/run", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("run prompt: %w", err)
	}
	var result PromptRunBundleResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("run prompt decode: %w", err)
	}
	return &result, nil
}

// ListRuns returns a paginated list of prompt runs.
func (a *PromptsAPI) ListRuns(ctx context.Context, page, pageSize int) (*Page[PromptRunResponse], error) {
	path := fmt.Sprintf("/v1/prompts/runs?page=%d&page_size=%d", page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list prompt runs: %w", err)
	}
	var result Page[PromptRunResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list prompt runs decode: %w", err)
	}
	return &result, nil
}

// ModelRuns returns a paginated list of prompt runs for a specific model.
func (a *PromptsAPI) ModelRuns(ctx context.Context, modelID string, page, pageSize int) (*Page[PromptRunResponse], error) {
	path := fmt.Sprintf("/v1/prompts/models/%s/runs?page=%d&page_size=%d", modelID, page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list model prompt runs: %w", err)
	}
	var result Page[PromptRunResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list model prompt runs decode: %w", err)
	}
	return &result, nil
}

// DeleteRun deletes a prompt run by ID.
func (a *PromptsAPI) DeleteRun(ctx context.Context, runID string) error {
	body := map[string]string{"run_id": runID}
	if err := a.client.request(ctx, "POST", "/v1/prompts/run/delete", body, nil, nil); err != nil {
		return fmt.Errorf("delete prompt run: %w", err)
	}
	return nil
}
