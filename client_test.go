package inouesdk

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	c := NewClient("https://api.example.com")
	if c.baseURL != "https://api.example.com" {
		t.Errorf("expected base URL %q, got %q", "https://api.example.com", c.baseURL)
	}
	if c.Internal == nil {
		t.Fatal("expected Internal API to be initialised")
	}
	if c.Schedule == nil {
		t.Fatal("expected Schedule API to be initialised")
	}
}

func TestNewClientWithOptions(t *testing.T) {
	c := NewClient("https://api.example.com/",
		WithAccessToken("tok_abc"),
		WithTimeout(10*time.Second),
	)
	if c.baseURL != "https://api.example.com" {
		t.Errorf("expected trailing slash to be trimmed, got %q", c.baseURL)
	}
	if c.accessToken != "tok_abc" {
		t.Errorf("expected access token %q, got %q", "tok_abc", c.accessToken)
	}
	if c.httpClient.Timeout != 10*time.Second {
		t.Errorf("expected timeout %v, got %v", 10*time.Second, c.httpClient.Timeout)
	}
}

func TestSetAccessToken(t *testing.T) {
	c := NewClient("https://api.example.com")
	c.SetAccessToken("new_token")
	if c.accessToken != "new_token" {
		t.Errorf("expected access token %q, got %q", "new_token", c.accessToken)
	}
}

func TestSdkErrorFormat(t *testing.T) {
	err := &SdkError{
		Code:    "not_found",
		Message: "resource not found",
		Status:  404,
		TraceID: "abc-123",
	}
	expected := "inoue-sdk: resource not found (code=not_found, status=404, trace=abc-123)"
	if err.Error() != expected {
		t.Errorf("expected error string %q, got %q", expected, err.Error())
	}
}

func TestSdkErrorFormatWithoutTraceID(t *testing.T) {
	err := &SdkError{
		Code:    "bad_request",
		Message: "invalid input",
		Status:  400,
	}
	expected := "inoue-sdk: invalid input (code=bad_request, status=400)"
	if err.Error() != expected {
		t.Errorf("expected error string %q, got %q", expected, err.Error())
	}
}

func TestGenerateTraceID(t *testing.T) {
	id := generateTraceID()
	if len(id) != 36 {
		t.Errorf("expected UUID length 36, got %d: %q", len(id), id)
	}
	if id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		t.Errorf("expected UUID format with dashes, got %q", id)
	}
}
