package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// ---------------------------------------------------------------------------
// Orgs types
// ---------------------------------------------------------------------------

// OrgResponse represents an organisation record returned by the API.
type OrgResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	OwnerUserID string `json:"owner_user_id,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// OrgOverviewResponse is the enriched organisation view that includes members
// and recent audit log entries.
type OrgOverviewResponse struct {
	Org       OrgResponse          `json:"org"`
	Members   []OrgOverviewMember  `json:"members"`
	AuditLogs []OrgOverviewAudit   `json:"audit_logs"`
}

// OrgOverviewMember is a member entry within an OrgOverviewResponse.
type OrgOverviewMember struct {
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name,omitempty"`
	Role        string `json:"role"`
}

// OrgOverviewAudit is an audit-log entry within an OrgOverviewResponse.
type OrgOverviewAudit struct {
	ID         string `json:"id"`
	ActionType string `json:"action_type"`
	CreatedAt  string `json:"created_at"`
}

// OrgCreateRequest is the JSON body for POST /v1/orgs.
type OrgCreateRequest struct {
	Name string `json:"name"`
}

// OrgDeleteRequest is the JSON body for POST /v1/orgs/delete.
type OrgDeleteRequest struct {
	OrgID string `json:"org_id"`
}

// InviteRequest is the JSON body for POST /v1/orgs/members (invite).
type InviteRequest struct {
	OrgID  string `json:"org_id"`
	Email  string `json:"email,omitempty"`
	UserID string `json:"user_id,omitempty"`
	Role   string `json:"role,omitempty"`
}

// RoleChangeRequest is the JSON body for POST /v1/orgs/members/role.
type RoleChangeRequest struct {
	OrgID  string `json:"org_id"`
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

// MembershipResponse represents a membership or invite record.
type MembershipResponse struct {
	ID          string `json:"id"`
	OrgID       string `json:"org_id"`
	UserID      string `json:"user_id"`
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Role        string `json:"role"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at,omitempty"`
}

// unexported request bodies ---------------------------------------------------

type acceptInviteRequest struct {
	InviteID string `json:"invite_id"`
}

type declineInviteRequest struct {
	InviteID string `json:"invite_id"`
}

type removeMemberRequest struct {
	OrgID  string `json:"org_id"`
	UserID string `json:"user_id"`
}

type transferOrgRequest struct {
	OrgID          string `json:"org_id"`
	NewOwnerUserID string `json:"new_owner_user_id"`
}

type quitRequest struct {
	OrgID string `json:"org_id"`
}

// ---------------------------------------------------------------------------
// OrgsAPI
// ---------------------------------------------------------------------------

// OrgsAPI provides access to organisation management endpoints.
type OrgsAPI struct {
	client *InoueClient
}

// List returns a paginated list of organisations the current user belongs to.
func (a *OrgsAPI) List(ctx context.Context, page, pageSize int) (*Page[OrgResponse], error) {
	path := fmt.Sprintf("/v1/orgs?page=%d&page_size=%d", page, pageSize)

	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", path, nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("list orgs: %w", err)
	}

	var result Page[OrgResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list orgs decode: %w", err)
	}
	return &result, nil
}

// Get returns a single organisation by ID.
func (a *OrgsAPI) Get(ctx context.Context, orgID string) (*OrgResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", "/v1/orgs/"+orgID, nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("get org: %w", err)
	}

	var result OrgResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get org decode: %w", err)
	}
	return &result, nil
}

// Overview returns an enriched view of an organisation including its members
// and recent audit logs.
func (a *OrgsAPI) Overview(ctx context.Context, orgID string) (*OrgOverviewResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", "/v1/orgs/"+orgID+"/overview", nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("org overview: %w", err)
	}

	var result OrgOverviewResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("org overview decode: %w", err)
	}
	return &result, nil
}

// Create creates a new organisation.
func (a *OrgsAPI) Create(ctx context.Context, req OrgCreateRequest) (*OrgResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/orgs", req, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("create org: %w", err)
	}

	var result OrgResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create org decode: %w", err)
	}
	return &result, nil
}

