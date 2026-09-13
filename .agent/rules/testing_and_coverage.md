# Rule: Testing Rigor, Dependency Injection & Coverage Invariants

## 1. Hand-in-Hand Testing Invariant

1. **Co-Located Tests:** Whenever a new feature, domain module, helper, or handler is created or modified, corresponding unit/integration tests **MUST** be written or updated in the **exact same change**. No code is complete without its tests.
2. **100% Domain Logic Coverage Objective:** Strive for 100% test coverage across all business logic, state transitions, validation rules, mathematical engines, and API handlers.
3. **Branch & Edge Case Rigor:** Every test suite must test:
   - Happy paths (nominal behavior).
   - Validation failure boundaries (min/max bounds, empty inputs, type mismatches).
   - Invariant breaches (duplicate keys, unauthorized tenancy access).
   - Error states and recovery flows.

---

## 2. Dependency Injection (DI) & Decoupled Architecture

1. **Inject Mockable Interfaces:**
   - Never hardcode database connections, HTTP clients, clock timestamps (`Date.now()`), or filesystem calls inside domain business logic.
   - Inject repository interfaces, cache clients, and transport adapters via constructors or factory parameters.

### ✅ Good: Interface-Based Dependency Injection
```typescript
export interface UserRepository {
  findById(id: string): Promise<User | null>;
  save(user: User): Promise<void>;
}

export class UserService {
  constructor(private readonly userRepo: UserRepository) {}

  async activateUser(userId: string): Promise<User> {
    const user = await this.userRepo.findById(userId);
    if (!user) throw new NotFoundError(`User ${userId} not found`);
    user.status = "ACTIVE";
    await this.userRepo.save(user);
    return user;
  }
}
```

2. **Deterministic, Isolated Unit Tests:**
   - Unit tests must execute in **milliseconds** without spawning external databases, Docker containers, or live networks.
   - Use in-memory test doubles, fakes, or spies to test business logic in isolation.

---

## 3. Test Structure: Arrange-Act-Assert (AAA)

All test cases must follow the AAA pattern with clear block descriptions:

```typescript
import { describe, it } from "node:test";
import assert from "node:assert/strict";

describe("BLOCK_TEST_USER_ACTIVATION_001: User Activation Flow", () => {
  it("should successfully activate an existing user", async () => {
    // Arrange
    const fakeRepo = new InMemoryUserRepo([{ id: "u-1", status: "PENDING" }]);
    const service = new UserService(fakeRepo);

    // Act
    const activated = await service.activateUser("u-1");

    // Assert
    assert.equal(activated.status, "ACTIVE");
  });

  it("should throw NotFoundError on non-existent user ID", async () => {
    // Arrange
    const fakeRepo = new InMemoryUserRepo([]);
    const service = new UserService(fakeRepo);

    // Act & Assert
    await assert.rejects(
      () => service.activateUser("u-ghost"),
      { name: "NotFoundError" }
    );
  });
});
```

---

## 4. Test Suite Execution & CI Gates

1. **Zero Failing Tests:** All tests must pass with 100% success rate.
2. **Zero Test Flakiness:** Tests must be idempotent and execution order independent.
3. **Fast Feedback Loop:** Unit test suites should complete in under 5 seconds locally.
