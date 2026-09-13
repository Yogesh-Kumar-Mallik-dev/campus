# Subsystem 15: Anonymity & Whistleblower System (`whistleblower`)

## 1. Overview & Business Context

The **Anonymity & Whistleblower System** delivers a secure, zero-knowledge grievance intake channel for the institution. Engineered to meet UGC anti-ragging mandates and institutional ethics oversight, it enables students and employees to submit reports on ragging, financial irregularities, academic corruption, and harassment with cryptographic anonymity.

Key architectural pillars:
- **Zero-Knowledge Token Architecture:** High-entropy secret token (`WB-tok-<hex>`) generated on submission; server stores only the SHA-256 hash `trackingHash`.
- **Decoupled Identity Guard:** Submissions are not tied to authenticated sessions, user profiles, or network IP identifiers.
- **Confidential 2-Way Channel:** Reporters can follow up, view status updates, and exchange encrypted dialogue with the committee using only their secret token.
- **Committee Investigation Workflow:** Structured state machine (`SUBMITTED` -> `UNDER_INVESTIGATION` -> `EVIDENCE_REQUESTED` -> `RESOLVED` / `REJECTED`).
- **Immutable Compliance Ledger:** All investigation events and disciplinary outcomes are audited anonymously without leaking complainant metadata.

---

## 2. Investigation State Machine

```mermaid
stateDiagram-v2
    [*] --> SUBMITTED: Anonymous Report Intake (Token Issued)
    SUBMITTED --> UNDER_INVESTIGATION: Assigned to Anti-Ragging / Vigilance Squad
    SUBMITTED --> REJECTED: Dismissed / Unsubstantiated
    UNDER_INVESTIGATION --> EVIDENCE_REQUESTED: Squad asks for dates/evidence
    EVIDENCE_REQUESTED --> UNDER_INVESTIGATION: Complainant submits anonymous reply
    UNDER_INVESTIGATION --> RESOLVED: Investigation Complete & Action Taken
    UNDER_INVESTIGATION --> REJECTED: Findings conclude no violation
    RESOLVED --> [*]
    REJECTED --> [*]
```

---

## 3. Cryptographic Token & Anonymous Dialogue Sequence

```mermaid
sequenceDiagram
    autonumber
    actor Complainant as Anonymous Student
    actor Committee as Vigilance / Anti-Ragging Committee
    participant Gateway as API Gateway / Chi
    participant Svc as Whistleblower Service
    participant Repo as Whistleblower Repository
    participant Audit as Central Audit Subsystem

    Complainant->>Gateway: POST /api/v1/whistleblower/reports (Category, Narrative)
    Gateway->>Svc: SubmitReport(req)
    Svc->>Svc: Generate rawToken & SHA-256(rawToken)
    Svc->>Repo: Persist WhistleblowerReport (TrackingHash)
    Svc->>Audit: Enqueue(whistleblower:report:submitted [ANONYMOUS])
    Gateway-->>Complainant: 201 Created (ReportNumber, rawToken)

    Complainant->>Gateway: POST /api/v1/whistleblower/reports/lookup (rawToken)
    Gateway->>Svc: LookupReportByToken(rawToken)
    Svc->>Repo: Fetch by SHA-256(rawToken)
    Gateway-->>Complainant: 200 OK (Status, Messages)

    Committee->>Gateway: POST /api/v1/whistleblower/reports/{id}/messages (Ask Evidence)
    Gateway->>Svc: AddMessage(Sender: COMMITTEE)
    Svc->>Repo: Persist WhistleblowerMessage

    Complainant->>Gateway: POST /api/v1/whistleblower/reports/{id}/messages (rawToken, Reply)
    Gateway->>Svc: AddMessage(Verify rawToken matches TrackingHash)
    Svc->>Repo: Persist WhistleblowerMessage
    Gateway-->>Complainant: 201 Created
```

---

## 4. REST API Transport Endpoints

| HTTP Method | Path | Description | Access Control |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/whistleblower/reports` | Submit anonymous grievance & receive tracking token | Public / Anonymous |
| `POST` | `/api/v1/whistleblower/reports/lookup` | Lookup report status & dialogue using secret token | Public / Anonymous |
| `POST` | `/api/v1/whistleblower/reports/{id}/messages` | Post anonymous reply (Token) or committee message (Auth) | Reporter Token / Committee |
| `GET` | `/api/v1/whistleblower/committee/reports` | List & filter confidential reports | Authorized Committee |
| `POST` | `/api/v1/whistleblower/committee/reports/{id}/assign` | Assign report to vigilance squad | Committee Head |
| `POST` | `/api/v1/whistleblower/committee/reports/{id}/status` | Update investigation status & findings log | Assigned Investigator |

---

## 5. PostgreSQL 18 / Prisma 8 Data Models

- `WhistleblowerReport`: Confidential report entity storing SHA-256 tracking token hash, severity, and investigation findings.
- `WhistleblowerMessage`: Anonymous conversation thread items exchanged between complainant and committee.
- `WhistleblowerEvidence`: Attached documentary evidence with SHA-256 file integrity hashing.
