# Product Requirements Document (PRD)

## Inoue AI SaaS Go SDK

### Overview

The Inoue AI SaaS Go SDK is a lightweight, zero-dependency Go client library for the Inoue AI SaaS Backend API. It provides typed access to worker management (registration, JWT token minting, status reporting) and content scheduling endpoints. The SDK is the Go counterpart to the existing Python SDK (`inoue_ai_saas_sdk`), purpose-built for Go-based workers and services in the Inoue AI platform. It handles JSON envelope parsing, automatic trace ID propagation, structured error extraction, and Bearer token lifecycle using only the Go standard library.

### Problem Statement

Go-based workers in the Inoue AI ecosystem need to communicate with the Backend API for worker bootstrap (registration + JWT minting), heartbeat status reporting, and content scheduling. Without a dedicated SDK, each Go worker would need to independently implement HTTP request construction, the Backend's JSON envelope parsing (`{ok, data, error, meta}`), trace ID generation, Bearer token management, error extraction from nested error payloads, and idempotent registration logic. This duplicated effort is error-prone, leads to inconsistent error handling across workers, and makes it difficult to add new API endpoints without updating every consumer. The Python SDK already solves this for Python workers — Go workers needed an equivalent.

### Goals

1. **Type-safe Backend API access** — Provide Go structs for all request parameters, response types, and error payloads so that Go workers interact with the Backend through compile-time-checked interfaces rather than raw HTTP and `map[string]interface{}`.
2. **Consistent envelope handling** — Parse the Backend's standard `{ok, data, error, meta}` JSON envelope in one place, with lazy `json.RawMessage` decoding of the `data` field, so sub-API methods only deal with their specific response types.
3. **Automatic trace ID propagation** — Generate a UUID v4 trace ID for every request and attach it as the `X-Trace-Id` header. Embed the same trace ID in any returned `SdkError` for end-to-end correlation across SDK caller, Backend API, and downstream services.
4. **Worker bootstrap in one call** — Expose `RegisterWorker()` (idempotent, 409 = success) and `WorkerToken()` (JWT minting) so Go workers can bootstrap themselves with two method calls and no conditional error handling for already-registered state.
5. **Zero external dependencies** — Use only the Go standard library to eliminate dependency management overhead, minimize binary size, and ensure compatibility across all Go environments and build constraints.
6. **Python SDK parity for covered endpoints** — Mirror the Python SDK's `client.internal.*` and `client.schedule.*` patterns so developers familiar with the Python SDK can use the Go SDK with the same mental model.

### Target Users

- Go-based workers in the Inoue AI platform (event workers, scheduled task runners, content pipeline services) that need to register with the Backend, obtain JWTs, report status, and create schedule entries
- Platform developers extending the SDK with new Backend API endpoints
- Any Go service in the Inoue AI ecosystem that communicates with the Backend API

### System Architecture

```
                    ┌──────────────────────────────────────────┐
                    │        Inoue AI SaaS Backend API          │
                    │                                           │
                    │  /internal/workers/register    (POST)     │
                    │  /internal/workers/token       (POST)     │
                    │  /internal/workers/me/status   (POST)     │
                    │  /v1/schedule/                 (POST)     │
                    │                                           │
                    │  Response envelope:                       │
                    │  {ok, data, error, meta}                  │
                    └──────────────────┬───────────────────────┘
                                       │
                                       │  HTTPS + JSON
                                       │  Authorization: Bearer {JWT}
                                       │  X-Trace-Id: {UUID v4}
                                       │
                    ┌──────────────────▼───────────────────────┐
                    │        Inoue AI SaaS Go SDK               │
                    │                                           │
                    │  InoueClient                              │
                    │  ├── request()       (HTTP engine)        │
                    │  ├── Internal        (InternalAPI)        │
                    │  │   ├── RegisterWorker()                 │
                    │  │   ├── WorkerToken()                    │
                    │  │   └── SetWorkerStatus()                │
                    │  └── Schedule        (ScheduleAPI)        │
                    │      └── Create()                         │
                    │                                           │
                    │  Zero external dependencies               │
                    │  (net/http, encoding/json, crypto/rand)   │
                    └──────────────────┬───────────────────────┘
                                       │
                                       │  Imported as Go module
                                       │
           ┌───────────────────────────┼───────────────────────────┐
           │                           │                           │
┌──────────▼──────────┐  ┌────────────▼────────────┐  ┌──────────▼──────────┐
│  Go Event Worker     │  │  Go Scheduled Runner    │  │  Future Go Workers  │
│                      │  │                         │  │                     │
│  1. RegisterWorker() │  │  1. RegisterWorker()    │  │  Same bootstrap +   │
│  2. WorkerToken()    │  │  2. WorkerToken()       │  │  API pattern        │
│  3. SetAccessToken() │  │  3. SetAccessToken()    │  │                     │
│  4. SetWorkerStatus()│  │  4. Schedule.Create()   │  │                     │
│  5. Business logic   │  │  5. Business logic      │  │                     │
└─────────────────────┘  └─────────────────────────┘  └─────────────────────┘
```

