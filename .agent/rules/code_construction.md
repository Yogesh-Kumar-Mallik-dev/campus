# Rule: Modular Code Construction, Block Standards & Early Returns

## 1. Modular Block Construction

Code must be written in modular, logical blocks. Every major logical unit (handler, algorithmic step, mathematical evaluation, transaction fold) must start with a descriptive comment header:

```typescript
/**
 * BLOCK_AUTH_TOKEN_ROTATION_001
 * Purpose: Rotates refresh token family and invalidates breached sessions upon token reuse.
 * Inputs:  rawRefreshToken (string), clientIp (string), userAgent (string)
 * Outputs: TokenPair { accessToken: string, refreshToken: string }
 * Errors:  ERR_TOKEN_EXPIRED (401), ERR_TOKEN_REUSE_DETECTED (403), ERR_SESSION_REVOKED (403)
 */
```

### Benefits of Block Structure:
1. **Traceability:** Errors report the exact block ID where failure occurred (`BLOCK_AUTH_TOKEN_ROTATION_001: refresh token reused`).
2. **AI & Human Comprehension:** Code reviews and AI reasoning can isolate and reason about bounded logic blocks without whole-file cognitive overload.
3. **Refactoring Safety:** Discrete blocks can be extracted into standalone helper functions with clean input/output contracts.

---

## 2. Flat Logic with Guard Clauses & Early Returns

Deeply nested `if/else` ladders (the "Pyramid of Doom") are strictly prohibited.

### ❌ Bad: Nested Conditional Ladder
```typescript
function processOrder(order: Order, user: User): Result {
  if (user.isActive) {
    if (order.items.length > 0) {
      if (user.balance >= order.total) {
        // Business logic buried 4 levels deep
        return executePayment(order, user);
      } else {
        return error("Insufficient balance");
      }
    } else {
      return error("Empty cart");
    }
  } else {
    return error("Inactive user");
  }
}
```

### ✅ Good: Guard Clauses with Early Returns
```typescript
function processOrder(order: Order, user: User): Result {
  if (!user.isActive) {
    return error("Inactive user");
  }
  if (order.items.length === 0) {
    return error("Empty cart");
  }
  if (user.balance < order.total) {
    return error("Insufficient balance");
  }

  // Happy path remains flat and unindented
  return executePayment(order, user);
}
```

---

## 3. Strict Typing & Zero Implicit `any`

1. **Explicit Return Types:** All exported functions, methods, and API handlers must specify explicit return types.
2. **Zero `any` Hand-Waving:** Use `unknown` with type narrowing (e.g. `zod`, `io-ts`, or type guards) instead of `any`.
3. **Immutability by Default:** Prefer `readonly` arrays, `const` bindings, and immutable data structures unless in-place mutation provides measurable performance gains in hot loops.

---

## 4. Error Message Specificity

Every error created or propagated must include:
1. **The Unique Block ID:** Identifies origin in logs and telemetry.
2. **Actionable Cause:** What specific check failed (not just _"An error occurred"_).
3. **Contextual Values:** The invalid input value (e.g. `invalid page_size: 250, maximum allowed: 100`).
