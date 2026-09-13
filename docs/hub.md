# Subsystem 17: The Hub Root Super-App (`hub`)

## 1. Overview & Business Context

The **Hub Root Super-App** is the central unifying cockpit and topological root of the Campus Management System (CMS). It aggregates live KPIs, metrics, and actionable alerts from all 16 modular subsystems (`auth`, `audit`, `onboarding`, `attendance`, `billing`, `notices`, `hostel`, `mess`, `library`, `studyhub`, `mentorship`, `events`, `helpdesk`, `sos`, `whistleblower`, `portal`), providing personalized cockpits for every institutional persona (`STUDENT`, `FACULTY`, `WARDEN`, `LIBRARIAN`, `ADMIN`, `SUPER_ADMIN`).

Key architectural pillars:
- **Unified Multi-Persona Cockpit:** Tailored dashboards surfacing domain-specific KPIs (e.g., student exam eligibility & balance, faculty grading workloads, warden curfew passes, librarian circulation metrics, admin financial inflows).
- **1-Tap Quick Action Shortcuts:** Context-aware shortcuts accelerating day-to-day interactions (e.g., geofenced attendance scan, 1-tap SOS distress trigger, meal QR code generation, fee payment, anonymous whistleblowing).
- **Dynamic Widget Layout Customizer:** End-user customization enabling modular cards to be reordered, resized (`SMALL`, `MEDIUM`, `LARGE`, `FULL_WIDTH`), and toggled based on preference.
- **Topological Telemetry Radar:** Real-time visibility into all 16 modular engine states, latency, and compliance status.
- **Audit & Activity Tracking:** Seamless audit emission (`hub:cockpit:viewed`, `hub:widgets:reconfigured`, `hub:shortcut:triggered`) enforcing full institutional traceability.

---

## 2. Super-App Architecture & Subsystem Aggregation Flow

```mermaid
flowchart TD
    subgraph ClientLayers["Client Experience Layers"]
        Web[Web Cockpit - Svelte 5 / Shadcn]
        Mobile[Mobile Super-App - React Native]
    end

    subgraph HubCore["Subsystem 17: The Hub Root Super-App"]
        Router[Persona Router & Authorizer]
        Aggregator[Cross-Subsystem KPI Aggregator]
        WidgetEngine[Widget Layout & Theme Engine]
        ShortcutEngine[1-Tap Quick Action Dispatcher]
    end

    subgraph Subsystems["Topological Modular Monolith (Ranks 1 to 16)"]
        S1[1. Auth & RBAC]
        S2[2. Audit & Merkle Ledger]
        S3[3. Onboarding & Sequence Engine]
        S4[4. Attendance & Geofencing]
        S5[5. Central Billing & Double-Entry Ledger]
        S6[6. Campus Notices & Acks]
        S7[7. Hostel & Curfew Passes]
        S8[8. Mess & QR Dining Tokens]
        S9[9. E-Library & RFID Borrows]
        S10[10. Study Hub & LMS]
        S11[11. Progress Tracker & Mentorship]
        S12[12. Campus Events & Gate Check-In]
        S13[13. Helpdesk & SLA Tracker]
        S14[14. SOS Emergency Radar]
        S15[15. Zero-Knowledge Whistleblower]
        S16[16. Public Web Portal & Prospect CRM]
    end

    Web --> Router
    Mobile --> Router
    Router --> Aggregator
    Router --> WidgetEngine
    Router --> ShortcutEngine
    Aggregator --> Subsystems
    ShortcutEngine -.-> Subsystems
```

---

## 3. Persona Cockpit Sequence Diagram

```mermaid
sequenceDiagram
    autonumber
    actor User as Campus User (Student / Faculty / Warden / Admin)
    participant Gateway as API Gateway / Chi
    participant HubSvc as Hub Service
    participant HubRepo as Hub Repository
    participant Subsystems as Subsystem Engines (1-16)
    participant Audit as Central Audit Subsystem

    User->>Gateway: GET /api/v1/hub/cockpit?persona=STUDENT
    Gateway->>HubSvc: GetPersonaCockpit(tenantID, userID, STUDENT)
    HubSvc->>HubRepo: GetOrCreateDashboard(tenantID, userID, STUDENT)
    HubSvc->>HubRepo: GetAggregatedMetrics(tenantID, userID, STUDENT)
    HubSvc->>HubRepo: ListShortcuts(tenantID, STUDENT)
    HubSvc->>Audit: Enqueue(hub:cockpit:viewed [STUDENT])
    HubSvc-->>Gateway: PersonaDashboardView (Dashboard, Metrics, Shortcuts)
    Gateway-->>User: 200 OK JSON

    User->>Gateway: POST /api/v1/hub/shortcuts/trigger (TRIGGER_SOS)
    Gateway->>HubSvc: TriggerShortcut(TRIGGER_SOS)
    HubSvc->>Audit: Enqueue(hub:shortcut:triggered [TRIGGER_SOS])
    HubSvc-->>Gateway: Shortcut Navigation Object (Route: /sos)
    Gateway-->>User: 200 OK
```

---

## 4. REST API Transport Endpoints

| HTTP Method | Path | Description | Access Control |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/hub/cockpit` | Retrieve consolidated persona dashboard, live KPIs, and shortcuts | Authenticated User |
| `POST` | `/api/v1/hub/widgets/reconfigure` | Save reordered and customized widget layout configuration | Authenticated User |
| `POST` | `/api/v1/hub/theme` | Update visual dashboard theme preference | Authenticated User |
| `GET` | `/api/v1/hub/shortcuts` | List persona-scoped 1-tap quick action shortcuts | Authenticated User |
| `POST` | `/api/v1/hub/shortcuts/trigger` | Execute and log 1-tap quick action trigger | Authenticated User |

---

## 5. PostgreSQL 18 / Prisma 8 Data Models

- `HubPersonaDashboard`: User-specific customized cockpit layout for a given persona.
- `HubWidgetConfig`: Individual widget card configuration (size, order index, enabled toggle, config JSON).
- `HubQuickActionShortcut`: 1-tap navigation and action trigger definitions per persona.
