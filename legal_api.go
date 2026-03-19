package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// LegalAPI provides access to the legal document endpoints.
type LegalAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// LegalDocumentResponse represents a legal document.
type LegalDocumentResponse struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Version   string `json:"version"`
	UpdatedAt string `json:"updated_at"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// Terms returns the terms of service document.
func (a *LegalAPI) Terms(ctx context.Context) (*LegalDocumentResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/legal/terms", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get terms: %w", err)
	}
	var result LegalDocumentResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get terms decode: %w", err)
	}
	return &result, nil
}

// Privacy returns the privacy policy document.
func (a *LegalAPI) Privacy(ctx context.Context) (*LegalDocumentResponse, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/legal/privacy", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("get privacy: %w", err)
	}
	var result LegalDocumentResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("get privacy decode: %w", err)
	}
	return &result, nil
}
