# Rule: Zero-Trust Security, Multi-Tenant Isolation & Sanitization

## 1. Zero-Trust Multi-Tenant Isolation

1. **Explicit Tenancy Scoping on Every Query:**
   - Every database query for tenant-owned data must explicitly include the authenticated tenant/user ID in the `WHERE` clause:
     ```sql
     -- ✅ Good: Tenancy bounded
     SELECT * FROM documents WHERE id = $1 AND tenant_id = $2;

     -- ❌ Bad: Vulnerable to Insecure Direct Object Reference (IDOR)
     SELECT * FROM documents WHERE id = $1;
     ```
2. **Never Trust Client-Provided Tenant Identifiers:**
   - The tenant ID must be resolved cryptographically from verified session tokens or authenticated context headers, never from untrusted query parameters or request JSON bodies.

---

## 2. Zero-Trust Input Sanitization & Machine Keys

1. **Lowercase Normalized Machine Keys:**
   - All machine identifiers, field keys, tags, and category codes must be strictly lowercased and sanitized (`^[a-z0-9_\.]+$`):
     ```typescript
     export function sanitizeMachineKey(input: string): string {
       return input
         .trim()
         .toLowerCase()
         .replace(/[^a-z0-9_\.]/g, "_")
         .replace(/_{2,}/g, "_");
     }
     ```
2. **Duplicate Machine Key Rejection:**
   - Schema and entity validation layers must reject payloads containing duplicate machine keys with explicit error codes (`DUPLICATE_KEY`).

---

## 3. Secret Management & Credential Hygiene

1. **Zero Hardcoded Secrets:**
   - Never commit passwords, API keys, private certificates, or JWT secrets to version control.
   - Use environment variables (`.env` with `.env.example` templates) or secret managers (Vault, AWS Secrets Manager, 1Password CLI).
2. **Constant-Time Cryptographic Comparisons:**
   - Always use constant-time comparison algorithms (e.g. `crypto.timingSafeEqual` in Node, `subtle.ConstantTimeCompare` in Go) for password hashes, HMAC signatures, and API tokens to prevent timing side-channel attacks.
3. **Session Token Family Rotation:**
   - Refresh tokens must be single-use and part of a token family. If a used refresh token is presented again, the entire session family must be immediately revoked (token reuse breach detection).
