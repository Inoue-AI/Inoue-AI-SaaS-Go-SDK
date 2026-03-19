package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// BillingAPI provides access to the billing and subscription endpoints.
type BillingAPI struct {
	client *InoueClient
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// BillingCheckoutCreateRequest is the body for checkout endpoints.
type BillingCheckoutCreateRequest struct {
	PriceID    string `json:"price_id"`
	Quantity   int    `json:"quantity,omitempty"`
	SuccessURL string `json:"success_url,omitempty"`
	CancelURL  string `json:"cancel_url,omitempty"`
}

// BillingPortalCreateRequest is the body for POST /v1/billing/portal.
type BillingPortalCreateRequest struct {
	ReturnURL string `json:"return_url,omitempty"`
}

// BillingSubscriptionChangeRequest is the body for POST /v1/billing/subscription/change.
type BillingSubscriptionChangeRequest struct {
	PriceID  string `json:"price_id"`
	Quantity int    `json:"quantity,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// BillingProduct represents a product available for purchase.
type BillingProduct struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Type        string                   `json:"type"`
	Prices      []map[string]interface{} `json:"prices"`
	Metadata    map[string]interface{}   `json:"metadata"`
}

// BillingSummary represents the user's billing summary.
type BillingSummary struct {
	CustomerID     string                 `json:"customer_id"`
	Subscription   map[string]interface{} `json:"subscription"`
	CurrentPlan    string                 `json:"current_plan"`
	BillingEmail   string                 `json:"billing_email"`
	PaymentMethod  map[string]interface{} `json:"payment_method"`
	NextInvoice    map[string]interface{} `json:"next_invoice"`
	CreditBalance  float64                `json:"credit_balance"`
}

// BillingCheckoutCreateResult is the response from a checkout creation.
type BillingCheckoutCreateResult struct {
	CheckoutURL string `json:"checkout_url"`
	SessionID   string `json:"session_id"`
}

// BillingPortalCreateResult is the response from a portal session creation.
type BillingPortalCreateResult struct {
	PortalURL string `json:"portal_url"`
}

// BillingSubscriptionChangeResult is the response from a subscription change.
type BillingSubscriptionChangeResult struct {
	SubscriptionID string `json:"subscription_id"`
	Status         string `json:"status"`
	PriceID        string `json:"price_id"`
}

// ---------------------------------------------------------------------------
// Methods
// ---------------------------------------------------------------------------

// Products returns the list of available billing products.
func (a *BillingAPI) Products(ctx context.Context) ([]BillingProduct, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/billing/products", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("billing products: %w", err)
	}
	var result []BillingProduct
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("billing products decode: %w", err)
	}
	return result, nil
}

// Summary returns the current user's billing summary.
func (a *BillingAPI) Summary(ctx context.Context) (*BillingSummary, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "GET", "/v1/billing/summary", nil, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("billing summary: %w", err)
	}
	var result BillingSummary
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("billing summary decode: %w", err)
	}
	return &result, nil
}

// CheckoutTopup creates a top-up checkout session.
func (a *BillingAPI) CheckoutTopup(ctx context.Context, req BillingCheckoutCreateRequest) (*BillingCheckoutCreateResult, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/billing/checkout/topup", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("checkout topup: %w", err)
	}
	var result BillingCheckoutCreateResult
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("checkout topup decode: %w", err)
	}
	return &result, nil
}

// CheckoutSubscription creates a subscription checkout session.
func (a *BillingAPI) CheckoutSubscription(ctx context.Context, req BillingCheckoutCreateRequest) (*BillingCheckoutCreateResult, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/billing/checkout/subscription", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("checkout subscription: %w", err)
	}
	var result BillingCheckoutCreateResult
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("checkout subscription decode: %w", err)
	}
	return &result, nil
}

// Portal creates a billing portal session.
func (a *BillingAPI) Portal(ctx context.Context, req BillingPortalCreateRequest) (*BillingPortalCreateResult, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/billing/portal", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("billing portal: %w", err)
	}
	var result BillingPortalCreateResult
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("billing portal decode: %w", err)
	}
	return &result, nil
}

// ChangeSubscription changes the user's active subscription.
func (a *BillingAPI) ChangeSubscription(ctx context.Context, req BillingSubscriptionChangeRequest) (*BillingSubscriptionChangeResult, error) {
	var apiResp ApiResponse
	if err := a.client.request(ctx, "POST", "/v1/billing/subscription/change", req, &apiResp, nil); err != nil {
		return nil, fmt.Errorf("change subscription: %w", err)
	}
	var result BillingSubscriptionChangeResult
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("change subscription decode: %w", err)
	}
	return &result, nil
}
