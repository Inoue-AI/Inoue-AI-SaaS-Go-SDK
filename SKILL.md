# SKILL.md — Inoue AI SaaS Go SDK

## What This Service Does

Go client library for the Inoue AI SaaS Backend API. Provides typed access to worker management (register, token minting, status reporting) and content scheduling endpoints. Every HTTP request is wrapped with automatic UUID v4 trace ID generation, JSON envelope parsing, Bearer token injection, and structured error extraction. Built entirely on the Go standard library with zero external dependencies.

## Repository Layout

```
Inoue-AI-SaaS-Go-SDK/
├── go.mod                # Module: github.com/Inoue-AI/Inoue-AI-SaaS-Go-SDK, Go 1.22, zero dependencies
├── client.go             # InoueClient, functional options, HTTP request engine, trace ID generation, error parsing
├── models.go             # ApiResponse envelope, ApiErrorPayload, generic Page[T] type
├── errors.go             # SdkError type implementing the error interface
├── internal_api.go       # InternalAPI: RegisterWorker, WorkerToken, SetWorkerStatus
├── schedule_api.go       # ScheduleAPI: Create, ScheduleEntryCreateParams, ScheduleEntry
├── client_test.go        # Unit tests: client construction, options, token setting, error formatting, trace ID
├── .gitignore            # Binaries, test artifacts, coverage files, .env, editor dirs
└── README.md             # Installation, quick start, API reference, error handling, development
```

## Key Components

| Component | File | Purpose |
|-----------|------|---------|
| `InoueClient` | `client.go` | Top-level client struct; holds base URL, access token, HTTP client, and sub-API accessors (`Internal`, `Schedule`) |
| `Option` | `client.go` | Functional option type (`func(*InoueClient)`) for configuring the client at construction time |
| `NewClient()` | `client.go` | Constructor: trims trailing slash from base URL, applies options, initializes sub-APIs |
| `request()` | `client.go` | Core HTTP engine: marshal body, build request with context, attach headers (Content-Type, Authorization, X-Trace-Id, extras), execute, read response, parse errors or decode JSON into `dest` |
| `parseErrorResponse()` | `client.go` | Extracts `SdkError` from API JSON envelope; falls back to raw body if envelope parsing fails |
| `generateTraceID()` | `client.go` | Produces a UUID v4 string using `crypto/rand` for request tracing |
| `ApiResponse` | `models.go` | Standard API envelope: `{ok, data, error, meta}` with `json.RawMessage` for lazy `data` decoding |
| `ApiErrorPayload` | `models.go` | Error payload inside the envelope: `{code, message, details}` |
| `Page[T]` | `models.go` | Generic paginated list type: `{items, total, page, page_size, total_pages}` |
| `SdkError` | `errors.go` | Structured error type with `Code`, `Message`, `Status`, `TraceID`, `Details`; implements `error` interface |
| `InternalAPI` | `internal_api.go` | Worker management: register, obtain JWT, set status |
| `RegisterWorker()` | `internal_api.go` | POST `/internal/workers/register` with `X-Admin-Token`; treats 409 as success (idempotent) |
| `WorkerToken()` | `internal_api.go` | POST `/internal/workers/token` with `X-Worker-Bootstrap`; unwraps JWT from `ApiResponse.Data` |
| `SetWorkerStatus()` | `internal_api.go` | POST `/internal/workers/me/status` with Bearer token auth |
| `ScheduleAPI` | `schedule_api.go` | Content scheduling endpoints |
| `Create()` | `schedule_api.go` | POST `/v1/schedule/`; returns deserialized `ScheduleEntry` |
| `ScheduleEntryCreateParams` | `schedule_api.go` | Create parameters: `OrgID`, `ModelID`, `PlatformID`, `ContentTypeID`, `ScheduledFor`, `AssetIDs`, `Notes` |
| `ScheduleEntry` | `schedule_api.go` | Response model: `ID`, `Status`, `ModelID`, `PlatformID`, `ContentTypeID`, `ScheduledFor`, `AssetIDs`, `Notes`, `CreatedAt` |

## Design Patterns

