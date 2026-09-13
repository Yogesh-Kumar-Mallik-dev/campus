# Rule: REST API Standards, RFC 7807 Problem Details & Query Semantics

## 1. Resource Modeling & URI Architecture

1. **Plural Lowercase Resource Nouns:**
   - Always use plural lowercase nouns for resource collections (`/api/v1/students`, `/api/v1/courses`, `/api/v1/hostels`).
   - Nested sub-resources model direct ownership (`/api/v1/departments/{deptId}/courses`).
   - Flatten any sub-resource hierarchy that can exist independently or when nesting exceeds 2 levels (`/api/v1/courses/{courseId}`).
2. **Strict HTTP Verb Semantics:**
   - `GET`: Safe, idempotent read operations. Must never mutate server state. Zero request body.
   - `POST`: Unsafe, non-idempotent creation or processing. Returns `201 Created` with `Location` header or `200 OK`.
   - `PUT`: Unsafe, idempotent full replacement. Omitted payload fields reset to defaults or null.
   - `PATCH`: Unsafe partial updates. Modifies only explicitly provided fields.
   - `DELETE`: Unsafe, idempotent removal. Returns `204 No Content` or `200 OK`.
3. **Non-CRUD Actions:**
   - Model non-CRUD actions as sub-resource state transitions or terminal nouns (e.g., `POST /api/v1/tickets/{id}/cancellation` or `PATCH /api/v1/tickets/{id}/status` with `{"status": "CANCELLED"}`).

---

## 2. Collection Query Empty Semantics (200 OK vs 404)

```mermaid
flowchart TD
    Req[GET /api/v1/projects/{id}/items] --> ValidateTenancy{User has access to project?}
    ValidateTenancy -- No --> Return403[403 Forbidden: Tenancy Breach]
    ValidateTenancy -- Yes --> QueryItems{Items exist in project?}
    QueryItems -- No --> Return200Empty[200 OK: data: [], pagination: totalItems: 0]
    QueryItems -- Yes --> Return200Data[200 OK: data: [...], pagination: totalItems: N]
```

1. **Collection Queries Return 200 OK with `[]`:**
   - When a client queries a collection (e.g. `GET /projects/proj-123/items`), if the parent resource exists and the user has access, but no items currently match, the endpoint **MUST return `200 OK` with an empty array (`{"data": [], "pagination": ...}`)**.
   - Returning `404 Not Found` for empty collections is a breaking anti-pattern.
2. **Individual Lookups Return 404:**
   - Single-item lookups (e.g. `GET /projects/proj-123/items/item-999`) MUST return `404 Not Found` if that specific item ID does not exist.
3. **Tenancy Violations Return 403 Forbidden:**
   - If a user queries across tenant boundaries, the API must return `403 Forbidden` (or masked `404 Not Found` if existence privacy is configured).

---

## 3. RFC 7807 Problem Details Specification

All HTTP error responses must adhere strictly to the **RFC 7807 Problem Details for HTTP APIs** specification with `Content-Type: application/problem+json`:

```json
{
  "type": "https://campus.institution.edu/errors/VALIDATION_FAILED",
  "title": "Unprocessable Content",
  "status": 422,
  "detail": "Field 'email' must be a valid RFC 5322 email address.",
  "instance": "/api/v1/users",
  "code": "VALIDATION_FAILED",
  "invalid_params": [
    {
      "name": "email",
      "reason": "Invalid email syntax"
    }
  ],
  "trace_id": "req_84a91c82-55db-4927-b50a-e24c5520e557"
}
```

### Standard HTTP Status Codes

| Code | Status | Usage Invariant |
| :--- | :--- | :--- |
| `200` | OK | Successful GET, PUT, PATCH, or sync collection query |
| `201` | Created | Successful POST creating a new persistent resource (with `Location` header) |
| `202` | Accepted | Asynchronous batch/job scheduled for background execution |
| `204` | No Content | Successful DELETE with empty body |
| `304` | Not Modified | Conditional ETag/If-None-Match cache hit (0 byte body) |
| `400` | Bad Request | Malformed JSON or unparseable query syntax |
| `401` | Unauthorized | Missing, invalid, or expired authentication token |
| `403` | Forbidden | Authenticated user lacks permission or crosses tenant boundary |
| `404` | Not Found | Specific resource ID does not exist |
| `409` | Conflict | State conflict, duplicate unique key, or concurrent mutation conflict |
| `412` | Precondition Failed | `If-Match` validation failure during Optimistic Concurrency Control |
| `413` | Payload Too Large | Request body exceeds gateway byte ceiling (e.g. 250 KB JSON limit) |
| `415` | Unsupported Media Type | Request lacks `Content-Type: application/json` on body mutations |
| `422` | Unprocessable Content | Semantic validation failure on syntactically valid JSON payload |
| `429` | Too Many Requests | Rate limit exceeded (must include `Retry-After` header) |
| `500` | Internal Server Error | Unhandled server exception (must log internal trace and return opaque ID) |
| `502` | Bad Gateway | Downstream dependency failure |
| `503` | Service Unavailable | Circuit breaker open or service in maintenance |
| `504` | Gateway Timeout | Database or downstream call breached context deadline |

---

## 4. Idempotency & Optimistic Concurrency Control (OCC)

1. **`Idempotency-Key` Header:**
   - Mandatory on financial, billing, and critical state-mutating POST endpoints.
   - Cached in distributed storage (e.g. Redis) for 24 hours. Replayed requests return cached terminal response with zero duplicate execution.
2. **ETag & `If-Match` Optimistic Locking:**
   - Mutation endpoints emit weak/strong `ETag` headers.
   - Concurrent updates must supply `If-Match: <ETag>`; outdated writes are rejected with `412 Precondition Failed`.

---

## 5. Security Headers, Rate Limiting & Health Probes

1. **Hardened Response Headers:**
   - Every response must include `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, `Referrer-Policy: strict-origin-when-cross-origin`, `Content-Security-Policy: default-src 'none'`, `Strict-Transport-Security: max-age=63072000`.
   - Information leak headers (`Server`, `X-Powered-By`) must be stripped unconditionally.
2. **IETF Rate Limit Headers:**
   - Rate-limited endpoints emit `RateLimit-Limit`, `RateLimit-Remaining`, `RateLimit-Reset`, and `Retry-After`.
3. **Health Probes:**
   - `GET /healthz/live` returns `200 OK` (memory-only liveness check).
   - `GET /healthz/ready` returns `200 OK` or `503 Service Unavailable` (verifies live DB/cache connectivity).
