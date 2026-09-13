# Subsystem 14: SOS & Emergency Response (`sos`)

## 1. Overview & Business Context

The **SOS & Emergency Response System** provides instantaneous distress alert handling and telemetry dispatch for the campus community. Engineered for high reliability, it supports 1-tap emergency broadcasts from mobile devices with geographic coordinates, automated siren/push notifications, multi-role responder team dispatch, and comprehensive post-incident resolution audits.

Key architectural pillars:
- **Instantaneous 1-Tap Trigger:** Single-touch distress broadcast capturing live GPS coordinates with latitude/longitude bounds validation.
- **Deterministic Alert Identifiers:** Partitioned sequence numbering (`SOS-YYYY-NNNNN`) per tenant.
- **Multi-Role Responder Fleet:** Dispatch and track campus security patrols, ambulance/paramedic teams, hostel wardens, and fire safety officers.
- **Arrival Checkpoint Surveillance:** Live transition tracking (`ASSIGNED` -> `EN_ROUTE` -> `ON_SCENE` -> `COMPLETED`).
- **Resolution & False-Alarm Logging:** Mandatory resolution reports and post-incident compliance logging into the immutable audit ledger.

---

## 2. Emergency Incident State Machine

```mermaid
stateDiagram-v2
    [*] --> TRIGGERED: 1-Tap Mobile SOS Broadcast
    TRIGGERED --> ACKNOWLEDGED: Control Room Operator Acknowledges
    TRIGGERED --> DISPATCHED: Immediate Responder Dispatch
    ACKNOWLEDGED --> DISPATCHED: Responder Teams Assigned
    DISPATCHED --> ON_SCENE: First Responder Arrives at GPS Pin
    ON_SCENE --> RESOLVED: Situation Neutralized & Logged
    ON_SCENE --> FALSE_ALARM: Accidental / Test Trigger
    TRIGGERED --> FALSE_ALARM: Caller Cancels
    RESOLVED --> [*]
    FALSE_ALARM --> [*]
```

---

## 3. Telemetry Broadcast & Dispatch Sequence

```mermaid
sequenceDiagram
    autonumber
    actor Caller as Student / Staff in Distress
    actor Operator as Control Room Officer
    actor Responder as Paramedic / Security Patrol
    participant Gateway as API Gateway / Chi
    participant Svc as SOS Service
    participant Repo as SOS Repository
    participant Audit as Central Audit Subsystem

    Caller->>Gateway: POST /api/v1/sos/trigger (Type: MEDICAL, Lat, Lng, Location)
    Gateway->>Svc: TriggerSOS(req)
    Svc->>Repo: Persist SOSIncident (Status: TRIGGERED)
    Svc->>Repo: Create SOSTelemetryBroadcast (PUSH)
    Svc->>Audit: Enqueue(sos:alert:triggered)
    Gateway-->>Caller: 201 Created (SOS-2026-00012)

    Operator->>Gateway: POST /api/v1/sos/incidents/{id}/acknowledge
    Gateway->>Svc: AcknowledgeIncident()
    Svc->>Audit: Enqueue(sos:alert:acknowledged)

    Operator->>Gateway: POST /api/v1/sos/incidents/{id}/dispatch (Role: PARAMEDIC)
    Gateway->>Svc: DispatchResponder()
    Svc->>Repo: Persist SOSDispatchResponder (Status: ASSIGNED)
    Svc->>Audit: Enqueue(sos:alert:dispatched)

    Responder->>Gateway: POST /api/v1/sos/incidents/{id}/responders/{id}/status (ON_SCENE)
    Gateway->>Svc: UpdateResponderStatus(ON_SCENE)
    Svc->>Repo: Update Arrival Timestamp & Incident ON_SCENE

    Responder->>Gateway: POST /api/v1/sos/incidents/{id}/resolve (Notes: Patient Stable)
    Gateway->>Svc: ResolveIncident()
    Svc->>Repo: Set resolvedAt & resolutionNotes
    Svc->>Audit: Enqueue(sos:alert:resolved)
    Gateway-->>Responder: 200 OK
```

---

## 4. REST API Transport Endpoints

| HTTP Method | Path | Description | Access Control |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/sos/trigger` | 1-Tap SOS Distress Trigger | Authenticated User |
| `GET` | `/api/v1/sos/incidents` | List & filter active/historical incidents | Security / Staff |
| `GET` | `/api/v1/sos/incidents/{id}` | Get incident telemetry & status | Requester / Staff |
| `POST` | `/api/v1/sos/incidents/{id}/acknowledge` | Acknowledge emergency alert | Control Room Operator |
| `POST` | `/api/v1/sos/incidents/{id}/dispatch` | Dispatch response units | Control Room Operator |
| `POST` | `/api/v1/sos/incidents/{id}/responders/{id}/status` | Update responder mission checkpoint (`EN_ROUTE`, `ON_SCENE`, `COMPLETED`) | Assigned Responder |
| `POST` | `/api/v1/sos/incidents/{id}/resolve` | Log resolution notes or false alarm | Responder / Operator |
| `GET` | `/api/v1/sos/incidents/{id}/responders` | List dispatched units for incident | Requester / Staff |

---

## 5. PostgreSQL 18 / Prisma 8 Data Models

- `SOSIncident`: Emergency alert entity with GPS coordinates, emergency type, status lifecycle, and resolution notes.
- `SOSDispatchResponder`: Responder missions mapping officers/paramedics to incidents with arrival timestamps.
- `SOSTelemetryBroadcast`: Multi-channel telemetry broadcasts (push notifications, SMS, public PA/sirens).
