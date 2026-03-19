package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// AnalyticsAPI provides access to the analytics endpoints.
type AnalyticsAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// AnalyticsRevenueResponse represents revenue analytics data.
type AnalyticsRevenueResponse struct {
	ModelID       string                   `json:"model_id"`
	OrgID         string                   `json:"org_id"`
	TotalRevenue  float64                  `json:"total_revenue"`
	Currency      string                   `json:"currency"`
	Breakdown     []map[string]interface{} `json:"breakdown"`
	Period        string                   `json:"period"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// Revenue returns revenue analytics for a specific model.
func (a *AnalyticsAPI) Revenue(ctx context.Context, modelID string) (*AnalyticsRevenueResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/analytics/revenue/models/"+modelID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("model revenue analytics: %w", err)
	}
	var result AnalyticsRevenueResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("model revenue analytics decode: %w", err)
	}
	return &result, nil
}

// OrgRevenue returns revenue analytics for an organization.
func (a *AnalyticsAPI) OrgRevenue(ctx context.Context, orgID string) (*AnalyticsRevenueResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/analytics/revenue/orgs/"+orgID, nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("org revenue analytics: %w", err)
	}
	var result AnalyticsRevenueResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("org revenue analytics decode: %w", err)
	}
	return &result, nil
}
