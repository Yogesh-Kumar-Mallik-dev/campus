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

- **Completed Subsystems & Infrastructure:**
  - **Rank 1: Central Auth & Permissions System (IAM)** (`auth`) & **Universal API Gateway & Transport Layer** — 100% Complete & Verified.
  - **Rank 2: Central Audit & Compliance System** (`audit`) — 100% Complete & Verified.
  - **Rank 3: Student & Staff Registration System (Onboarding)** (`onboarding`) — 100% Complete & Verified (KYC state machine, document attachment & verification, deterministic roll number & employee ID generation, cohort provisioning, and admissions queues on Web and Mobile).
  - **3 Frontends Feature Parity & UI Modernization (Web, Mobile, Desktop):**
    - **Web (SvelteKit 2 + Svelte 5 + Vite 6 + shadcn-svelte):** Upgraded UI architecture with complete shadcn-svelte component library, unified OKLCH theme and dark mode toggling (`+layout.css`), and modernized views: The Hub Home (`/`), Auth Login with TOTP (`/auth/login`), Audit Ledger Explorer (`/audit`), Institutional Accreditation Reports (`/audit/compliance`), Public KYC Application (`/onboarding`), and Admissions & Staff Verification Desk (`/admin/onboarding`).
    - **Mobile (React Native + Expo Router):** Configured `app.json`, root layouts, auth screens (`app/(auth)/login.tsx`), The Hub Dashboard (`app/(app)/index.tsx`), Mobile Audit Trail (`app/(app)/audit.tsx`), Accreditation Reports (`app/(app)/compliance.tsx`), and Mobile KYC Registration & Tracking (`app/(app)/onboarding.tsx`).
    - **Desktop (Tauri 2):** Configured `src-tauri/tauri.conf.json`, `Cargo.toml`, and native Rust entrypoints (`lib.rs`, `main.rs`) wrapping the SvelteKit frontend.
  - **Universal Lifecycle Scripts & Verification:**
    - Updated `./script.sh build`, `dev`, `check`, `test`, `deps` with Vite 6 SvelteKit compilation, Go backend server build, `go vet`, and `svelte-check` (0 errors, 0 warnings).
    - 64 unit and integration tests passing with 100% success across 9 packages.
- **Next Subsystem:** **Rank 4: Attendance Management System** (`attendance`).
- **Immediate Next Step:** Formulate domain model for multi-mode attendance (faculty manual roll call, biometric terminal sync, RFID/QR, geofenced mobile), shortage calculation engine (<75% threshold alerts), medical leave approvals, and daily timetable period sessions.




