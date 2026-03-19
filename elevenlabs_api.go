package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// ElevenLabsAPI provides access to the ElevenLabs integration endpoints.
type ElevenLabsAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// ElevenLabsKeyCreateRequest is the body for creating an ElevenLabs key.
type ElevenLabsKeyCreateRequest struct {
	Name       string `json:"name"`
	APIKey     string `json:"api_key"`
	OwnerOrgID string `json:"owner_org_id,omitempty"`
}

// ElevenLabsKeyUpdateRequest is the body for updating an ElevenLabs key.
type ElevenLabsKeyUpdateRequest struct {
	Name   string `json:"name,omitempty"`
	APIKey string `json:"api_key,omitempty"`
}

// ElevenLabsShareRequest is the body for sharing an ElevenLabs key.
type ElevenLabsShareRequest struct {
	UserID     string `json:"user_id,omitempty"`
	OrgID      string `json:"org_id,omitempty"`
	Permission string `json:"permission,omitempty"`
}

// ElevenLabsVoiceModelLinksRequest is the body for updating voice-model links.
type ElevenLabsVoiceModelLinksRequest struct {
	Links []map[string]interface{} `json:"links"`
}

// ElevenLabsCloneVoiceRequest is the body for cloning a voice.
type ElevenLabsCloneVoiceRequest struct {
	KeyID       string   `json:"key_id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	AssetIDs    []string `json:"asset_ids"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// ElevenLabsDesignVoiceRequest is the body for designing a voice.
type ElevenLabsDesignVoiceRequest struct {
	KeyID       string                 `json:"key_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// ElevenLabsVoiceUpdateRequest is the body for updating a voice.
type ElevenLabsVoiceUpdateRequest struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// ElevenLabsVoiceSettingsRequest is the body for updating voice settings.
type ElevenLabsVoiceSettingsRequest struct {
	Stability       float64 `json:"stability"`
	SimilarityBoost float64 `json:"similarity_boost"`
	Style           float64 `json:"style,omitempty"`
	UseSpeakerBoost bool    `json:"use_speaker_boost,omitempty"`
}

// ElevenLabsTTSRequest is the body for text-to-speech.
type ElevenLabsTTSRequest struct {
	KeyID    string                 `json:"key_id"`
	VoiceID  string                 `json:"voice_id"`
	Text     string                 `json:"text"`
	ModelID  string                 `json:"model_id,omitempty"`
	Settings map[string]interface{} `json:"settings,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// ElevenLabsKeyResponse represents an ElevenLabs API key.
type ElevenLabsKeyResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	OwnerUserID string `json:"owner_user_id"`
	OwnerOrgID  string `json:"owner_org_id"`
	Masked      string `json:"masked"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ElevenLabsModelResponse represents an ElevenLabs model.
type ElevenLabsModelResponse struct {
	ModelID     string `json:"model_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ElevenLabsVoiceResponse represents an ElevenLabs voice.
type ElevenLabsVoiceResponse struct {
	VoiceID     string                 `json:"voice_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	Labels      map[string]string      `json:"labels"`
	PreviewURL  string                 `json:"preview_url"`
	Settings    map[string]interface{} `json:"settings"`
}

// ElevenLabsVoiceSettingsResponse represents voice settings.
type ElevenLabsVoiceSettingsResponse struct {
	Stability       float64 `json:"stability"`
	SimilarityBoost float64 `json:"similarity_boost"`
	Style           float64 `json:"style"`
	UseSpeakerBoost bool    `json:"use_speaker_boost"`
}

// ElevenLabsTTSResponse represents the result of a text-to-speech request.
type ElevenLabsTTSResponse struct {
	AudioURL  string `json:"audio_url"`
	AssetID   string `json:"asset_id"`
	Duration  float64 `json:"duration"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// ListKeys returns the list of ElevenLabs API keys.
func (a *ElevenLabsAPI) ListKeys(ctx context.Context) ([]ElevenLabsKeyResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/elevenlabs/keys", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list elevenlabs keys: %w", err)
	}
	var result []ElevenLabsKeyResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list elevenlabs keys decode: %w", err)
	}
	return result, nil
}

