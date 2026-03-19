package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// ---------------------------------------------------------------------------
// Model types
// ---------------------------------------------------------------------------

// ModelResponse represents a model record returned by the API.
type ModelResponse struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name,omitempty"`
	OwnerUserID       string                 `json:"owner_user_id,omitempty"`
	OwnerOrgID        string                 `json:"owner_org_id,omitempty"`
	Status            string                 `json:"status,omitempty"`
	CreationMode      string                 `json:"creation_mode,omitempty"`
	CanonicalImageURL string                 `json:"canonical_image_url,omitempty"`
	TriggerWord       string                 `json:"trigger_word,omitempty"`
	BaseModel         string                 `json:"base_model,omitempty"`
	InputJSON         map[string]interface{} `json:"input_json,omitempty"`
	ConfigJSON        map[string]interface{} `json:"config_json,omitempty"`
	CreatedAt         string                 `json:"created_at,omitempty"`
	UpdatedAt         string                 `json:"updated_at,omitempty"`
	DeletedAt         string                 `json:"deleted_at,omitempty"`
}

// ModelListParams contains query parameters for listing models.
type ModelListParams struct {
	Page       int    `json:"page,omitempty"`
	PageSize   int    `json:"page_size,omitempty"`
	OwnerOrgID string `json:"owner_org_id,omitempty"`
	Status     string `json:"status,omitempty"`
	Scope      string `json:"scope,omitempty"`
}

// ModelCreateRequest is the JSON body for POST /v1/models.
type ModelCreateRequest struct {
	Name         string                 `json:"name"`
	OwnerUserID  string                 `json:"owner_user_id,omitempty"`
	OwnerOrgID   string                 `json:"owner_org_id,omitempty"`
	CreationMode string                 `json:"creation_mode,omitempty"`
	TriggerWord  string                 `json:"trigger_word,omitempty"`
	BaseModel    string                 `json:"base_model,omitempty"`
	InputJSON    map[string]interface{} `json:"input_json,omitempty"`
	ConfigJSON   map[string]interface{} `json:"config_json,omitempty"`
}

// ModelUpdateRequest is the JSON body for POST /v1/models/update.
type ModelUpdateRequest struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name,omitempty"`
	TriggerWord string                 `json:"trigger_word,omitempty"`
	InputJSON   map[string]interface{} `json:"input_json,omitempty"`
	ConfigJSON  map[string]interface{} `json:"config_json,omitempty"`
}

// ModelPatchRequest is the JSON body for PATCH /v1/models/{id}.
type ModelPatchRequest struct {
	Name              *string                `json:"name,omitempty"`
	Status            *string                `json:"status,omitempty"`
	TriggerWord       *string                `json:"trigger_word,omitempty"`
	CanonicalImageURL *string                `json:"canonical_image_url,omitempty"`
	InputJSON         map[string]interface{} `json:"input_json,omitempty"`
	ConfigJSON        map[string]interface{} `json:"config_json,omitempty"`
}

// ModelForkRequest is the JSON body for POST /v1/models/fork.
type ModelForkRequest struct {
	ModelID    string `json:"model_id"`
	OwnerOrgID string `json:"owner_org_id,omitempty"`
}

// ModelTransferRequest is the JSON body for POST /v1/models/transfer.
type ModelTransferRequest struct {
	ModelID        string `json:"model_id"`
	TargetOrgID    string `json:"target_org_id,omitempty"`
	TargetUserID   string `json:"target_user_id,omitempty"`
}

// ModelShareRequest is the JSON body for POST /v1/models/share.
type ModelShareRequest struct {
	ModelID     string `json:"model_id"`
	GranteeType string `json:"grantee_type,omitempty"`
	GranteeID   string `json:"grantee_id,omitempty"`
	Permission  string `json:"permission,omitempty"`
}

