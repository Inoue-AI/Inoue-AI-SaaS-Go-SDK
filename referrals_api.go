package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// ReferralsAPI provides access to the referral program endpoints.
type ReferralsAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// ReferralCodeCreateRequest is the body for creating a referral code.
type ReferralCodeCreateRequest struct {
	Code string `json:"code,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// ReferralSummaryResponse represents the user's referral summary.
type ReferralSummaryResponse struct {
	TotalReferrals  int     `json:"total_referrals"`
	ActiveReferrals int     `json:"active_referrals"`
	TotalEarnings   float64 `json:"total_earnings"`
	PendingEarnings float64 `json:"pending_earnings"`
	Currency        string  `json:"currency"`
}

// ReferralCodeResponse represents a referral code.
type ReferralCodeResponse struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	UserID    string `json:"user_id"`
	Uses      int    `json:"uses"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// Summary returns the referral summary for the current user.
func (a *ReferralsAPI) Summary(ctx context.Context) (*ReferralSummaryResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/referrals/summary", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("referral summary: %w", err)
	}
	var result ReferralSummaryResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("referral summary decode: %w", err)
	}
	return &result, nil
}

// ListCodes returns the list of referral codes.
func (a *ReferralsAPI) ListCodes(ctx context.Context) ([]ReferralCodeResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/referrals/codes", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list referral codes: %w", err)
	}
	var result []ReferralCodeResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list referral codes decode: %w", err)
	}
	return result, nil
}

// CreateCode creates a new referral code.
func (a *ReferralsAPI) CreateCode(ctx context.Context, req ReferralCodeCreateRequest) (*ReferralCodeResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/referrals/codes", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("create referral code: %w", err)
	}
	var result ReferralCodeResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("create referral code decode: %w", err)
	}
	return &result, nil
}

// RevokeCode revokes a referral code.
func (a *ReferralsAPI) RevokeCode(ctx context.Context, codeID string) error {
	body := map[string]string{"code_id": codeID}
	if err := a.client.request(ctx, "POST", "/v1/referrals/codes/revoke", body, nil, nil); err != nil {
		return fmt.Errorf("revoke referral code: %w", err)
	}
	return nil
}