// InviteMember sends an invitation to join an organisation.
func (a *OrgsAPI) InviteMember(ctx context.Context, req InviteRequest) (*MembershipResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/orgs/members", req, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("invite member: %w", err)
	}

	var result MembershipResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("invite member decode: %w", err)
	}
	return &result, nil
}

// ListMembers returns a paginated list of members for an organisation.
func (a *OrgsAPI) ListMembers(ctx context.Context, orgID string, page, pageSize int) (*Page[MembershipResponse], error) {
	path := fmt.Sprintf("/v1/orgs/%s/members?page=%d&page_size=%d", orgID, page, pageSize)

	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", path, nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}

	var result Page[MembershipResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list members decode: %w", err)
	}
	return &result, nil
}

// ListInvites returns a paginated list of pending invites for an organisation.
func (a *OrgsAPI) ListInvites(ctx context.Context, orgID string, page, pageSize int) (*Page[MembershipResponse], error) {
	path := fmt.Sprintf("/v1/orgs/%s/invites?page=%d&page_size=%d", orgID, page, pageSize)

	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", path, nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("list invites: %w", err)
	}

	var result Page[MembershipResponse]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list invites decode: %w", err)
	}
	return &result, nil
}

// ChangeRole updates the role of a member within an organisation.
func (a *OrgsAPI) ChangeRole(ctx context.Context, req RoleChangeRequest) (*MembershipResponse, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/orgs/members/role", req, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("change role: %w", err)
	}

	var result MembershipResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("change role decode: %w", err)
	}
	return &result, nil
}

// AcceptInvite accepts a pending organisation invite.
func (a *OrgsAPI) AcceptInvite(ctx context.Context, inviteID string) (*MembershipResponse, error) {
	body := acceptInviteRequest{InviteID: inviteID}

	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/orgs/invites/accept", body, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("accept invite: %w", err)
	}

	var result MembershipResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("accept invite decode: %w", err)
	}
	return &result, nil
}

// DeclineInvite declines a pending organisation invite.
func (a *OrgsAPI) DeclineInvite(ctx context.Context, inviteID string) (*MembershipResponse, error) {
	body := declineInviteRequest{InviteID: inviteID}

	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/orgs/invites/decline", body, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("decline invite: %w", err)
	}

	var result MembershipResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("decline invite decode: %w", err)
	}
	return &result, nil
}

// RemoveMember removes a user from an organisation.
func (a *OrgsAPI) RemoveMember(ctx context.Context, orgID, userID string) error {
	body := removeMemberRequest{OrgID: orgID, UserID: userID}

	err := a.client.request(ctx, "POST", "/v1/orgs/members/remove", body, nil, nil)
	if err != nil {
		return fmt.Errorf("remove member: %w", err)
	}
	return nil
}

// Transfer transfers ownership of an organisation to another user.
func (a *OrgsAPI) Transfer(ctx context.Context, orgID, newOwnerUserID string) (*OrgResponse, error) {
	body := transferOrgRequest{OrgID: orgID, NewOwnerUserID: newOwnerUserID}

	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/orgs/transfer", body, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("transfer org: %w", err)
	}

	var result OrgResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("transfer org decode: %w", err)
	}
	return &result, nil
}

// Quit removes the current user from an organisation (self-removal).
func (a *OrgsAPI) Quit(ctx context.Context, orgID string) error {
	body := quitRequest{OrgID: orgID}

	err := a.client.request(ctx, "POST", "/v1/orgs/members/quit", body, nil, nil)
	if err != nil {
		return fmt.Errorf("quit org: %w", err)
	}
	return nil
}

// Delete permanently deletes an organisation.
func (a *OrgsAPI) Delete(ctx context.Context, req OrgDeleteRequest) error {
	err := a.client.request(ctx, "POST", "/v1/orgs/delete", req, nil, nil)
	if err != nil {
		return fmt.Errorf("delete org: %w", err)
	}
	return nil
}

// OrgStorage returns storage usage information for an organisation.
func (a *OrgsAPI) OrgStorage(ctx context.Context, orgID string) (map[string]interface{}, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", "/v1/orgs/"+orgID+"/storage", nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("org storage: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("org storage decode: %w", err)
	}
	return result, nil
}
