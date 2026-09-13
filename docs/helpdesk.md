# Subsystem 13: Application & Query Helpdesk (`helpdesk`)

## 1. Overview & Business Context

The **Application & Query Helpdesk** is the central institutional support engine for the Campus Management System (CMS). It streamlines student and staff queries across academics, hostel maintenance, central fee billing, and IT infrastructure into unified, SLA-monitored resolution workflows.

Key architectural pillars:
- **Deterministic Ticket Identification:** Human-readable format (`HD-YYYY-NNNNN`) partitioned per tenant.
- **SLA Breach Surveillance Radar:** Dynamic countdown timers based on category response & resolution hour quotas.
- **Staff-Only Internal Notes:** Role-gated messages hidden from applicants to facilitate internal deliberation before official responses.
- **Dynamic Conversation State Engine:** Auto-transitions ticket status between `IN_PROGRESS` and `WAITING_FOR_APPLICANT` upon replies.
- **SLA Escalation Protocol:** 1-tap urgent escalation bumping priority to `URGENT` with audit trail dispatching.
- **Customer Satisfaction (CSAT) Rating:** 5-star rating and feedback loop gated to resolved or closed tickets.

---

## 2. Ticket State Machine & Lifecycle

```mermaid
stateDiagram-v2
    [*] --> OPEN: Student / Applicant creates ticket
    OPEN --> IN_PROGRESS: Staff assigned or begins review
    IN_PROGRESS --> WAITING_FOR_APPLICANT: Staff posts public question
    WAITING_FOR_APPLICANT --> IN_PROGRESS: Applicant responds
    IN_PROGRESS --> RESOLVED: Staff provides resolution
    WAITING_FOR_APPLICANT --> RESOLVED: Staff resolves ticket
    OPEN --> RESOLVED: Direct resolution
    RESOLVED --> CLOSED: Requester confirms / Auto-close
    RESOLVED --> IN_PROGRESS: Requester disputes resolution (Reopen)
    CLOSED --> [*]
```

---

## 3. SLA Escalation & Live Resolution Sequence

```mermaid
sequenceDiagram
    autonumber
    actor Student as Student / Applicant
    actor Staff as Helpdesk Staff
    participant Gateway as API Gateway / Chi
    participant Svc as Helpdesk Service
    participant Repo as Helpdesk Repository
    participant Audit as Central Audit Subsystem

    Student->>Gateway: POST /api/v1/helpdesk/tickets (Create Ticket)
    Gateway->>Svc: CreateTicket(req)
    Svc->>Repo: Fetch Category SLA & Sequence
    Svc->>Repo: Persist HelpdeskTicket (Status: OPEN, SLA Due At)
    Svc->>Audit: Enqueue(helpdesk:ticket:created)
    Gateway-->>Student: 201 Created (HD-2026-00042)

    Staff->>Gateway: POST /api/v1/helpdesk/tickets/{id}/assign
    Gateway->>Svc: AssignTicket(staffID)
    Svc->>Repo: Update Ticket (Assigned, Status: IN_PROGRESS)
    Svc->>Audit: Enqueue(helpdesk:ticket:assigned)

    Staff->>Gateway: POST /api/v1/helpdesk/tickets/{id}/messages (isInternalNote: true)
    Gateway->>Svc: AddMessage(Internal Note)
    Svc->>Repo: Persist HelpdeskMessage
    Note over Staff, Svc: Internal note is masked from student queries

    Staff->>Gateway: POST /api/v1/helpdesk/tickets/{id}/status (RESOLVED)
    Gateway->>Svc: UpdateTicketStatus(RESOLVED)
    Svc->>Repo: Set resolvedAt timestamp

    Student->>Gateway: POST /api/v1/helpdesk/tickets/{id}/rate
    Gateway->>Svc: RateTicket(Rating: 5 Stars, Feedback)
    Svc->>Repo: Update Ticket CSAT
    Svc->>Audit: Enqueue(helpdesk:ticket:rated)
    Gateway-->>Student: 200 OK
```

---

## 4. REST API Transport Endpoints

| HTTP Method | Path | Description | Access Control |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/helpdesk/categories` | Create Helpdesk Category & SLA config | Staff / Admin |
| `GET` | `/api/v1/helpdesk/categories` | List Active Categories | Public / Authenticated |
| `POST` | `/api/v1/helpdesk/tickets` | Submit new inquiry ticket | Authenticated User |
| `GET` | `/api/v1/helpdesk/tickets` | Filter tickets (status, priority, category) | Authenticated User |
| `GET` | `/api/v1/helpdesk/tickets/{id}` | Get ticket details and SLA timer | Requester / Staff |
| `POST` | `/api/v1/helpdesk/tickets/{id}/assign` | Assign ticket to support officer | Staff / Admin |
| `POST` | `/api/v1/helpdesk/tickets/{id}/status` | Transition ticket state | Staff / Requester |
| `POST` | `/api/v1/helpdesk/tickets/{id}/messages` | Post message (public or internal note) | Requester / Staff |
| `GET` | `/api/v1/helpdesk/tickets/{id}/messages` | Get message thread (internal notes masked for students) | Requester / Staff |
| `POST` | `/api/v1/helpdesk/tickets/{id}/escalate` | Trigger SLA Urgent Escalation | Requester / Staff |
| `POST` | `/api/v1/helpdesk/tickets/{id}/rate` | Submit CSAT satisfaction score (1-5 stars) | Ticket Requester |

---

## 5. PostgreSQL 18 / Prisma 8 Data Models

- `HelpdeskCategory`: Category taxonomy with default priority, SLA response and resolution hour thresholds.
- `HelpdeskTicket`: Core ticket entity tracking requester, assignee, SLA deadline, status lifecycle, and CSAT rating.
- `HelpdeskMessage`: Conversation thread elements with staff reply flag, internal note masking, and attachments.
- `HelpdeskEscalation`: Escalation incident logs tracking reasons, supervisor alerts, and escalation timestamps.