### Core Features

#### 1. Client Construction with Functional Options

The `InoueClient` is created via `NewClient(baseURL, opts...)` with optional functional options:

- `WithAccessToken(token)` — Sets the initial Bearer token so requests are authenticated from the first call
- `WithTimeout(duration)` — Overrides the default 30-second HTTP timeout
- `WithHTTPClient(client)` — Replaces the entire `http.Client` for custom transport, TLS configuration, or proxy support

The base URL has its trailing slash trimmed to prevent double-slash issues when concatenating paths. Sub-API structs (`InternalAPI`, `ScheduleAPI`) are initialized after options are applied and hold a back-reference to the parent client.

#### 2. Authentication & Token Lifecycle

The SDK supports two authentication flows:

- **Direct token** — Pass a known token via `WithAccessToken()` at construction or `SetAccessToken()` at any time
- **Worker bootstrap** — Call `RegisterWorker()` with an admin token to ensure the worker row exists, then call `WorkerToken()` with the bootstrap secret to mint a JWT. Set the returned JWT with `SetAccessToken()` for all subsequent calls.

The Bearer token is injected as the `Authorization` header on every request where the token is non-empty. Endpoints that use alternative auth headers (`X-Admin-Token`, `X-Worker-Bootstrap`) bypass the Bearer token and pass their auth via the extra headers mechanism.

#### 3. Internal Worker Management API

Three endpoints for the worker lifecycle:

- **RegisterWorker** — `POST /internal/workers/register` with `X-Admin-Token`. Sends `{id, name}`. Idempotent: HTTP 409 (worker already registered) is silently treated as success, returning `nil`.
- **WorkerToken** — `POST /internal/workers/token` with `X-Worker-Bootstrap`. Sends `{worker_id}`. Unwraps the `worker_token` field from the `ApiResponse.Data` envelope and returns the JWT string.
- **SetWorkerStatus** — `POST /internal/workers/me/status` with Bearer token. Sends `{status}`. Used for heartbeat pings (`"active"`) and graceful shutdown (`"disabled"`).

#### 4. Schedule API

Content scheduling via `POST /v1/schedule/`:

- **Create** — Accepts `ScheduleEntryCreateParams` with required fields (`ModelID`, `PlatformID`, `ContentTypeID`, `ScheduledFor`) and optional fields (`OrgID`, `AssetIDs`, `Notes`). Returns a `*ScheduleEntry` with the server-assigned `ID`, `Status`, and `CreatedAt`.
- Optional fields use pointer types (`*string`) with `omitempty` JSON tags so they are excluded from the request body when nil.

#### 5. Structured Error Handling

Two error categories with distinct handling:

- **API errors** (`*SdkError`) — Returned when HTTP status >= 400. Contains `Code` (app-level error code), `Message` (human-readable), `Status` (HTTP code), `TraceID` (request trace), and `Details` (optional structured data). The error is extracted from the Backend's JSON envelope when possible; falls back to raw body with `code: "unknown_error"` when the body is not valid JSON.
- **Transport errors** — Returned as standard `error` values wrapped with `fmt.Errorf` and the `%w` verb for `errors.Is()`/`errors.As()` compatibility. Prefixed with `"inoue-sdk:"` for grep-ability.