// CreateKey creates a new ElevenLabs API key.
func (a *ElevenLabsAPI) CreateKey(ctx context.Context, req ElevenLabsKeyCreateRequest) (*ElevenLabsKeyResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/elevenlabs/keys", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("create elevenlabs key: %w", err)
	}
	var result ElevenLabsKeyResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create elevenlabs key decode: %w", err)
	}
	return &result, nil
}

// UpdateKey updates an existing ElevenLabs API key.
func (a *ElevenLabsAPI) UpdateKey(ctx context.Context, keyID string, req ElevenLabsKeyUpdateRequest) (*ElevenLabsKeyResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "PATCH", "/v1/elevenlabs/keys/"+keyID, req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("update elevenlabs key: %w", err)
	}
	var result ElevenLabsKeyResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("update elevenlabs key decode: %w", err)
	}
	return &result, nil
}

// DeleteKey deletes an ElevenLabs API key.
func (a *ElevenLabsAPI) DeleteKey(ctx context.Context, keyID string) error {
	body := map[string]string{"key_id": keyID}
	if err := a.client.request(ctx, "POST", "/v1/elevenlabs/keys/delete", body, nil, nil); err != nil {
		return fmt.Errorf("delete elevenlabs key: %w", err)
	}
	return nil
}

// ShareKey shares an ElevenLabs key with a user or organization.
func (a *ElevenLabsAPI) ShareKey(ctx context.Context, keyID string, req ElevenLabsShareRequest) (*ModelShareGrant, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/elevenlabs/keys/"+keyID+"/share", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("share elevenlabs key: %w", err)
	}
	var result ModelShareGrant
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("share elevenlabs key decode: %w", err)
	}
	return &result, nil
}

// RevokeKey revokes a share on an ElevenLabs key.
func (a *ElevenLabsAPI) RevokeKey(ctx context.Context, keyID, grantID string) error {
	path := fmt.Sprintf("/v1/elevenlabs/keys/%s/share/%s", keyID, grantID)
	if err := a.client.request(ctx, "DELETE", path, nil, nil, nil); err != nil {
		return fmt.Errorf("revoke elevenlabs key share: %w", err)
	}
	return nil
}

// ListShares returns the list of shares for an ElevenLabs key.
func (a *ElevenLabsAPI) ListShares(ctx context.Context, keyID string) ([]ModelShareGrant, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/elevenlabs/keys/"+keyID+"/shares", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list elevenlabs key shares: %w", err)
	}
	var result []ModelShareGrant
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list elevenlabs key shares decode: %w", err)
	}
	return result, nil
}

// ListModels returns the list of available ElevenLabs models.
func (a *ElevenLabsAPI) ListModels(ctx context.Context, keyID string) ([]ElevenLabsModelResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/elevenlabs/models?key_id="+keyID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list elevenlabs models: %w", err)
	}
	var result []ElevenLabsModelResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list elevenlabs models decode: %w", err)
	}
	return result, nil
}

// ListVoices returns the list of available ElevenLabs voices.
func (a *ElevenLabsAPI) ListVoices(ctx context.Context, keyID string) ([]ElevenLabsVoiceResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/elevenlabs/voices?key_id="+keyID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list elevenlabs voices: %w", err)
	}
	var result []ElevenLabsVoiceResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list elevenlabs voices decode: %w", err)
	}
	return result, nil
}

// GetVoice retrieves a specific ElevenLabs voice.
func (a *ElevenLabsAPI) GetVoice(ctx context.Context, keyID, voiceID string) (*ElevenLabsVoiceResponse, error) {
	path := fmt.Sprintf("/v1/elevenlabs/voices/%s?key_id=%s", voiceID, keyID)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get elevenlabs voice: %w", err)
	}
	var result ElevenLabsVoiceResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get elevenlabs voice decode: %w", err)
	}
	return &result, nil
}