// ModelShareGrant represents a share grant record on a model.
type ModelShareGrant struct {
	ID          string `json:"id"`
	ModelID     string `json:"model_id"`
	GranteeType string `json:"grantee_type,omitempty"`
	GranteeID   string `json:"grantee_id,omitempty"`
	Permission  string `json:"permission,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

// DatasetImagePayload represents a single image to upload to a model's dataset.
type DatasetImagePayload struct {
	URL     string `json:"url,omitempty"`
	Caption string `json:"caption,omitempty"`
	Base64  string `json:"base64,omitempty"`
}

// ModelDatasetImageResponse represents a dataset image after upload.
type ModelDatasetImageResponse struct {
	ID        string `json:"id"`
	ModelID   string `json:"model_id"`
	URL       string `json:"url,omitempty"`
	Caption   string `json:"caption,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

// TrainLoraJobRequest is the JSON body for POST /v1/models/{id}/jobs/train-lora.
type TrainLoraJobRequest struct {
	ConfigJSON map[string]interface{} `json:"config_json,omitempty"`
}

// GenerateCandidatesRequest is the JSON body for POST /v1/models/{id}/jobs/generate-candidates.
type GenerateCandidatesRequest struct {
	Prompt     string                 `json:"prompt,omitempty"`
	Count      int                    `json:"count,omitempty"`
	ConfigJSON map[string]interface{} `json:"config_json,omitempty"`
}

// ModelJobResponse represents a model job record returned by the API.
type ModelJobResponse struct {
	ID           string                 `json:"id"`
	ModelID      string                 `json:"model_id"`
	JobType      string                 `json:"job_type,omitempty"`
	Status       string                 `json:"status,omitempty"`
	ProgressJSON map[string]interface{} `json:"progress_json,omitempty"`
	ResultJSON   map[string]interface{} `json:"result_json,omitempty"`
	ErrorJSON    map[string]interface{} `json:"error_json,omitempty"`
	CreatedAt    string                 `json:"created_at,omitempty"`
	UpdatedAt    string                 `json:"updated_at,omitempty"`
}

// ModelCandidateResponse represents a candidate image for a model.
type ModelCandidateResponse struct {
	ID        string `json:"id"`
	ModelID   string `json:"model_id"`
	ImageURL  string `json:"image_url,omitempty"`
	Prompt    string `json:"prompt,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

// selectCanonicalRequest is the JSON body for POST /v1/models/{id}/canonical.
type selectCanonicalRequest struct {
	CandidateID string `json:"candidate_id"`
}

// GenerateAvatarRequest is the JSON body for POST /v1/models/{id}/jobs/generate-avatar.
type GenerateAvatarRequest struct {
	Prompt     string                 `json:"prompt,omitempty"`
	ConfigJSON map[string]interface{} `json:"config_json,omitempty"`
}

// AvatarResponse represents a generated avatar for a model.
type AvatarResponse struct {
	ID        string `json:"id"`
	ModelID   string `json:"model_id"`
	ImageURL  string `json:"image_url,omitempty"`
	Prompt    string `json:"prompt,omitempty"`
	Status    string `json:"status,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

// AvatarProfileResponse represents an avatar profile attached to a model.
type AvatarProfileResponse struct {
	ID          string                 `json:"id"`
	ModelID     string                 `json:"model_id"`
	Name        string                 `json:"name,omitempty"`
	ImageURL    string                 `json:"image_url,omitempty"`
	ConfigJSON  map[string]interface{} `json:"config_json,omitempty"`
	CreatedAt   string                 `json:"created_at,omitempty"`
	UpdatedAt   string                 `json:"updated_at,omitempty"`
}

// AvatarProfileCreateRequest is the JSON body for POST /v1/models/avatars.
type AvatarProfileCreateRequest struct {
	ModelID    string                 `json:"model_id"`
	Name       string                 `json:"name,omitempty"`
	ImageURL   string                 `json:"image_url,omitempty"`
	ConfigJSON map[string]interface{} `json:"config_json,omitempty"`
}

// CharacterCreationRequest is the JSON body for POST /v1/models/character-creations.
type CharacterCreationRequest struct {
	ModelID    string                 `json:"model_id,omitempty"`
	OwnerOrgID string                `json:"owner_org_id,omitempty"`
	InputJSON  map[string]interface{} `json:"input_json,omitempty"`
	ConfigJSON map[string]interface{} `json:"config_json,omitempty"`
}

// CharacterCreationEnqueueResponse is returned when a character creation is enqueued.
type CharacterCreationEnqueueResponse struct {
	ID     string `json:"id"`
	Status string `json:"status,omitempty"`
}

// CharacterCreationResponse represents a character creation record.
type CharacterCreationResponse struct {
	ID         string                 `json:"id"`
	ModelID    string                 `json:"model_id,omitempty"`
	Status     string                 `json:"status,omitempty"`
	InputJSON  map[string]interface{} `json:"input_json,omitempty"`
	ConfigJSON map[string]interface{} `json:"config_json,omitempty"`
	ResultJSON map[string]interface{} `json:"result_json,omitempty"`
	ErrorJSON  map[string]interface{} `json:"error_json,omitempty"`
	CreatedAt  string                 `json:"created_at,omitempty"`
	UpdatedAt  string                 `json:"updated_at,omitempty"`
}

// CharacterCreationListParams contains query parameters for listing character creations.
type CharacterCreationListParams struct {
	Page       int    `json:"page,omitempty"`
	PageSize   int    `json:"page_size,omitempty"`
	ModelID    string `json:"model_id,omitempty"`
	OwnerOrgID string `json:"owner_org_id,omitempty"`
	Status     string `json:"status,omitempty"`
}

// CharacterCreationUpdateRequest is the JSON body for PATCH /v1/models/character-creations/{id}.
type CharacterCreationUpdateRequest struct {
	Status     *string                `json:"status,omitempty"`
	InputJSON  map[string]interface{} `json:"input_json,omitempty"`
	ConfigJSON map[string]interface{} `json:"config_json,omitempty"`
	ResultJSON map[string]interface{} `json:"result_json,omitempty"`
}

// unexported request bodies ---------------------------------------------------

type revokeShareRequest struct {
	GrantID string `json:"grant_id"`
}

type datasetImagesRequest struct {
	Images []DatasetImagePayload `json:"images"`
}

// ---------------------------------------------------------------------------
// ModelsAPI
// ---------------------------------------------------------------------------

// ModelsAPI provides access to model management endpoints.
type ModelsAPI struct {
	client *InoueClient
}

// List returns a paginated list of models, optionally filtered by the given params.
func (a *ModelsAPI) List(ctx context.Context, params ModelListParams) (*Page[ModelResponse], error) {
	q := url.Values{}
	if params.Page > 0 {
		q.Set("page", fmt.Sprintf("%d", params.Page))
	}
	if params.PageSize > 0 {
		q.Set("page_size", fmt.Sprintf("%d", params.PageSize))
	}
	if params.OwnerOrgID != "" {
		q.Set("owner_org_id", params.OwnerOrgID)
	}
	if params.Status != "" {
		q.Set("status", params.Status)
	}
	if params.Scope != "" {
		q.Set("scope", params.Scope)
	}

	path := "/v1/models"
	if encoded := q.Encode(); encoded != "" {
		path += "?" + encoded
	}

	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", path, nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("list models: %w", err)
	}

	var result Page[ModelResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list models decode: %w", err)
	}
	return &result, nil
}

// Get returns a single model by ID.
func (a *ModelsAPI) Get(ctx context.Context, modelID string) (*ModelResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", "/v1/models/"+modelID, nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("get model: %w", err)
	}

	var result ModelResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get model decode: %w", err)
	}
	return &result, nil
}

