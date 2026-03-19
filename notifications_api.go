package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// NotificationsAPI provides access to the notification endpoints.
type NotificationsAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// NotificationCreateRequest is the body for creating a notification.
type NotificationCreateRequest struct {
	Title      string                 `json:"title"`
	Body       string                 `json:"body,omitempty"`
	Type       string                 `json:"type,omitempty"`
	ActionURL  string                 `json:"action_url,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// NotificationResponse represents a notification record.
type NotificationResponse struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	Title     string                 `json:"title"`
	Body      string                 `json:"body"`
	Type      string                 `json:"type"`
	ActionURL string                 `json:"action_url"`
	Read      bool                   `json:"read"`
	Muted     bool                   `json:"muted"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// List returns a paginated list of notifications.
func (a *NotificationsAPI) List(ctx context.Context, page, pageSize int) (*Page[NotificationResponse], error) {
	path := fmt.Sprintf("/v1/notifications?page=%d&page_size=%d", page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}
	var result Page[NotificationResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list notifications decode: %w", err)
	}
	return &result, nil
}

// Create creates a new notification.
func (a *NotificationsAPI) Create(ctx context.Context, req NotificationCreateRequest) (*NotificationResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/notifications", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("create notification: %w", err)
	}
	var result NotificationResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create notification decode: %w", err)
	}
	return &result, nil
}

// Delete deletes a notification by ID.
func (a *NotificationsAPI) Delete(ctx context.Context, notificationID string) error {
	body := map[string]string{"notification_id": notificationID}
	if err := a.client.request(ctx, "POST", "/v1/notifications/delete", body, nil, nil); err != nil {
		return fmt.Errorf("delete notification: %w", err)
	}
	return nil
}

// MarkRead marks a notification as read.
func (a *NotificationsAPI) MarkRead(ctx context.Context, notificationID string) error {
	body := map[string]string{"notification_id": notificationID}
	if err := a.client.request(ctx, "POST", "/v1/notifications/read", body, nil, nil); err != nil {
		return fmt.Errorf("mark notification read: %w", err)
	}
	return nil
}

// MarkUnread marks a notification as unread.
func (a *NotificationsAPI) MarkUnread(ctx context.Context, notificationID string) error {
	body := map[string]string{"notification_id": notificationID}
	if err := a.client.request(ctx, "POST", "/v1/notifications/unread", body, nil, nil); err != nil {
		return fmt.Errorf("mark notification unread: %w", err)
	}
	return nil
}

// Mute mutes a notification.
func (a *NotificationsAPI) Mute(ctx context.Context, notificationID string) error {
	body := map[string]string{"notification_id": notificationID}
	if err := a.client.request(ctx, "POST", "/v1/notifications/mute", body, nil, nil); err != nil {
		return fmt.Errorf("mute notification: %w", err)
	}
	return nil
}

// Unmute unmutes a notification.
func (a *NotificationsAPI) Unmute(ctx context.Context, notificationID string) error {
	body := map[string]string{"notification_id": notificationID}
	if err := a.client.request(ctx, "POST", "/v1/notifications/unmute", body, nil, nil); err != nil {
		return fmt.Errorf("unmute notification: %w", err)
	}
	return nil
}
