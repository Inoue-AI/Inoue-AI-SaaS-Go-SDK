package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// CreditsAPI provides access to the credits and usage endpoints.
type CreditsAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// UserWallet represents a user's credit wallet.
type UserWallet struct {
	ID              string  `json:"id"`
	UserID          string  `json:"user_id"`
	Balance         float64 `json:"balance"`
	LifetimeCredits float64 `json:"lifetime_credits"`
	Currency        string  `json:"currency"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// UsageRecord represents a single credit usage record.
type UsageRecord struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	OrgID       string  `json:"org_id"`
	Action      string  `json:"action"`
	Credits     float64 `json:"credits"`
	Description string  `json:"description"`
	ReferenceID string  `json:"reference_id"`
	CreatedAt   string  `json:"created_at"`
}

// RateCardEntry represents a single rate card item.
type RateCardEntry struct {
	Action      string  `json:"action"`
	Credits     float64 `json:"credits"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// Wallet returns the current user's credit wallet.
func (a *CreditsAPI) Wallet(ctx context.Context) (*UserWallet, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/credits/wallet", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get wallet: %w", err)
	}
	var result UserWallet
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get wallet decode: %w", err)
	}
	return &result, nil
}

// Usage returns a paginated list of credit usage records.
func (a *CreditsAPI) Usage(ctx context.Context, page, pageSize int) (*Page[UsageRecord], error) {
	path := fmt.Sprintf("/v1/credits/usage?page=%d&page_size=%d", page, pageSize)
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", path, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("list usage: %w", err)
	}
	var result Page[UsageRecord]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list usage decode: %w", err)
	}
	return &result, nil
}

// UsageSeries returns time-series usage data.
func (a *CreditsAPI) UsageSeries(ctx context.Context) (map[string]interface{}, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/credits/usage/series", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("usage series: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("usage series decode: %w", err)
	}
	return result, nil
}

// RateCards returns the list of credit rate cards.
func (a *CreditsAPI) RateCards(ctx context.Context) ([]RateCardEntry, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/credits/rate-cards", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("rate cards: %w", err)
	}
	var result []RateCardEntry
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("rate cards decode: %w", err)
	}
	return result, nil
}
