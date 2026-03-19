package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// FanvueAPI provides access to the Fanvue integration endpoints.
type FanvueAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// FanvueConnectStartRequest is the body for initiating a Fanvue connection.
type FanvueConnectStartRequest struct {
	ModelID     string `json:"model_id"`
	RedirectURI string `json:"redirect_uri,omitempty"`
}

// FanvueCallbackRequest is the body for the OAuth callback.
type FanvueCallbackRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

// FanvueRefreshTokenRequest is the body for refreshing a Fanvue token.
type FanvueRefreshTokenRequest struct {
	ModelID string `json:"model_id"`
}

// FanvueWebhookRequest is the body for registering a Fanvue webhook.
type FanvueWebhookRequest struct {
	ModelID  string   `json:"model_id"`
	URL      string   `json:"url"`
	Events   []string `json:"events,omitempty"`
}

// FanvueMessageRequest is the body for sending a Fanvue message.
type FanvueMessageRequest struct {
	ModelID        string   `json:"model_id"`
	ConversationID string   `json:"conversation_id"`
	Text           string   `json:"text,omitempty"`
	AssetIDs       []string `json:"asset_ids,omitempty"`
	Price          float64  `json:"price,omitempty"`
}

// FanvueLockRequest is the body for managing Fanvue locks.
type FanvueLockRequest struct {
	ModelID    string  `json:"model_id"`
	ContentID  string  `json:"content_id"`
	Price      float64 `json:"price,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// FanvueConnectResponse represents the Fanvue connection initiation response.
type FanvueConnectResponse struct {
	AuthURL string `json:"auth_url"`
	State   string `json:"state"`
}

// FanvueAccountResponse represents a connected Fanvue account.
type FanvueAccountResponse struct {
	ID          string                 `json:"id"`
	ModelID     string                 `json:"model_id"`
	Username    string                 `json:"username"`
	DisplayName string                 `json:"display_name"`
	AvatarURL   string                 `json:"avatar_url"`
	Connected   bool                   `json:"connected"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

// FanvueTokenResponse represents a token refresh result.
type FanvueTokenResponse struct {
	Success   bool   `json:"success"`
	ExpiresAt string `json:"expires_at"`
}

// FanvueWebhookResponse represents a Fanvue webhook.
type FanvueWebhookResponse struct {
	ID       string   `json:"id"`
	ModelID  string   `json:"model_id"`
	URL      string   `json:"url"`
	Events   []string `json:"events"`
	Active   bool     `json:"active"`
	CreatedAt string  `json:"created_at"`
}

// FanvueEventResponse represents a Fanvue event.
type FanvueEventResponse struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt string                 `json:"created_at"`
}

// FanvueConversationResponse represents a Fanvue conversation.
type FanvueConversationResponse struct {
	ID             string                 `json:"id"`
	ModelID        string                 `json:"model_id"`
	ParticipantID  string                 `json:"participant_id"`
	Username       string                 `json:"username"`
	LastMessage    string                 `json:"last_message"`
	UnreadCount    int                    `json:"unread_count"`
	Metadata       map[string]interface{} `json:"metadata"`
	CreatedAt      string                 `json:"created_at"`
	UpdatedAt      string                 `json:"updated_at"`
}

// FanvueMessageResponse represents a Fanvue message.
type FanvueMessageResponse struct {
	ID             string                 `json:"id"`
	ConversationID string                 `json:"conversation_id"`
	SenderID       string                 `json:"sender_id"`
	Text           string                 `json:"text"`
	Attachments    []map[string]interface{} `json:"attachments"`
	Price          float64                `json:"price"`
	CreatedAt      string                 `json:"created_at"`
}