- **Functional Options** — `NewClient(baseURL, opts...)` accepts variadic `Option` functions (`WithAccessToken`, `WithTimeout`, `WithHTTPClient`). Options mutate the client struct during construction, allowing callers to configure only what they need without a builder or config struct.

- **Sub-API Accessor Pattern** — Domain-specific endpoints are grouped into dedicated structs (`InternalAPI`, `ScheduleAPI`) attached as public fields on `InoueClient`. This mirrors the Python SDK's `client.internal.*` / `client.schedule.*` pattern and keeps the top-level client clean while allowing each domain to grow independently.

- **JSON Envelope Parsing** — All Backend API responses follow a `{ok, data, error, meta}` envelope (`ApiResponse`). The `data` field uses `json.RawMessage` for lazy deserialization — the request engine decodes the envelope first, checks for errors, then each sub-API method unmarshals `Data` into its specific response type. This avoids generics at the HTTP layer while keeping type safety at the call site.

- **Trace ID Propagation** — Every request generates a UUID v4 trace ID via `crypto/rand` and attaches it as the `X-Trace-Id` header. The same trace ID is embedded in any `SdkError` returned, enabling end-to-end request correlation across SDK caller, Backend API, and downstream services.

- **Idempotent Registration** — `RegisterWorker()` treats HTTP 409 (Conflict) as success rather than an error. This makes the bootstrap flow safe to call on every startup without conditional logic in the caller.

- **Extra Headers Pattern** — The internal `request()` method accepts an `extraHeaders map[string]string` parameter. This allows sub-API methods to inject endpoint-specific headers (`X-Admin-Token`, `X-Worker-Bootstrap`) without exposing header manipulation to SDK consumers.

- **Structured Error Extraction** — On HTTP >= 400, `parseErrorResponse()` attempts to decode the Backend's JSON error envelope into a `SdkError` with typed fields (`Code`, `Message`, `Status`, `TraceID`, `Details`). If the body is not valid JSON or does not match the envelope schema, the raw body is used as the message. This ensures errors are always actionable.

- **Zero Dependencies** — The entire SDK uses only Go standard library packages (`net/http`, `encoding/json`, `crypto/rand`, `context`, `fmt`, `io`, `bytes`, `strings`, `time`). This eliminates dependency management overhead and ensures compatibility across all Go environments.

## API Reference

### Exported Types

| Type | Kind | Package | Description |
|------|------|---------|-------------|
| `InoueClient` | struct | `inouesdk` | Top-level client with `Internal *InternalAPI` and `Schedule *ScheduleAPI` fields |
| `Option` | `func(*InoueClient)` | `inouesdk` | Functional option for client configuration |
| `InternalAPI` | struct | `inouesdk` | Worker management endpoints |
| `ScheduleAPI` | struct | `inouesdk` | Content scheduling endpoints |
| `ApiResponse` | struct | `inouesdk` | JSON envelope: `Ok bool`, `Data json.RawMessage`, `Error *ApiErrorPayload`, `Meta map[string]interface{}` |
| `ApiErrorPayload` | struct | `inouesdk` | Error payload: `Code string`, `Message string`, `Details map[string]interface{}` |
| `Page[T]` | generic struct | `inouesdk` | Paginated list: `Items []T`, `Total int`, `Page int`, `PageSize int`, `TotalPages int` |
| `SdkError` | struct | `inouesdk` | Error type: `Code string`, `Message string`, `Status int`, `TraceID string`, `Details map[string]interface{}` |
| `ScheduleEntryCreateParams` | struct | `inouesdk` | Create params: `OrgID *string`, `ModelID string`, `PlatformID string`, `ContentTypeID string`, `ScheduledFor string`, `AssetIDs []string`, `Notes *string` |
| `ScheduleEntry` | struct | `inouesdk` | Response: `ID string`, `Status string`, `ModelID string`, `PlatformID string`, `ContentTypeID string`, `ScheduledFor string`, `AssetIDs []string`, `Notes *string`, `CreatedAt string` |

