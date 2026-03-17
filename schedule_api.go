package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// ScheduleAPI provides access to the content scheduling endpoints.
type ScheduleAPI struct {
	client *InoueClient
}

// ScheduleEntryCreateParams contains the parameters for creating a new schedule entry
// via the internal worker endpoint. OwnerUserID and OwnerOrgID identify who the entry
// belongs to — the worker acts on behalf of the dispatcher which already validated ownership.
type ScheduleEntryCreateParams struct {
	OwnerUserID   *string  `json:"owner_user_id,omitempty"`
	OwnerOrgID    *string  `json:"owner_org_id,omitempty"`
	ModelID       string   `json:"model_id,omitempty"`
	PlatformID    string   `json:"platform_id,omitempty"`
	ContentTypeID string   `json:"content_type_id,omitempty"`
	ScheduledFor  string   `json:"scheduled_for"`
	AssetIDs      []string `json:"asset_ids,omitempty"`
	Notes         *string  `json:"notes,omitempty"`
}

// ScheduleEntry represents a scheduled content entry returned by the API.
type ScheduleEntry struct {
	ID            string   `json:"id"`
	Status        string   `json:"status"`
	ModelID       string   `json:"model_id"`
	PlatformID    string   `json:"platform_id"`
	ContentTypeID string   `json:"content_type_id"`
	ScheduledFor  string   `json:"scheduled_for"`
	AssetIDs      []string `json:"asset_ids"`
	Notes         *string  `json:"notes"`
	CreatedAt     string   `json:"created_at"`
}

// Create creates a new schedule entry via the internal worker endpoint.
// This endpoint trusts the OwnerUserID/OwnerOrgID from the request body
// and does not require org membership checks — the worker JWT is sufficient.
func (a *ScheduleAPI) Create(ctx context.Context, params ScheduleEntryCreateParams) (*ScheduleEntry, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/internal/schedule/create", params, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("create schedule entry: %w", err)
	}

	var entry ScheduleEntry
	if err := json.Unmarshal(apiResp.Data, &entry); err != nil {
		return nil, fmt.Errorf("create schedule entry: failed to decode response: %w", err)
	}

	return &entry, nil
}
