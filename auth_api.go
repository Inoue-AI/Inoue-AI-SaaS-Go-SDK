package inouesdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// ---------------------------------------------------------------------------
// Auth types
// ---------------------------------------------------------------------------

// TokenPairResponse contains the access and refresh tokens returned by
// the register, login, and refresh endpoints.
type TokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// AuthMeOrg represents an organisation membership entry inside AuthMeResult.
type AuthMeOrg struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// AuthMeResult is the profile payload returned by GET /v1/auth/me.
type AuthMeResult struct {
	ID               string      `json:"id"`
	Email            string      `json:"email"`
	DisplayName      string      `json:"display_name,omitempty"`
	IsAdmin          bool        `json:"is_admin"`
	TwoFactorEnabled bool        `json:"two_factor_enabled"`
	Orgs             []AuthMeOrg `json:"orgs"`
}

// TwoFactorSetupResult is returned when setting up 2FA for a user.
type TwoFactorSetupResult struct {
	Secret        string   `json:"secret"`
	QRCode        string   `json:"qr_code"`
	RecoveryCodes []string `json:"recovery_codes"`
}

// TwoFactorVerifyResult is returned after successfully verifying a 2FA code.
type TwoFactorVerifyResult struct {
	RecoveryCodes []string `json:"recovery_codes"`
}

// UserSession represents a single authenticated session for the current user.
type UserSession struct {
	ID           string `json:"id"`
	DeviceName   string `json:"device_name,omitempty"`
	IPAddress    string `json:"ip_address,omitempty"`
	UserAgent    string `json:"user_agent,omitempty"`
	CreatedAt    string `json:"created_at"`
	LastActiveAt string `json:"last_active_at,omitempty"`
}

// UpdateMeRequest is the JSON body for PATCH /v1/auth/me.
type UpdateMeRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
}

// unexported request bodies ---------------------------------------------------

type registerRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name,omitempty"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type setup2FARequest struct {
	UserID string `json:"user_id,omitempty"`
}

type verify2FARequest struct {
	UserID string `json:"user_id,omitempty"`
	Code   string `json:"code,omitempty"`
}

type renameSessionRequest struct {
	DeviceName string `json:"device_name"`
}

// ---------------------------------------------------------------------------
// AuthAPI
// ---------------------------------------------------------------------------

// AuthAPI provides access to authentication endpoints.
type AuthAPI struct {
	client *InoueClient
}

// Register creates a new user account and returns a token pair.
func (a *AuthAPI) Register(ctx context.Context, email, password, displayName string) (*TokenPairResponse, error) {
	body := registerRequest{Email: email, Password: password, DisplayName: displayName}

	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/auth/register", body, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	var result TokenPairResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("register decode: %w", err)
	}
	return &result, nil
}

// Login authenticates with email and password and returns a token pair.
func (a *AuthAPI) Login(ctx context.Context, email, password string) (*TokenPairResponse, error) {
	body := loginRequest{Email: email, Password: password}

	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/auth/login", body, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	var result TokenPairResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("login decode: %w", err)
	}
	return &result, nil
}

// Refresh exchanges a refresh token for a new token pair.
func (a *AuthAPI) Refresh(ctx context.Context, refreshToken string) (*TokenPairResponse, error) {
	body := refreshRequest{RefreshToken: refreshToken}

	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/auth/refresh", body, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("refresh: %w", err)
	}

	var result TokenPairResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("refresh decode: %w", err)
	}
	return &result, nil
}

// Logout invalidates the current session on the server side.
func (a *AuthAPI) Logout(ctx context.Context) error {
	err := a.client.request(ctx, "POST", "/v1/auth/logout", nil, nil, nil)
	if err != nil {
		return fmt.Errorf("logout: %w", err)
	}
	return nil
}

// Me returns the profile of the currently authenticated user.
func (a *AuthAPI) Me(ctx context.Context) (*AuthMeResult, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", "/v1/auth/me", nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("me: %w", err)
	}

	var result AuthMeResult
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("me decode: %w", err)
	}
	return &result, nil
}

// UpdateMe patches the current user's profile.
func (a *AuthAPI) UpdateMe(ctx context.Context, req UpdateMeRequest) (*AuthMeResult, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "PATCH", "/v1/auth/me", req, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("update me: %w", err)
	}

	var result AuthMeResult
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("update me decode: %w", err)
	}
	return &result, nil
}

// RegisterSettings returns the public registration configuration such as
// whether registration is open, allowed domains, etc.
func (a *AuthAPI) RegisterSettings(ctx context.Context) (map[string]interface{}, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", "/v1/auth/register/settings", nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("register settings: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("register settings decode: %w", err)
	}
	return result, nil
}

// Setup2FA initiates two-factor authentication setup for the given user.
func (a *AuthAPI) Setup2FA(ctx context.Context, userID string) (*TwoFactorSetupResult, error) {
	body := setup2FARequest{UserID: userID}

	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/auth/2fa/setup", body, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("setup 2fa: %w", err)
	}

	var result TwoFactorSetupResult
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("setup 2fa decode: %w", err)
	}
	return &result, nil
}

// Verify2FA verifies a TOTP code and completes the 2FA enrollment.
func (a *AuthAPI) Verify2FA(ctx context.Context, userID, code string) (*TwoFactorVerifyResult, error) {
	body := verify2FARequest{UserID: userID, Code: code}

	var apiResp ApiResponse
	err := a.client.request(ctx, "POST", "/v1/auth/2fa/verify", body, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("verify 2fa: %w", err)
	}

	var result TwoFactorVerifyResult
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("verify 2fa decode: %w", err)
	}
	return &result, nil
}

// ListSessions returns a paginated list of active sessions for the current user.
func (a *AuthAPI) ListSessions(ctx context.Context, page, pageSize int) (*Page[UserSession], error) {
	path := fmt.Sprintf("/v1/auth/me/sessions?page=%d&page_size=%d", page, pageSize)

	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", path, nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}

	var result Page[UserSession]
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("list sessions decode: %w", err)
	}
	return &result, nil
}

// RenameSession updates the device name of a session.
func (a *AuthAPI) RenameSession(ctx context.Context, sessionID, deviceName string) (*UserSession, error) {
	body := renameSessionRequest{DeviceName: deviceName}

	var apiResp ApiResponse
	err := a.client.request(ctx, "PATCH", "/v1/auth/me/sessions/"+sessionID, body, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("rename session: %w", err)
	}

	var result UserSession
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("rename session decode: %w", err)
	}
	return &result, nil
}

// RevokeSession invalidates the given session, logging it out.
func (a *AuthAPI) RevokeSession(ctx context.Context, sessionID string) error {
	err := a.client.request(ctx, "POST", "/v1/auth/me/sessions/"+sessionID+"/revoke", nil, nil, nil)
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	return nil
}

// MeAnalytics returns usage analytics for the current user.
func (a *AuthAPI) MeAnalytics(ctx context.Context) (map[string]interface{}, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", "/v1/auth/me/analytics", nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("me analytics: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("me analytics decode: %w", err)
	}
	return result, nil
}

// MeStorage returns storage usage information for the current user.
func (a *AuthAPI) MeStorage(ctx context.Context) (map[string]interface{}, error) {
	var apiResp ApiResponse
	err := a.client.request(ctx, "GET", "/v1/auth/me/storage", nil, &apiResp, nil)
	if err != nil {
		return nil, fmt.Errorf("me storage: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("me storage decode: %w", err)
	}
	return result, nil
}