### Exported Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `NewClient` | `func NewClient(baseURL string, opts ...Option) *InoueClient` | Create a client; trims trailing slash, default 30s timeout |
| `WithAccessToken` | `func WithAccessToken(token string) Option` | Set initial Bearer token |
| `WithTimeout` | `func WithTimeout(d time.Duration) Option` | Set HTTP client timeout |
| `WithHTTPClient` | `func WithHTTPClient(client *http.Client) Option` | Replace the `http.Client` |

### Exported Methods

| Receiver | Method | Signature | Description |
|----------|--------|-----------|-------------|
| `*InoueClient` | `SetAccessToken` | `func (c *InoueClient) SetAccessToken(token string)` | Update Bearer token for subsequent requests |
| `*InternalAPI` | `RegisterWorker` | `func (a *InternalAPI) RegisterWorker(ctx context.Context, workerID, name, adminToken string) error` | POST `/internal/workers/register`; 409 = success |
| `*InternalAPI` | `WorkerToken` | `func (a *InternalAPI) WorkerToken(ctx context.Context, workerID, bootstrapSecret string) (string, error)` | POST `/internal/workers/token`; returns JWT string |
| `*InternalAPI` | `SetWorkerStatus` | `func (a *InternalAPI) SetWorkerStatus(ctx context.Context, status string) error` | POST `/internal/workers/me/status` |
| `*ScheduleAPI` | `Create` | `func (a *ScheduleAPI) Create(ctx context.Context, params ScheduleEntryCreateParams) (*ScheduleEntry, error)` | POST `/v1/schedule/`; returns created entry |

## Backend API Contract

The Go SDK communicates with the Backend via these endpoints:

| Endpoint | Method | Auth Header | Purpose |
|----------|--------|-------------|---------|
| `/internal/workers/register` | POST | `X-Admin-Token` | Register a worker (idempotent, 409 = already registered) |
| `/internal/workers/token` | POST | `X-Worker-Bootstrap` | Mint a worker JWT using the bootstrap secret |
| `/internal/workers/me/status` | POST | `Authorization: Bearer {token}` | Report worker status (e.g. `"active"`, `"disabled"`) |
| `/v1/schedule/` | POST | `Authorization: Bearer {token}` | Create a new content schedule entry |

### Request/Response Envelope

All endpoints return the standard envelope:

```json
{
  "ok": true,
  "data": { ... },
  "error": null,
  "meta": { ... }
}
```

On error (HTTP >= 400):

```json
{
  "ok": false,
  "data": null,
  "error": {
    "code": "not_found",
    "message": "Resource not found",
    "details": { ... }
  },
  "meta": { ... }
}
```

### Register Worker Request

```json
POST /internal/workers/register
X-Admin-Token: {admin_token}
Content-Type: application/json

{
  "id": "worker-uuid",
  "name": "worker-name"
}
```

### Worker Token Request

```json
POST /internal/workers/token
X-Worker-Bootstrap: {bootstrap_secret}
Content-Type: application/json

{
  "worker_id": "worker-uuid"
}
```

Response `data`:

```json
{
  "worker_token": "eyJhbGciOi..."
}
```

### Set Worker Status Request

```json
POST /internal/workers/me/status
Authorization: Bearer {worker_token}
Content-Type: application/json

{
  "status": "active"
}
```

### Create Schedule Entry Request

```json
POST /v1/schedule/
Authorization: Bearer {token}
Content-Type: application/json

{
  "org_id": "org-uuid",
  "model_id": "model-uuid",
  "platform_id": "platform-uuid",
  "content_type_id": "content-type-uuid",
  "scheduled_for": "2026-03-20T10:00:00Z",
  "asset_ids": ["asset-uuid-1"],
  "notes": "Campaign launch post"
}
```

Response `data`:

```json
{
  "id": "entry-uuid",
  "status": "pending",
  "model_id": "model-uuid",
  "platform_id": "platform-uuid",
  "content_type_id": "content-type-uuid",
  "scheduled_for": "2026-03-20T10:00:00Z",
  "asset_ids": ["asset-uuid-1"],
  "notes": "Campaign launch post",
  "created_at": "2026-03-17T12:00:00Z"
}
```

## Error Handling

The SDK has two error categories:

### 1. API Errors (`*SdkError`)

Returned when the server responds with HTTP status >= 400. The SDK parses the JSON error envelope and produces a `*SdkError` with:

- `Code` — Application-level error code from the API (e.g. `"not_found"`, `"validation_error"`, `"unauthorized"`)
- `Message` — Human-readable error message
- `Status` — HTTP status code
- `TraceID` — UUID v4 trace ID that was sent with the request
- `Details` — Optional structured details from the API

If the response body cannot be parsed as a JSON envelope, `Code` defaults to `"unknown_error"` and `Message` contains the raw response body.

Special cases:
- `RegisterWorker()` treats 409 as success (worker already registered)

### 2. Transport/Encoding Errors

Returned as standard Go errors (wrapped with `fmt.Errorf`) for:
- Request body JSON marshaling failures (`"inoue-sdk: failed to marshal request body"`)
- HTTP request creation failures (`"inoue-sdk: failed to create request"`)
- Network/transport failures (`"inoue-sdk: request failed"`)
- Response body read failures (`"inoue-sdk: failed to read response body"`)
- Response JSON decoding failures (`"inoue-sdk: failed to decode response"`)

All transport errors are wrapped with `%w` for compatibility with `errors.Is()` and `errors.As()`.

## Adding a New API

To extend the SDK with a new service (e.g. a Jobs API):

1. **Create the API file** (e.g. `jobs_api.go`):
   ```go
   package inouesdk

   import (
       "context"
       "encoding/json"
       "fmt"
   )

   // JobsAPI provides access to job management endpoints.
   type JobsAPI struct {
       client *InoueClient
   }

   // GetJob fetches a job by ID.
   func (a *JobsAPI) GetJob(ctx context.Context, jobID string) (*Job, error) {
       var apiResp ApiResponse
       path := fmt.Sprintf("/internal/jobs/%s", jobID)
       err := a.client.request(ctx, "GET", path, nil, &apiResp, nil)
       if err != nil {
           return nil, fmt.Errorf("get job: %w", err)
       }

       var job Job
       if err := json.Unmarshal(apiResp.Data, &job); err != nil {
           return nil, fmt.Errorf("get job: failed to decode response: %w", err)
       }
       return &job, nil
   }
   ```

2. **Define request/response types** in the same file or `models.go`:
   ```go
   type Job struct {
       ID     string `json:"id"`
       Status string `json:"status"`
       // ...
   }
   ```

3. **Wire the sub-API into `InoueClient`** in `client.go`:
   ```go
   type InoueClient struct {
       // ... existing fields ...
       Jobs *JobsAPI
   }
   ```
   And in `NewClient()`:
   ```go
   c.Jobs = &JobsAPI{client: c}
   ```

4. **Add tests** in `client_test.go` or a new `jobs_api_test.go` file.

Key conventions:
- Use `a.client.request()` for all HTTP calls — it handles JSON marshaling, headers, trace IDs, and error parsing
- Decode `ApiResponse.Data` with `json.Unmarshal` into the specific response type
- Wrap errors with the method name prefix (e.g. `"get job: %w"`)
- Use `extraHeaders` map for endpoint-specific headers (admin token, bootstrap secret)
- For endpoints that return no data, pass `nil` as the `dest` parameter

## Development

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run with race detector
go test -race ./...

# Vet the code
go vet ./...

# Format the code
gofmt -w .
```

## Test Coverage

| Test Function | Covers |
|---|---|
| `TestNewClient` | Client construction: base URL set, Internal and Schedule sub-APIs initialized |
| `TestNewClientWithOptions` | Functional options: trailing slash trimmed, access token set, timeout applied |
| `TestSetAccessToken` | `SetAccessToken()` updates the Bearer token on an existing client |
| `TestSdkErrorFormat` | `SdkError.Error()` formats with code, status, and trace ID |
| `TestSdkErrorFormatWithoutTraceID` | `SdkError.Error()` omits trace ID when empty |
| `TestGenerateTraceID` | `generateTraceID()` produces valid UUID v4 format (36 chars, dashes at positions 8/13/18/23) |
| **Total** | **6 tests** |