// Create creates a new model.
func (a *ModelsAPI) Create(ctx context.Context, req ModelCreateRequest) (*ModelResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/models", req, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("create model: %w", err)
	}

	var result ModelResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create model decode: %w", err)
	}
	return &result, nil
}

// Update replaces model fields via POST /v1/models/update.
func (a *ModelsAPI) Update(ctx context.Context, req ModelUpdateRequest) (*ModelResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/models/update", req, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("update model: %w", err)
	}

	var result ModelResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("update model decode: %w", err)
	}
	return &result, nil
}

// Patch partially updates a model by ID.
func (a *ModelsAPI) Patch(ctx context.Context, modelID string, req ModelPatchRequest) (*ModelResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "PATCH", "/v1/models/"+modelID, req, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("patch model: %w", err)
	}

	var result ModelResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("patch model decode: %w", err)
	}
	return &result, nil
}

// Delete removes a model by ID.
func (a *ModelsAPI) Delete(ctx context.Context, modelID string) error {
	err := a.client.request(ctx, "DELETE", "/v1/models/"+modelID, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("delete model: %w", err)
	}
	return nil
}

// Fork creates a copy of an existing model under a new owner.
func (a *ModelsAPI) Fork(ctx context.Context, req ModelForkRequest) (*ModelResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/models/fork", req, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("fork model: %w", err)
	}

	var result ModelResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("fork model decode: %w", err)
	}
	return &result, nil
}

// Transfer moves a model to a different owner.
func (a *ModelsAPI) Transfer(ctx context.Context, req ModelTransferRequest) (*ModelResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/models/transfer", req, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("transfer model: %w", err)
	}

	var result ModelResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("transfer model decode: %w", err)
	}
	return &result, nil
}

