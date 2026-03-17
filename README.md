# Inoue AI SaaS Go SDK

Go client library for the Inoue AI SaaS Backend API, purpose-built for Go-based workers in the Inoue AI platform.

The SDK provides a typed, idiomatic Go interface to the Backend's internal worker management and content scheduling endpoints. It handles JSON envelope parsing, automatic UUID v4 trace ID propagation, structured error extraction, and Bearer token lifecycle — all with **zero external dependencies** (Go standard library only).

## Installation

```bash
go get github.com/Inoue-AI/Inoue-AI-SaaS-Go-SDK
```

Requires **Go 1.22+**.

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	inouesdk "github.com/Inoue-AI/Inoue-AI-SaaS-Go-SDK"
)

func main() {
	ctx := context.Background()

	// Create a client pointed at the Backend API
	client := inouesdk.NewClient("https://api.inoue.ai")

	// Register this worker (idempotent — 409 is treated as success)
	err := client.Internal.RegisterWorker(ctx, "worker-uuid", "my-worker", "admin-token")
	if err != nil {
		log.Fatal(err)
	}

	// Obtain a worker JWT using the bootstrap secret
	token, err := client.Internal.WorkerToken(ctx, "worker-uuid", "bootstrap-secret")
	if err != nil {
		log.Fatal(err)
	}

	// Set the token for all subsequent authenticated requests
	client.SetAccessToken(token)

	// Report worker status
	err = client.Internal.SetWorkerStatus(ctx, "active")
	if err != nil {
		log.Fatal(err)
	}

	// Create a schedule entry
	notes := "Campaign launch post"
	entry, err := client.Schedule.Create(ctx, inouesdk.ScheduleEntryCreateParams{
		ModelID:       "model-uuid",
		PlatformID:    "platform-uuid",
		ContentTypeID: "content-type-uuid",
		ScheduledFor:  "2026-03-20T10:00:00Z",
		AssetIDs:      []string{"asset-uuid-1"},
		Notes:         &notes,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created schedule entry: %s (status: %s)\n", entry.ID, entry.Status)
}
```

## API Reference

### Client Construction

| Function / Method | Signature | Purpose |
|---|---|---|
| `NewClient` | `NewClient(baseURL string, opts ...Option) *InoueClient` | Create a new client pointed at the given base URL |
| `WithAccessToken` | `WithAccessToken(token string) Option` | Set the initial Bearer token during construction |
| `WithTimeout` | `WithTimeout(d time.Duration) Option` | Set the HTTP client timeout (default: 30s) |
| `WithHTTPClient` | `WithHTTPClient(client *http.Client) Option` | Replace the underlying `http.Client` entirely |
| `SetAccessToken` | `(c *InoueClient) SetAccessToken(token string)` | Update the Bearer token for subsequent requests |

### Internal API (`client.Internal`)

| Method | Signature | Purpose |
|---|---|---|
| `RegisterWorker` | `(a *InternalAPI) RegisterWorker(ctx context.Context, workerID, name, adminToken string) error` | Register a worker with the backend; 409 (already registered) is treated as success |
| `WorkerToken` | `(a *InternalAPI) WorkerToken(ctx context.Context, workerID, bootstrapSecret string) (string, error)` | Obtain a worker JWT using the bootstrap secret |
| `SetWorkerStatus` | `(a *InternalAPI) SetWorkerStatus(ctx context.Context, status string) error` | Update this worker's status (e.g. `"active"`, `"disabled"`) |

### Schedule API (`client.Schedule`)

| Method | Signature | Purpose |
|---|---|---|
| `Create` | `(a *ScheduleAPI) Create(ctx context.Context, params ScheduleEntryCreateParams) (*ScheduleEntry, error)` | Create a new content schedule entry |

### Types

| Type | Purpose |
|---|---|
| `InoueClient` | Top-level client; holds `Internal` and `Schedule` sub-API accessors |
| `Option` | Functional option for configuring the client during construction |
| `InternalAPI` | Worker management endpoints (register, token, status) |
| `ScheduleAPI` | Content scheduling endpoints |
| `ApiResponse` | Standard JSON envelope: `{ok, data, error, meta}` |
| `ApiErrorPayload` | Structured error inside `ApiResponse`: `{code, message, details}` |
| `Page[T]` | Generic paginated list: `{items, total, page, page_size, total_pages}` |
| `ScheduleEntryCreateParams` | Parameters for creating a schedule entry |
| `ScheduleEntry` | A scheduled content entry returned by the API |
| `SdkError` | Structured error with HTTP status, error code, message, trace ID, and details |

## Configuration

The SDK does not read environment variables. All configuration is passed explicitly at construction time:

| Parameter | Passed Via | Default | Description |
|---|---|---|---|
| Base URL | `NewClient(baseURL)` | (required) | Backend API base URL |
| Access token | `WithAccessToken()` or `SetAccessToken()` | `""` (unauthenticated) | Bearer token for `Authorization` header |
| HTTP timeout | `WithTimeout()` | 30 seconds | Timeout for all HTTP requests |
| HTTP client | `WithHTTPClient()` | Default `http.Client` | Full control over transport, TLS, proxy, etc. |

## Error Handling

All API methods return errors as `*SdkError` when the server responds with HTTP status >= 400. The error contains structured information extracted from the API's JSON envelope:

```go
entry, err := client.Schedule.Create(ctx, params)
if err != nil {
	var sdkErr *inouesdk.SdkError
	if errors.As(err, &sdkErr) {
		fmt.Printf("API error: code=%s status=%d trace=%s\n",
			sdkErr.Code, sdkErr.Status, sdkErr.TraceID)
		fmt.Printf("Message: %s\n", sdkErr.Message)
		if sdkErr.Details != nil {
			fmt.Printf("Details: %v\n", sdkErr.Details)
		}
	} else {
		// Transport-level error (DNS, timeout, etc.)
		fmt.Printf("Transport error: %v\n", err)
	}
}
```

`SdkError` fields:

| Field | Type | Description |
|---|---|---|
| `Code` | `string` | Application-level error code (e.g. `"not_found"`, `"validation_error"`) |
| `Message` | `string` | Human-readable error message |
| `Status` | `int` | HTTP status code (400, 401, 404, 409, 500, etc.) |
| `TraceID` | `string` | UUID v4 trace ID for the request, useful for log correlation |
| `Details` | `map[string]interface{}` | Optional additional error details from the API |

`SdkError` implements the `error` interface. Its `Error()` method returns a formatted string:
```
inoue-sdk: resource not found (code=not_found, status=404, trace=abc-123)
```

Non-API errors (JSON marshaling failures, network errors, response decoding errors) are returned as standard wrapped errors with the `inoue-sdk:` prefix.

## Development

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run vet
go vet ./...

# Run tests with race detector
go test -race ./...
```

## License

MIT