// FanvueLockResponse represents a Fanvue content lock.
type FanvueLockResponse struct {
	ID        string  `json:"id"`
	ContentID string  `json:"content_id"`
	ModelID   string  `json:"model_id"`
	Price     float64 `json:"price"`
	Locked    bool    `json:"locked"`
	CreatedAt string  `json:"created_at"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// ConnectStart initiates a Fanvue OAuth connection.
func (a *FanvueAPI) ConnectStart(ctx context.Context, req FanvueConnectStartRequest) (*FanvueConnectResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/fanvue/connect/start", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("fanvue connect start: %w", err)
	}
	var result FanvueConnectResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("fanvue connect start decode: %w", err)
	}
	return &result, nil
}

// Callback handles the Fanvue OAuth callback.
func (a *FanvueAPI) Callback(ctx context.Context, req FanvueCallbackRequest) (*FanvueAccountResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/fanvue/connect/callback", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("fanvue callback: %w", err)
	}
	var result FanvueAccountResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("fanvue callback decode: %w", err)
	}
	return &result, nil
}

// RefreshToken refreshes the Fanvue access token for a model.
func (a *FanvueAPI) RefreshToken(ctx context.Context, req FanvueRefreshTokenRequest) (*FanvueTokenResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/fanvue/refresh", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("fanvue refresh token: %w", err)
	}
	var result FanvueTokenResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("fanvue refresh token decode: %w", err)
	}
	return &result, nil
}

// RegisterWebhook registers a Fanvue webhook.
func (a *FanvueAPI) RegisterWebhook(ctx context.Context, req FanvueWebhookRequest) (*FanvueWebhookResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/fanvue/webhooks", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("register fanvue webhook: %w", err)
	}
	var result FanvueWebhookResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("register fanvue webhook decode: %w", err)
	}
	return &result, nil
}

// ListEvents returns Fanvue events for a model.
func (a *FanvueAPI) ListEvents(ctx context.Context, modelID string, page, pageSize int) (*Page[FanvueEventResponse], error) {
	path := fmt.Sprintf("/v1/fanvue/events/%s?page=%d&page_size=%d", modelID, page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list fanvue events: %w", err)
	}
	var result Page[FanvueEventResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list fanvue events decode: %w", err)
	}
	return &result, nil
}

// Accounts returns the list of connected Fanvue accounts.
func (a *FanvueAPI) Accounts(ctx context.Context) ([]FanvueAccountResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/fanvue/accounts", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list fanvue accounts: %w", err)
	}
	var result []FanvueAccountResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list fanvue accounts decode: %w", err)
	}
	return result, nil
}

// Account retrieves a specific Fanvue account by model ID.
func (a *FanvueAPI) Account(ctx context.Context, modelID string) (*FanvueAccountResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/fanvue/accounts/"+modelID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get fanvue account: %w", err)
	}
	var result FanvueAccountResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get fanvue account decode: %w", err)
	}
	return &result, nil
}

// Disconnect disconnects a Fanvue account.
func (a *FanvueAPI) Disconnect(ctx context.Context, modelID string) error {
	body := map[string]string{"model_id": modelID}
	if err := a.client.request(ctx, "POST", "/v1/fanvue/disconnect", body, nil, nil); err != nil {
		return fmt.Errorf("fanvue disconnect: %w", err)
	}
	return nil
}

// ListConversations returns Fanvue conversations for a model.
func (a *FanvueAPI) ListConversations(ctx context.Context, modelID string, page, pageSize int) (*Page[FanvueConversationResponse], error) {
	path := fmt.Sprintf("/v1/fanvue/conversations/%s?page=%d&page_size=%d", modelID, page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list fanvue conversations: %w", err)
	}
	var result Page[FanvueConversationResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list fanvue conversations decode: %w", err)
	}
	return &result, nil
}

// ListMessages returns messages in a Fanvue conversation.
func (a *FanvueAPI) ListMessages(ctx context.Context, modelID, conversationID string, page, pageSize int) (*Page[FanvueMessageResponse], error) {
	path := fmt.Sprintf("/v1/fanvue/conversations/%s/%s/messages?page=%d&page_size=%d", modelID, conversationID, page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list fanvue messages: %w", err)
	}
	var result Page[FanvueMessageResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list fanvue messages decode: %w", err)
	}
	return &result, nil
}

// SendMessage sends a message in a Fanvue conversation.
func (a *FanvueAPI) SendMessage(ctx context.Context, req FanvueMessageRequest) (*FanvueMessageResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/fanvue/messages", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("send fanvue message: %w", err)
	}
	var result FanvueMessageResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("send fanvue message decode: %w", err)
	}
	return &result, nil
}

// LockContent creates a lock on Fanvue content.
func (a *FanvueAPI) LockContent(ctx context.Context, req FanvueLockRequest) (*FanvueLockResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/fanvue/locks", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("fanvue lock content: %w", err)
	}
	var result FanvueLockResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("fanvue lock content decode: %w", err)
	}
	return &result, nil
}

// UnlockContent removes a lock from Fanvue content.
func (a *FanvueAPI) UnlockContent(ctx context.Context, modelID, contentID string) error {
	body := map[string]string{"model_id": modelID, "content_id": contentID}
	if err := a.client.request(ctx, "POST", "/v1/fanvue/locks/remove", body, nil, nil); err != nil {
		return fmt.Errorf("fanvue unlock content: %w", err)
	}
	return nil
}

// ListLocks returns Fanvue locks for a model.
func (a *FanvueAPI) ListLocks(ctx context.Context, modelID string) ([]FanvueLockResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/fanvue/locks/"+modelID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list fanvue locks: %w", err)
	}
	var result []FanvueLockResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list fanvue locks decode: %w", err)
	}
	return result, nil
}
