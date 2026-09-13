# Current Agent Context & State

**Repository:** Campus Management System  
**Last Updated:** 2026-09-13  
**Status:** Architecture Specification & Project Initialization  

---

## 1. System Overview & Scope

The Campus Management System is an enterprise-grade platform unifying 17 core campus subsystems under strict modular monolith architecture and universal engineering standards:

1. **Public Web Portal System** (`portal`): Admissions showcase, course catalog, public events, campus news, public prospect onboarding.
2. **Central Auth & Permissions System (IAM)** (`auth`): Multi-tenant RBAC, session rotation, TOTP MFA, constant-time comparisons.
3. **Mess Management System** (`mess`): Meal scheduling, menu publishing, dietary preferences, QR/token validation, rebate tracking, mess billing.
4. **Hostel Management System** (`hostel`): Room/bed allocation, warden approvals, curfew logs, gate passes, maintenance requests.
5. **Notice & Announcement System** (`notices`): Targeted role/department broadcasts, rich text, priority pinning, expiration lifecycles, push delivery.
6. **Application & Query Resolution System (Helpdesk)** (`helpdesk`): Ticket lifecycle, SLA timer tracking, department triage, grievance escalation.
7. **E-Library System** (`library`): Catalog management (ISBN/barcode), physical issue/return/renew, digital media viewer, automated fine calculation.
8. **The Hub System (Root Super-App Shell)** (`hub`): Universal root host container providing role-specialized hubs: Student Hub, Teacher Hub, Admin Hub, Librarian Hub, Parent Hub, Warden Hub.
9. **Anonymity & Whistleblower System** (`whistleblower`): Zero-knowledge encrypted incident reporting (anti-ragging, safety, compliance), anonymous 2-way messaging, anti-retaliation logs.
10. **Attendance Management System** (`attendance`): Multi-mode attendance (faculty manual mark, biometric terminal sync, RFID/QR, geofenced mobile), shortage calculation (<75% threshold alerts), medical leaves.
11. **Study Hub System** (`studyhub`): Course syllabus, lecture notes repository, assignments submission & grading, past question papers, peer discussion threads.
12. **Event Organisation System** (`events`): Campus calendar, venue/hall bookings & approvals, RSVP/ticketing, volunteer staffing, digital certificate generation.
13. **Progress Tracker & Resolution System (Mentorship/Remediation)** (`mentorship`): Faculty mentor allocation, academic risk flags, remedial action plans, 1-on-1 counseling logs, milestone tracking.
14. **SOS & Emergency Response System** (`sos`): 1-tap emergency trigger, GPS location broadcast, campus security alerts, emergency broadcast siren/SMS/push, dispatch tracking.
15. **Student & Staff Registration System (Onboarding)** (`onboarding`): Applicant KYC, document upload & verification, roll number/employee ID generation, department/batch assignment.
16. **Central Payment & Billing System** (`billing`): Unified invoice engine (tuition, hostel, mess, library fines), payment gateway integration, receipt generation, split installment plans.
17. **Central Audit & Compliance System** (`audit`): Immutable tamper-evident audit logs, accreditation compliance reporting (NAAC, NIRF, ABET), suspicious activity triggers.

---

## 2. Definitive Execution Ranking (Rank 1 to 17)

1. **Rank 1: Central Auth & Permissions System (IAM)** (`auth`) — Identity, token families, RBAC, TOTP MFA.
2. **Rank 2: Central Audit & Compliance System** (`audit`) — Tamper-evident ledger, event subscriber.
3. **Rank 3: Student & Staff Registration System (Onboarding)** (`onboarding`) — Verified roll numbers, staff profiles, cohorts.
4. **Rank 4: Attendance Management System** (`attendance`) — Multi-mode roll call, shortage calculations (<75%).
5. **Rank 5: Central Payment & Billing System** (`billing`) — Universal invoicing, double-entry ledger, gateways.
6. **Rank 6: Notice & Announcement System** (`notices`) — Targeted broadcast alerts, priority pinning, push alerts.
7. **Rank 7: Hostel Management System** (`hostel`) — Room/bed inventory, warden approvals, curfew, gate passes.
8. **Rank 8: Mess Management System** (`mess`) — Dining menus, QR validation, leave rebates linked to hostel leaves.
9. **Rank 9: E-Library System** (`library`) — Book cataloging, issue/return state machine, fine generator.
10. **Rank 10: Study Hub System** (`studyhub`) — Course syllabus, notes repository, assignment submission & grading.
11. **Rank 11: Progress Tracker & Resolution System (Mentorship)** (`mentorship`) — At-risk detection, counseling logs.
12. **Rank 12: Event Organisation System** (`events`) — Institutional calendar, venue conflict checks, RSVP ticketing.
13. **Rank 13: Application & Query Resolution System (Helpdesk)** (`helpdesk`) — Support tickets, SLA timers, routing.
14. **Rank 14: SOS & Emergency Response System** (`sos`) — 1-tap mobile SOS, GPS dispatch, emergency sirens.
15. **Rank 15: Anonymity & Whistleblower System** (`whistleblower`) — Zero-knowledge reporting, anti-ragging channel.
16. **Rank 16: Public Web Portal System** (`portal`) — External showcase, programs catalog, admissions lead capture.
17. **Rank 17: The Hub System (Root Super-App Shell)** (`hub`) — Universal root host shell uniting Student Hub, Teacher Hub, Admin Hub, Librarian Hub, Parent Hub, Warden Hub.

---

## 3. Technology Stack & Sizing Targets

- **Frontend Web & Desktop:** Svelte 5 (Runes) + SvelteKit 2 + shadcn-svelte + Tauri 2.
- **Frontend Mobile:** React Native + Expo Router + NativeWind + React Native Reusables (`@rn-primitives`).
- **Backend Core:** Go 1.23+ Modular Monolith with `go-chi` and interface-driven Dependency Injection.
- **Detached Data Layer:** PostgreSQL 18 + Prisma 6 schema & TypeScript tooling (`database/`).
- **Monorepo Toolchain:** Turborepo + pnpm workspace + `./script.sh` / `.\script.ps1`.
- **Target Scale:** 30,000 total users, 5,000 peak concurrent active sessions (baseline: ~1,200 users).

---

## 4. Active Focus & Immediate Next Steps

- **Completed Subsystem:** **Rank 1: Central Auth & Permissions System (IAM)** (`auth`) — 100% Complete & Verified:
  - Step 1: PostgreSQL 18 schema defined in `database/schema.prisma`.
  - Step 2: Go domain core, Argon2id hasher, JWT signer, RFC 6238 TOTP engine, and mock test suite in `backend/auth/`.
  - Step 3: RFC 7807 Problem Details HTTP engine (`api/problem/`) and Chi REST handlers & tests (`api/http/auth/`).
  - Step 4: Svelte 5 / shadcn-svelte web UI and React Native Reusables mobile UI components for IAM.
  - Step 5: Synced domain documentation created at `docs/auth.md`. All 19 unit and integration tests passing with 0 failures.
- **Next Subsystem:** **Rank 2: Central Audit & Compliance System** (`audit`).
- **Immediate Next Step:** Formulate domain model, immutable audit ledger schema, and async event subscriber contracts for `audit`.