`SdkError` implements the `error` interface with a formatted string that includes the message, code, status, and trace ID.

#### 6. Trace ID Propagation

Every request generates a UUID v4 trace ID using `crypto/rand`:

- 16 random bytes with version 4 and variant bits set per RFC 4122
- Formatted as `xxxxxxxx-xxxx-4xxx-Nxxx-xxxxxxxxxxxx`
- Attached as the `X-Trace-Id` header on the outgoing request
- Embedded in any `SdkError` returned from that request
- Enables correlation between SDK caller logs, Backend API logs, and downstream service logs

### Non-Functional Requirements

- **Go 1.22+** minimum version (uses generic `Page[T]` type)
- **Zero external dependencies** — Only Go standard library packages: `net/http`, `encoding/json`, `crypto/rand`, `context`, `fmt`, `io`, `bytes`, `strings`, `time`
- **Context propagation** — All API methods accept `context.Context` as first parameter for cancellation, timeout, and deadline support
- **Thread safety** — The `http.Client` is safe for concurrent use; token updates via `SetAccessToken()` are not synchronized (callers should set the token before concurrent use or use external synchronization)
- **Testability** — `WithHTTPClient()` allows injecting a custom `http.Client` with a test transport for unit testing without network calls
- **Idiomatic Go** — Functional options pattern, `error` interface implementation, `context.Context` first parameter, pointer receivers, `json` struct tags
- **Module versioning** — Tagged as `v0.1.x` following Go module semantic versioning

### Configuration

The SDK takes no environment variables. All configuration is explicit:

| Parameter | How It Is Passed | Default | Description |
|---|---|---|---|
| Base URL | `NewClient(baseURL)` | Required, no default | Backend API base URL (e.g. `https://api.inoue.ai`) |
| Access token | `WithAccessToken(token)` or `SetAccessToken(token)` | Empty string (unauthenticated) | Bearer token for `Authorization` header |
| HTTP timeout | `WithTimeout(duration)` | 30 seconds | Timeout for all HTTP requests |
| HTTP client | `WithHTTPClient(client)` | `&http.Client{Timeout: 30s}` | Complete control over transport, TLS, proxy |
| Admin token | `RegisterWorker(..., adminToken)` parameter | N/A (per-call) | `X-Admin-Token` header for worker registration |
| Bootstrap secret | `WorkerToken(..., bootstrapSecret)` parameter | N/A (per-call) | `X-Worker-Bootstrap` header for JWT minting |

### Success Metrics

- All 6 unit tests pass (`go test ./...`)
- Go vet reports no issues (`go vet ./...`)
- Client construction trims trailing slashes and initializes both sub-APIs
- Functional options correctly apply access token and timeout overrides
- `SdkError.Error()` produces the expected formatted string with and without trace ID
- `generateTraceID()` produces valid UUID v4 format (36 characters, dashes at positions 8/13/18/23)
- Worker bootstrap flow (register + token + set token) works end-to-end against the Backend API
- Schedule entry creation returns a fully populated `ScheduleEntry` with server-assigned fields
- 409 Conflict on `RegisterWorker()` returns `nil` (idempotent)
- HTTP >= 400 responses produce `*SdkError` with code, message, status, and trace ID extracted from the JSON envelope
- Non-JSON error responses produce `*SdkError` with `code: "unknown_error"` and raw body as message

### Future Considerations

- Additional sub-APIs mirroring the Python SDK's coverage (Jobs, Assets, Models, Auth, etc.)
- Retry middleware with configurable backoff for transient errors (5xx, timeouts)
- Token auto-refresh on 401 responses (re-mint JWT and retry the request)
- Connection pooling tuning via `http.Transport` defaults
- OpenTelemetry trace context propagation alongside the custom `X-Trace-Id` header
- Rate limiting / circuit breaker for high-throughput workers
- `Page[T]` iterator helpers for paginated endpoint consumption
- Code generation from Backend OpenAPI spec to keep the SDK in sync automatically
- Webhook signature verification utilities
- Context-based structured logging integration
