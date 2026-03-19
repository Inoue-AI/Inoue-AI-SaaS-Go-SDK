package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// AuditAPI provides access to the audit log endpoints.
type AuditAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// AuditLogResponse represents audit log data.
type AuditLogResponse struct {
	Entries []AuditLogEntry `json:"entries"`
}

// AuditLogEntry represents a single audit log entry.
type AuditLogEntry struct {
	ID         string                 `json:"id"`
	UserID     string                 `json:"user_id"`
	OrgID      string                 `json:"org_id"`
	Action     string                 `json:"action"`
	Resource   string                 `json:"resource"`
	ResourceID string                 `json:"resource_id"`
	Details    map[string]interface{} `json:"details"`
	IPAddress  string                 `json:"ip_address"`
	UserAgent  string                 `json:"user_agent"`
	CreatedAt  string                 `json:"created_at"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// Me returns the audit log for the current user.
func (a *AuditAPI) Me(ctx context.Context) (*AuditLogResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/audit/me", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("audit me: %w", err)
	}
	var result AuditLogResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("audit me decode: %w", err)
	}
	return &result, nil
}

// Org returns the audit log for an organization.
func (a *AuditAPI) Org(ctx context.Context, orgID string) (*AuditLogResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/audit/orgs/"+orgID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("audit org: %w", err)
	}
	var result AuditLogResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("audit org decode: %w", err)
	}
	return &result, nil
}
