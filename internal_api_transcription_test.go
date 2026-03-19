package inouesdk

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClientWithServer(handler http.HandlerFunc) (*InoueClient, *httptest.Server) {
	srv := httptest.NewServer(handler)
	client := NewClient(srv.URL, WithAccessToken("test-token"))
	return client, srv
}

func TestGetTranscription(t *testing.T) {
	client, srv := newTestClientWithServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/internal/transcription-requests/req-123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing or wrong Authorization header")
		}

		resp := ApiResponse{
			Ok: true,
			Data: json.RawMessage(`{
				"id": "req-123",
				"caption_project_id": "proj-456",
				"asset_id": "asset-789",
				"status": "queued",
				"provider": "deepgram",
				"created_at": "2026-03-18T10:00:00Z",
				"updated_at": "2026-03-18T10:00:00Z"
			}`),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	defer srv.Close()

	req, err := client.Internal.GetTranscription(context.Background(), "req-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.ID != "req-123" {
		t.Errorf("expected ID %q, got %q", "req-123", req.ID)
	}
	if req.CaptionProjectID != "proj-456" {
		t.Errorf("expected CaptionProjectID %q, got %q", "proj-456", req.CaptionProjectID)
	}
	if req.Status != "queued" {
		t.Errorf("expected status %q, got %q", "queued", req.Status)
	}
	if req.Provider != "deepgram" {
		t.Errorf("expected provider %q, got %q", "deepgram", req.Provider)
	}
}

func TestClaimTranscription(t *testing.T) {
	client, srv := newTestClientWithServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/internal/transcription-requests/req-123/claim" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := ApiResponse{
			Ok: true,
			Data: json.RawMessage(`{
				"id": "req-123",
				"status": "running",
				"claimed_by_worker_id": "worker-1",
				"claimed_at": "2026-03-18T10:01:00Z"
			}`),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	defer srv.Close()

	claim, err := client.Internal.ClaimTranscription(context.Background(), "req-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claim.Status != "running" {
		t.Errorf("expected status %q, got %q", "running", claim.Status)
	}
	if claim.ClaimedByWorker != "worker-1" {
		t.Errorf("expected worker %q, got %q", "worker-1", claim.ClaimedByWorker)
	}
}

func TestClaimTranscriptionConflict(t *testing.T) {
	client, srv := newTestClientWithServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(409)
		resp := ApiResponse{
			Ok:    false,
			Error: &ApiErrorPayload{Code: "conflict", Message: "already claimed"},
		}
		json.NewEncoder(w).Encode(resp)
	})
	defer srv.Close()

	_, err := client.Internal.ClaimTranscription(context.Background(), "req-123")
	if err == nil {
		t.Fatal("expected error for 409 conflict")
	}
	var sdkErr *SdkError
	if !errors.As(err, &sdkErr) {
		t.Fatalf("expected *SdkError in error chain, got %T", err)
	}
	if sdkErr.Status != 409 {
		t.Errorf("expected status 409, got %d", sdkErr.Status)
	}
}

func TestReclaimTranscription(t *testing.T) {
	client, srv := newTestClientWithServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/internal/transcription-requests/req-123/reclaim" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := ApiResponse{
			Ok: true,
			Data: json.RawMessage(`{
				"id": "req-123",
				"status": "running",
				"claimed_by_worker_id": "worker-2",
				"claimed_at": "2026-03-18T10:05:00Z"
			}`),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	defer srv.Close()

	claim, err := client.Internal.ReclaimTranscription(context.Background(), "req-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claim.ClaimedByWorker != "worker-2" {
		t.Errorf("expected worker %q, got %q", "worker-2", claim.ClaimedByWorker)
	}
}

func TestTranscriptionHeartbeat(t *testing.T) {
	client, srv := newTestClientWithServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/internal/transcription-requests/req-123/heartbeat" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(200)
	})
	defer srv.Close()

	err := client.Internal.TranscriptionHeartbeat(context.Background(), "req-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTranscriptionProgress(t *testing.T) {
	var receivedBody map[string]interface{}

	client, srv := newTestClientWithServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/internal/transcription-requests/req-123/progress" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedBody)

		w.WriteHeader(200)
	})
	defer srv.Close()

	req := TranscriptionProgressRequest{
		ProgressJSON: map[string]interface{}{
			"phase":   "transcribing",
			"percent": float64(42),
		},
	}
	err := client.Internal.TranscriptionProgress(context.Background(), "req-123", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	progressJSON, ok := receivedBody["progress_json"].(map[string]interface{})
	if !ok {
		t.Fatal("expected progress_json in request body")
	}
	if progressJSON["phase"] != "transcribing" {
		t.Errorf("expected phase %q, got %v", "transcribing", progressJSON["phase"])
	}
}

func TestCompleteTranscription(t *testing.T) {
	var receivedBody map[string]interface{}

	client, srv := newTestClientWithServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/internal/transcription-requests/req-123/complete" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedBody)

		w.WriteHeader(200)
	})
	defer srv.Close()

	req := TranscriptionCompleteRequest{
		WordsJSON: []TranscriptionWord{
			{Text: "Hello", StartTime: 0.0, EndTime: 0.5, Confidence: 0.99, Speaker: "A"},
			{Text: "world", StartTime: 0.5, EndTime: 1.0, Confidence: 0.95, Speaker: "A"},
		},
		ChunksJSON:       []map[string]interface{}{{"start_time": 0.0, "end_time": 1.0}},
		LanguageDetected: "en",
		DurationSeconds:  1.0,
	}
	err := client.Internal.CompleteTranscription(context.Background(), "req-123", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedBody["language_detected"] != "en" {
		t.Errorf("expected language_detected %q, got %v", "en", receivedBody["language_detected"])
	}
	words, ok := receivedBody["words_json"].([]interface{})
	if !ok {
		t.Fatal("expected words_json array in request body")
	}
	if len(words) != 2 {
		t.Errorf("expected 2 words, got %d", len(words))
	}
}

func TestFailTranscription(t *testing.T) {
	var receivedBody map[string]interface{}

	client, srv := newTestClientWithServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/internal/transcription-requests/req-123/fail" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedBody)

		w.WriteHeader(200)
	})
	defer srv.Close()

	req := TranscriptionFailRequest{
		ErrorJSON: map[string]interface{}{
			"code":    "provider_error",
			"message": "Deepgram returned 503",
		},
	}
	err := client.Internal.FailTranscription(context.Background(), "req-123", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	errorJSON, ok := receivedBody["error_json"].(map[string]interface{})
	if !ok {
		t.Fatal("expected error_json in request body")
	}
	if errorJSON["code"] != "provider_error" {
		t.Errorf("expected code %q, got %v", "provider_error", errorJSON["code"])
	}
}

func TestTranscriptionServerError(t *testing.T) {
	client, srv := newTestClientWithServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		resp := ApiResponse{
			Ok:    false,
			Error: &ApiErrorPayload{Code: "internal_error", Message: "database unavailable"},
		}
		json.NewEncoder(w).Encode(resp)
	})
	defer srv.Close()

	_, err := client.Internal.GetTranscription(context.Background(), "req-123")
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
	var sdkErr *SdkError
	if !errors.As(err, &sdkErr) {
		t.Fatalf("expected *SdkError in error chain, got %T", err)
	}
	if sdkErr.Status != 500 {
		t.Errorf("expected status 500, got %d", sdkErr.Status)
	}
}
