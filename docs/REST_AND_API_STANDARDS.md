# Universal REST & API Engineering Standards

**Status:** Authoritative Standard (Non-Negotiable)  
**Version:** 2.0.0  
**Scope:** Universal HTTP API design, RFC 7807 problem details, idempotency, security headers, rate limiting, and observability.

---

## 1. API Design & Resource Modeling

### Resource Nouns & Hierarchies
- Always use **plural lowercase nouns** for resources (`/api/v1/students`, `/api/v1/courses`, `/api/v1/hostels`).
- Model nested sub-resources strictly to indicate parent-child ownership:
  - Top-level collection: `GET /api/v1/departments`
  - Sub-resource under parent: `GET /api/v1/departments/{deptId}/courses`
- **Flat Fallback Invariant:** If a sub-resource can exist independently or nesting exceeds 2 levels, flatten the route:
  - ✅ Good: `GET /api/v1/courses/{courseId}`
  - ❌ Bad: `GET /api/v1/universities/{uid}/faculties/{fid}/departments/{did}/courses/{cid}`

### Strict HTTP Verb Semantics

| HTTP Verb | Safety | Idempotent | Semantic Purpose & Response Invariants |
| :--- | :---: | :---: | :--- |
| `GET` | ✅ Safe | ✅ Yes | Read operations. Must never mutate server state. Zero request body. |
| `POST` | ❌ Unsafe | ❌ No | Resource creation or state processing. Returns `201 Created` with `Location` header or `200 OK`. |
| `PUT` | ❌ Unsafe | ✅ Yes | Full resource replacement. Missing payload fields reset to defaults/null. |
| `PATCH` | ❌ Unsafe | ❌ No* | Partial update. Modifies only provided fields (*idempotent when paired with `Idempotency-Key`). |
| `DELETE` | ❌ Unsafe | ✅ Yes | Resource removal. Returns `204 No Content` or `200 OK` with deleted payload. |
| `HEAD` | ✅ Safe | ✅ Yes | Same as `GET` without response body. Used for existence and cache checks. |
| `OPTIONS` | ✅ Safe | ✅ Yes | Pre-flight CORS negotiation; returns supported verbs in `Allow` header. |

### Non-CRUD Action Modeling
When operations do not fit standard CRUD verbs, model actions as sub-resource state controllers or terminal nouns:
- ❌ **Anti-Pattern:** `POST /api/v1/cancelStudentTicket?id=123`
- ✅ **Standard:** `POST /api/v1/tickets/{id}/cancellation` or `PATCH /api/v1/tickets/{id}/status` with `{"status": "CANCELLED"}`.

---

## 2. Versioning & Lifecycle Governance

- **URI Path Versioning:** Explicit `/api/v1/`, `/api/v2/` prefix.
- **Additive Evolution (Zero Major Bumps):** Adding optional request fields, new response fields, or new endpoints must never trigger `/v2/`.
- **Breaking Changes:** Renaming/removing fields, altering types, or changing status codes requires a major version path increment.
- **RFC 8594 Deprecation & Sunset Headers:**
  ```http
  Deprecation: @1759276800
  Sunset: Wed, 11 Nov 2026 00:00:00 GMT
  Link: <https://api.campus.institution.edu/docs/v2-migration>; rel="sunset"
  ```

---

## 3. Pagination, Filtering, Sorting & Projections

### Cursor-Based (Keyset) Pagination (Default for High-Velocity & Large Datasets)
- Request: `GET /api/v1/audit/logs?limit=50&cursor=eyJpZCI6MTA0MiwiY3JlYXRlZEF0IjoxNzI2MjExfQ==`
- Envelope:
  ```json
  {
    "data": [...],
    "pagination": {
      "limit": 50,
      "has_more": true,
      "next_cursor": "eyJpZCI6MTA5MiwiY3JlYXRlZEF0IjoxNzI2MzAwfQ==",
      "prev_cursor": "eyJpZCI6MTA0MiwiY3JlYXRlZEF0IjoxNzI2MjExfQ=="
    }
  }
  ```

