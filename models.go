package inouesdk

import "encoding/json"

// ApiResponse represents the standard envelope returned by all Inoue AI API endpoints.
type ApiResponse struct {
	Ok    bool                   `json:"ok"`
	Data  json.RawMessage        `json:"data"`
	Error *ApiErrorPayload       `json:"error,omitempty"`
	Meta  map[string]interface{} `json:"meta,omitempty"`
}

// ApiErrorPayload carries the structured error information inside an ApiResponse.
type ApiErrorPayload struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Page represents a paginated list of items returned by the API.
type Page[T any] struct {
	Items      []T `json:"items"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}