// Share grants access to a model for another user or organisation.
func (a *ModelsAPI) Share(ctx context.Context, req ModelShareRequest) (*ModelShareGrant, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/models/share", req, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("share model: %w", err)
	}

	var result ModelShareGrant
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("share model decode: %w", err)
	}
	return &result, nil
}

// RevokeShare removes a share grant by its ID.
func (a *ModelsAPI) RevokeShare(ctx context.Context, grantID string) error {
	body := revokeShareRequest{GrantID: grantID}

	err := a.client.request(ctx, "POST", "/v1/models/revoke", body, nil, nil)
	if err != nil {
		return fmt.Errorf("revoke share: %w", err)
	}
	return nil
}

// ListShares returns a paginated list of share grants for a model.
func (a *ModelsAPI) ListShares(ctx context.Context, modelID string, page, pageSize int) (*Page[ModelShareGrant], error) {
	path := fmt.Sprintf("/v1/models/%s/shares?page=%d&page_size=%d", modelID, page, pageSize)

	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", path, nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("list shares: %w", err)
	}

	var result Page[ModelShareGrant]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list shares decode: %w", err)
	}
	return &result, nil
}

// UploadDatasetImages uploads one or more images to a model's training dataset.
func (a *ModelsAPI) UploadDatasetImages(ctx context.Context, modelID string, images []DatasetImagePayload) ([]ModelDatasetImageResponse, error) {
	body := datasetImagesRequest{Images: images}

	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/models/"+modelID+"/dataset/images", body, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("upload dataset images: %w", err)
	}

	var result []ModelDatasetImageResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("upload dataset images decode: %w", err)
	}
	return result, nil
}

// TrainLoraJob enqueues a LoRA training job for the given model.
func (a *ModelsAPI) TrainLoraJob(ctx context.Context, modelID string, req TrainLoraJobRequest) (*ModelJobResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/models/"+modelID+"/jobs/train-lora", req, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("train lora job: %w", err)
	}

	var result ModelJobResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("train lora job decode: %w", err)
	}
	return &result, nil
}

// GenerateCandidatesJob enqueues a candidate image generation job for a model.
func (a *ModelsAPI) GenerateCandidatesJob(ctx context.Context, modelID string, req GenerateCandidatesRequest) (*ModelJobResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/models/"+modelID+"/jobs/generate-candidates", req, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("generate candidates job: %w", err)
	}

	var result ModelJobResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("generate candidates job decode: %w", err)
	}
	return &result, nil
}

// ListModelJobs returns all jobs associated with a model.
func (a *ModelsAPI) ListModelJobs(ctx context.Context, modelID string) ([]ModelJobResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", "/v1/models/"+modelID+"/jobs", nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("list model jobs: %w", err)
	}

	var result []ModelJobResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list model jobs decode: %w", err)
	}
	return result, nil
}

// GetModelJob returns a single model job by its ID.
func (a *ModelsAPI) GetModelJob(ctx context.Context, jobID string) (*ModelJobResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", "/v1/models/model-jobs/"+jobID, nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("get model job: %w", err)
	}

	var result ModelJobResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get model job decode: %w", err)
	}
	return &result, nil
}

// ListCandidates returns all candidate images for a model.
func (a *ModelsAPI) ListCandidates(ctx context.Context, modelID string) ([]ModelCandidateResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", "/v1/models/"+modelID+"/candidates", nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("list candidates: %w", err)
	}

	var result []ModelCandidateResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list candidates decode: %w", err)
	}
	return result, nil
}

// SelectCanonical sets the canonical (primary) image for a model from a candidate.
func (a *ModelsAPI) SelectCanonical(ctx context.Context, modelID, candidateID string) (*ModelResponse, error) {
	body := selectCanonicalRequest{CandidateID: candidateID}

	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/models/"+modelID+"/canonical", body, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("select canonical: %w", err)
	}

	var result ModelResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("select canonical decode: %w", err)
	}
	return &result, nil
}

// GenerateAvatarJob enqueues an avatar generation job for a model.
func (a *ModelsAPI) GenerateAvatarJob(ctx context.Context, modelID string, req GenerateAvatarRequest) (*ModelJobResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/models/"+modelID+"/jobs/generate-avatar", req, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("generate avatar job: %w", err)
	}

	var result ModelJobResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("generate avatar job decode: %w", err)
	}
	return &result, nil
}

// ListGeneratedAvatars returns all generated avatars for a model.
func (a *ModelsAPI) ListGeneratedAvatars(ctx context.Context, modelID string) ([]AvatarResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", "/v1/models/"+modelID+"/generated-avatars", nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("list generated avatars: %w", err)
	}

	var result []AvatarResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list generated avatars decode: %w", err)
	}
	return result, nil
}