### Offset/Limit Pagination (UI Tables with Static Page Jumping)
- Request: `GET /api/v1/courses?page=2&pageSize=10`
- Maximum page cap: `maxLimit = 100` to prevent deep offset memory exhaustion.
- **Collection Query Invariant:** Always return `200 OK` with `{"data": [], "pagination": ...}` when no items match.

### Standardized Filter & Sort Syntax
- **Exact Match:** `?status=ACTIVE&department=CSE`
- **Range & Comparative:** `?created_at[gte]=2026-01-01&created_at[lte]=2026-06-30`
- **Multi-Value IN:** `?status=ACTIVE,PENDING,PROBATION`
- **Signed Field Sorting:** `?sort=-batchYear,lastName` (minus prefix = descending).
- **Sparse Fieldsets & Embeds:** `?fields=id,rollNumber,email` and `?include=profile,hostelAllocation`.

---

## 4. Idempotency & Optimistic Concurrency Control (OCC)

### `Idempotency-Key` Header
- Mandatory on critical financial, booking, and mutation POST routes (`Idempotency-Key: <UUID>`).
- Keys are cached in Redis for 24 hours alongside the in-flight lock state and terminal response. Replayed requests return the exact cached response with zero duplicate execution.

### Optimistic Concurrency Control (ETags)
- Server returns weak ETag: `ETag: W/"v4-9f8a892b"`.
- Clients send: `If-Match: W/"v4-9f8a892b"`.
- Concurrent overwrite attempts fail immediately with `412 Precondition Failed` or `409 Conflict`.

---

## 5. RFC 7807 Standardized Problem Details Error Format

All error payloads use `Content-Type: application/problem+json`:

```json
{
  "type": "https://campus.institution.edu/errors/VALIDATION_FAILED",
  "title": "Unprocessable Content",
  "status": 422,
  "detail": "The request payload failed semantic validation rules.",
  "instance": "/api/v1/hostel/gate-passes",
  "code": "VALIDATION_FAILED",
  "invalid_params": [
    {
      "name": "returnExpected",
      "reason": "Must be a timestamp later than departExpected"
    }
  ],
  "trace_id": "req_84a91c82-55db-4927-b50a-e24c5520e557"
}
```

### Canonical HTTP Status Code Directory

| Status Code | Standard | Usage Invariant |
| :--- | :--- | :--- |
| `200 OK` | RFC 9110 | Successful GET, PUT, PATCH, or collection query with body |
| `201 Created` | RFC 9110 | Successful POST creating persistent resource; includes `Location` header |
| `202 Accepted` | RFC 9110 | Asynchronous job queued for background execution |
| `204 No Content` | RFC 9110 | Successful execution with empty body (typical for DELETE) |
| `304 Not Modified` | RFC 9110 | Conditional ETag/If-None-Match cache hit (0 byte body) |
| `400 Bad Request` | RFC 9110 | Malformed JSON syntax or unparseable framing |
| `401 Unauthorized` | RFC 9110 | Missing, invalid, or expired authentication token |
| `403 Forbidden` | RFC 9110 | Authenticated user lacks permission or crosses tenant boundary |
| `404 Not Found` | RFC 9110 | Specific resource URI does not exist |
| `409 Conflict` | RFC 9110 | Duplicate unique key, lock violation, or room double-booking |
| `412 Precondition Failed`| RFC 9110 | If-Match validation failure during optimistic locking |
| `413 Payload Too Large` | RFC 9110 | Request body exceeds gateway byte ceiling (250KB JSON, 10MB upload) |
| `415 Unsupported Media` | RFC 9110 | Request lacks `Content-Type: application/json` on mutation |
| `422 Unprocessable` | RFC 9110 | Semantic validation failure on syntactically valid JSON |
| `429 Too Many Requests` | RFC 6585 | Rate limit exceeded; must include `Retry-After` header |
| `500 Internal Error` | RFC 9110 | Unhandled server runtime exception (logs trace ID) |
| `502 Bad Gateway` | RFC 9110 | Downstream service returned invalid response |
| `503 Service Unavailable` | RFC 9110 | Circuit breaker open or maintenance mode |
| `504 Gateway Timeout` | RFC 9110 | Database query or downstream call breached context timeout |

