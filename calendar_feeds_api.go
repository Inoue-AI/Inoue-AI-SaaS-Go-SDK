package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// CalendarFeedsAPI provides access to the calendar feed endpoints.
type CalendarFeedsAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// CalendarFeedCreateRequest is the body for creating a calendar feed.
type CalendarFeedCreateRequest struct {
	Name       string `json:"name"`
	ModelID    string `json:"model_id,omitempty"`
	OwnerOrgID string `json:"owner_org_id,omitempty"`
}

// CalendarFeedUpdateRequest is the body for updating a calendar feed.
type CalendarFeedUpdateRequest struct {
	Name string `json:"name,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// CalendarFeedResponse represents a calendar feed record.
type CalendarFeedResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ModelID    string `json:"model_id"`
	OwnerUserID string `json:"owner_user_id"`
	OwnerOrgID string `json:"owner_org_id"`
	Token      string `json:"token"`
	FeedURL    string `json:"feed_url"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// List returns the list of calendar feeds.
func (a *CalendarFeedsAPI) List(ctx context.Context) ([]CalendarFeedResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/calendar-feeds", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list calendar feeds: %w", err)
	}
	var result []CalendarFeedResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list calendar feeds decode: %w", err)
	}
	return result, nil
}

// Create creates a new calendar feed.
func (a *CalendarFeedsAPI) Create(ctx context.Context, req CalendarFeedCreateRequest) (*CalendarFeedResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/calendar-feeds", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("create calendar feed: %w", err)
	}
	var result CalendarFeedResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create calendar feed decode: %w", err)
	}
	return &result, nil
}

// Update updates an existing calendar feed.
func (a *CalendarFeedsAPI) Update(ctx context.Context, feedID string, req CalendarFeedUpdateRequest) (*CalendarFeedResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "PATCH", "/v1/calendar-feeds/"+feedID, req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("update calendar feed: %w", err)
	}
	var result CalendarFeedResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("update calendar feed decode: %w", err)
	}
	return &result, nil
}

// Delete deletes a calendar feed by ID.
func (a *CalendarFeedsAPI) Delete(ctx context.Context, feedID string) error {
	body := map[string]string{"feed_id": feedID}
	if err := a.client.request(ctx, "POST", "/v1/calendar-feeds/delete", body, nil, nil); err != nil {
		return fmt.Errorf("delete calendar feed: %w", err)
	}
	return nil
}

// RegenerateToken regenerates the token for a calendar feed.
func (a *CalendarFeedsAPI) RegenerateToken(ctx context.Context, feedID string) (*CalendarFeedResponse, error) {
	body := map[string]string{"feed_id": feedID}
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/calendar-feeds/regenerate-token", body, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("regenerate calendar feed token: %w", err)
	}
	var result CalendarFeedResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("regenerate calendar feed token decode: %w", err)
	}
	return &result, nil
}