// GetGeneratedAvatar returns a single generated avatar by its ID.
func (a *ModelsAPI) GetGeneratedAvatar(ctx context.Context, avatarID string) (*AvatarResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", "/v1/models/avatars/"+avatarID, nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("get generated avatar: %w", err)
	}

	var result AvatarResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get generated avatar decode: %w", err)
	}
	return &result, nil
}

// ListAvatars returns a paginated list of avatar profiles for a model.
func (a *ModelsAPI) ListAvatars(ctx context.Context, modelID string, page, pageSize int) (*Page[AvatarProfileResponse], error) {
	path := fmt.Sprintf("/v1/models/%s/avatars?page=%d&page_size=%d", modelID, page, pageSize)

	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", path, nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("list avatars: %w", err)
	}

	var result Page[AvatarProfileResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list avatars decode: %w", err)
	}
	return &result, nil
}

// GetAvatar returns a single avatar profile by its ID.
func (a *ModelsAPI) GetAvatar(ctx context.Context, avatarID string) (*AvatarProfileResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", "/v1/models/avatar/"+avatarID, nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("get avatar: %w", err)
	}

	var result AvatarProfileResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get avatar decode: %w", err)
	}
	return &result, nil
}

// CreateAvatar creates a new avatar profile.
func (a *ModelsAPI) CreateAvatar(ctx context.Context, req AvatarProfileCreateRequest) (*AvatarProfileResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/models/avatars", req, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("create avatar: %w", err)
	}

	var result AvatarProfileResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create avatar decode: %w", err)
	}
	return &result, nil
}

// DeleteAvatar removes an avatar profile by its ID.
func (a *ModelsAPI) DeleteAvatar(ctx context.Context, avatarID string) error {
	err := a.client.request(ctx, "DELETE", "/v1/models/avatar/"+avatarID, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("delete avatar: %w", err)
	}
	return nil
}

// CharacterSchema returns the JSON schema for the character designer.
func (a *ModelsAPI) CharacterSchema(ctx context.Context) (map[string]interface{}, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", "/v1/models/character-designer/schema", nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("character schema: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("character schema decode: %w", err)
	}
	return result, nil
}

// CreateCharacterCreation enqueues a new character creation.
func (a *ModelsAPI) CreateCharacterCreation(ctx context.Context, req CharacterCreationRequest) (*CharacterCreationEnqueueResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/models/character-creations", req, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("create character creation: %w", err)
	}

	var result CharacterCreationEnqueueResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create character creation decode: %w", err)
	}
	return &result, nil
}

// ListCharacterCreations returns a paginated list of character creations.
func (a *ModelsAPI) ListCharacterCreations(ctx context.Context, params CharacterCreationListParams) (*Page[CharacterCreationResponse], error) {
	q := url.Values{}
	if params.Page > 0 {
		q.Set("page", fmt.Sprintf("%d", params.Page))
	}
	if params.PageSize > 0 {
		q.Set("page_size", fmt.Sprintf("%d", params.PageSize))
	}
	if params.ModelID != "" {
		q.Set("model_id", params.ModelID)
	}
	if params.OwnerOrgID != "" {
		q.Set("owner_org_id", params.OwnerOrgID)
	}
	if params.Status != "" {
		q.Set("status", params.Status)
	}

	path := "/v1/models/character-creations"
	if encoded := q.Encode(); encoded != "" {
		path += "?" + encoded
	}

	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", path, nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("list character creations: %w", err)
	}

	var result Page[CharacterCreationResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list character creations decode: %w", err)
	}
	return &result, nil
}

// GetCharacterCreation returns a single character creation by ID.
func (a *ModelsAPI) GetCharacterCreation(ctx context.Context, creationID string) (*CharacterCreationResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", "/v1/models/character-creations/"+creationID, nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("get character creation: %w", err)
	}

	var result CharacterCreationResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get character creation decode: %w", err)
	}
	return &result, nil
}

// PatchCharacterCreation partially updates a character creation by ID.
func (a *ModelsAPI) PatchCharacterCreation(ctx context.Context, creationID string, req CharacterCreationUpdateRequest) (*CharacterCreationResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "PATCH", "/v1/models/character-creations/"+creationID, req, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("patch character creation: %w", err)
	}

	var result CharacterCreationResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("patch character creation decode: %w", err)
	}
	return &result, nil
}