---

## 6. Security, Zero Trust & Hardened Headers

### Hardened Response Headers on Every HTTP Response
```http
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
Referrer-Policy: strict-origin-when-cross-origin
Content-Security-Policy: default-src 'none'; frame-ancestors 'none'
Cross-Origin-Opener-Policy: same-origin
Cross-Origin-Resource-Policy: same-origin
Strict-Transport-Security: max-age=63072000; includeSubDomains; preload
```
*Note: `Server`, `X-Powered-By`, and runtime banners are stripped unconditionally.*

### Zero-Trust Multi-Tenancy & Sanitization
- Every database query scopes tenant boundary explicitly (`WHERE tenant_id = $1`).
- Gateway sets strict JSON body limit (**250 KB**) to prevent memory starvation attacks.
- Non-whitelisted payload fields are stripped to block mass-assignment privilege escalation.

---

## 7. Rate Limiting & Resilience Architecture

### IETF Draft Rate Limit Headers
```http
RateLimit-Limit: 100
RateLimit-Remaining: 24
RateLimit-Reset: 1726211400
Retry-After: 60
```

### Context Deadlines & Cancellation Propagation
- Every request propagates `context.WithTimeout(ctx, 5*time.Second)`.
- If client disconnects or times out, Go cancels the context and PostgreSQL terminates the running SQL query immediately.

### Circuit Breakers & Full Jitter Exponential Backoff
- External integrations (payment gateways, SMS dispatch) wrap with circuit breakers (Closed $\to$ Open $\to$ Half-Open).
- Retries use exponential backoff with full jitter to avoid thundering herds:
  $$T_{\text{wait}} = \min(T_{\text{max}}, T_{\text{base}} \times 2^{\text{attempt}}) \pm \text{jitter}$$

---

## 8. HTTP Performance & Caching

- **ETag Validation:** Content hash / update timestamp ETags return `304 Not Modified` on cache hits.
- **Cache-Control Directives:**
  - User private data: `Cache-Control: private, no-cache, no-store, must-revalidate`
  - Public catalogs: `Cache-Control: public, max-age=300, stale-while-revalidate=60`
  - Static immutable assets: `Cache-Control: public, max-age=31536000, immutable`
- **Connection Pooling:** Bounded `pgxpool` connections:
  $$\text{Max Connections} = (\text{CPU Cores} \times 2) + \text{Disk Spindles}$$
- **$O(1)$ Batch Fetching:** Keyset pagination and `WHERE id IN (...)` eliminate $N+1$ query cascades.

---

## 9. Observability, Telemetry & Health Probes

### Request Tracing & Correlation IDs
- Edge gateway assigns or validates `X-Request-ID` (and W3C `traceparent`).
- Propagated to all logs, SQL query comments (`/* trace_id=xyz */`), and returned in response headers.

### Structured JSON Logging (Zero Credentials)
```json
{
  "timestamp": "2026-09-13T05:49:41.102Z",
  "level": "WARN",
  "trace_id": "req_d4b2e8a1-7c9f-4b08-8e6d-6211f42d591b",
  "method": "PATCH",
  "route": "/api/v1/hostel/rooms/{id}",
  "status": 409,
  "latency_ms": 18.4,
  "user_id": "usr_941a8c82",
  "client_ip": "203.0.113.195",
  "error_code": "ROOM_CAPACITY_EXCEEDED"
}
```

### Health Probes Architecture
- **Liveness Probe:** `GET /healthz/live` $\to$ `200 OK` (memory-only check validating process responsiveness).
- **Readiness Probe:** `GET /healthz/ready` $\to$ `200 OK` or `503 Service Unavailable` (validates active PostgreSQL socket and Redis connections).
