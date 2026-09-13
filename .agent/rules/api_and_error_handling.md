# Rule: REST API Standards, RFC 7807 Problem Details & Query Semantics

## 1. RFC 7807 Problem Details Specification

All HTTP error responses must adhere strictly to the **RFC 7807 Problem Details for HTTP APIs** specification with `Content-Type: application/problem+json`:

```json
{
  "type": "https://api.example.com/errors/validation-failed",
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
  ]
}
```

### Standard HTTP Status Codes

| Code | Status | Usage Invariant |
| :--- | :--- | :--- |
| `200` | OK | Successful GET, PUT, PATCH, or sync collection query |
| `201` | Created | Successful POST creating a new persistent resource |
| `204` | No Content | Successful DELETE with empty body |
| `400` | Bad Request | Malformed JSON or unparseable query syntax |
| `401` | Unauthorized | Missing, invalid, or expired authentication token |
| `403` | Forbidden | Authenticated user lacks permission or crosses tenant isolation boundaries |
| `404` | Not Found | Specific resource ID does not exist |
| `409` | Conflict | State conflict, duplicate key, or lock violation |
| `422` | Unprocessable Content | Semantic validation failure on syntactically valid payload |
| `429` | Too Many Requests | Rate limit exceeded (must include `Retry-After` header) |
| `500` | Internal Server Error | Unhandled server exception (must log internal trace and return opaque ID) |

---

## 2. Collection Query Empty Semantics (200 OK vs 404)

1. **Collection Queries Return 200 OK with `[]`:**
   - When a client queries a collection (e.g. `GET /projects/proj-123/items`), if the parent resource exists and the user has access, but no items currently exist, the endpoint **MUST return `200 OK` with an empty array (`{"data": [], "pagination": ...}`)**.
   - Returning `404 Not Found` for empty collections is an anti-pattern that breaks frontend initialization and creates cascading UI errors.
2. **Individual Lookups Return 404:**
   - Single-item lookups (e.g. `GET /projects/proj-123/items/item-999`) MUST return `404 Not Found` if that specific item ID does not exist.
3. **Tenancy Violations Return 403 Forbidden:**
   - If a user queries a parent resource belonging to another tenant, the API must return `403 Forbidden` (or `404 Not Found` if multi-tenant existence masking is configured).

---

## 3. Standardized Pagination Envelope

All collection endpoints must provide pagination envelopes:

```json
{
  "data": [
    { "id": "item-1", "name": "Alpha" }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 10,
    "totalItems": 1,
    "totalPages": 1,
    "hasNextPage": false,
    "hasPrevPage": false
  }
}
```