// VoiceModelLinks returns the voice-model links for a key.
func (a *ElevenLabsAPI) VoiceModelLinks(ctx context.Context, keyID string) ([]map[string]interface{}, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/elevenlabs/voice-model-links?key_id="+keyID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("voice model links: %w", err)
	}
	var result []map[string]interface{}
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("voice model links decode: %w", err)
	}
	return result, nil
}

// UpdateVoiceModelLinks updates the voice-model links.
func (a *ElevenLabsAPI) UpdateVoiceModelLinks(ctx context.Context, req ElevenLabsVoiceModelLinksRequest) ([]map[string]interface{}, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "PUT", "/v1/elevenlabs/voice-model-links", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("update voice model links: %w", err)
	}
	var result []map[string]interface{}
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("update voice model links decode: %w", err)
	}
	return result, nil
}

// CloneVoice clones a voice using uploaded audio samples.
func (a *ElevenLabsAPI) CloneVoice(ctx context.Context, req ElevenLabsCloneVoiceRequest) (*ElevenLabsVoiceResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/elevenlabs/voices/clone", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("clone voice: %w", err)
	}
	var result ElevenLabsVoiceResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("clone voice decode: %w", err)
	}
	return &result, nil
}

// DesignVoice designs a new voice using AI parameters.
func (a *ElevenLabsAPI) DesignVoice(ctx context.Context, req ElevenLabsDesignVoiceRequest) (*ElevenLabsVoiceResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/elevenlabs/voices/design", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("design voice: %w", err)
	}
	var result ElevenLabsVoiceResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("design voice decode: %w", err)
	}
	return &result, nil
}

// UpdateVoice updates a voice's properties.
func (a *ElevenLabsAPI) UpdateVoice(ctx context.Context, keyID, voiceID string, req ElevenLabsVoiceUpdateRequest) (*ElevenLabsVoiceResponse, error) {
	path := fmt.Sprintf("/v1/elevenlabs/voices/%s?key_id=%s", voiceID, keyID)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "PATCH", path, req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("update voice: %w", err)
	}
	var result ElevenLabsVoiceResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("update voice decode: %w", err)
	}
	return &result, nil
}

// DeleteVoice deletes a voice.
func (a *ElevenLabsAPI) DeleteVoice(ctx context.Context, keyID, voiceID string) error {
	body := map[string]string{"key_id": keyID, "voice_id": voiceID}
	if err := a.client.request(ctx, "POST", "/v1/elevenlabs/voices/delete", body, nil, nil); err != nil {
		return fmt.Errorf("delete voice: %w", err)
	}
	return nil
}

// VoiceSettings returns the settings for a voice.
func (a *ElevenLabsAPI) VoiceSettings(ctx context.Context, keyID, voiceID string) (*ElevenLabsVoiceSettingsResponse, error) {
	path := fmt.Sprintf("/v1/elevenlabs/voices/%s/settings?key_id=%s", voiceID, keyID)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("voice settings: %w", err)
	}
	var result ElevenLabsVoiceSettingsResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("voice settings decode: %w", err)
	}
	return &result, nil
}

// UpdateVoiceSettings updates the settings for a voice.
func (a *ElevenLabsAPI) UpdateVoiceSettings(ctx context.Context, keyID, voiceID string, req ElevenLabsVoiceSettingsRequest) (*ElevenLabsVoiceSettingsResponse, error) {
	path := fmt.Sprintf("/v1/elevenlabs/voices/%s/settings?key_id=%s", voiceID, keyID)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "PUT", path, req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("update voice settings: %w", err)
	}
	var result ElevenLabsVoiceSettingsResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("update voice settings decode: %w", err)
	}
	return &result, nil
}

// TextToSpeech converts text to speech using a voice.
func (a *ElevenLabsAPI) TextToSpeech(ctx context.Context, req ElevenLabsTTSRequest) (*ElevenLabsTTSResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/elevenlabs/tts", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("text to speech: %w", err)
	}
	var result ElevenLabsTTSResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("text to speech decode: %w", err)
	}
	return &result, nil
}
