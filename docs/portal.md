# Subsystem 16: Public Web Portal (`portal`)

## 1. Overview & Business Context

The **Public Web Portal** subsystem serves as the institution's public-facing digital storefront and prospective student admissions gateway. It manages institutional landing pages, academic program catalogs across undergraduate, postgraduate, diploma, and doctoral tracks, admissions cycle announcements, and prospect lead capture.

Key architectural pillars:
- **Public Prospectus Catalog:** Dynamic program listing including degree levels (`UG`, `PG`, `DIPLOMA`, `PHD`), duration, total semesters, minimum eligibility criteria, and annual fee schedules.
- **Inquiry & Lead Triage Workflow:** Captures prospective candidate contact information and areas of interest, transitioning leads through a structured status machine (`NEW` -> `CONTACTED` -> `CONVERTED` -> `CLOSED`).
- **Counselor Assignment & Notes:** Equips admissions counselors with CRM capabilities to assign leads, record call transcripts, and track enrollment conversions.
- **Tenant Institutional Brand Showcase:** Publishes campus location details, accreditation badges, and admissions cycle status.
- **Audit & Conversion Analytics:** Publishes audit events on lead capture and conversion for marketing ROI and compliance tracking.

---

## 2. Admissions Lead State Machine

```mermaid
stateDiagram-v2
    [*] --> NEW: Prospect Submits Inquiry Lead Form
    NEW --> CONTACTED: Admissions Counselor Reaches Out (Call/Email)
    CONTACTED --> CONVERTED: Candidate Submits Application / Secures Admission
    CONTACTED --> CLOSED: Lead Ineligible / Opted Out
    NEW --> CLOSED: Duplicate / Invalid Lead
    CONVERTED --> [*]
    CLOSED --> [*]
```

---

## 3. Prospect Inquiry & Counselor Follow-Up Sequence

```mermaid
sequenceDiagram
    autonumber
    actor Prospect as Prospective Student
    actor Counselor as Admissions Counselor
    participant Gateway as API Gateway / Chi
    participant Svc as Portal Service
    participant Repo as Portal Repository
    participant Audit as Central Audit Subsystem

    Prospect->>Gateway: POST /api/v1/portal/inquiries (Name, Email, Phone, Program)
    Gateway->>Svc: CreateInquiry(req)
    Svc->>Repo: Persist PortalPublicInquiry (Status: NEW)
    Svc->>Audit: Enqueue(portal:inquiry:created)
    Gateway-->>Prospect: 201 Created (InquiryID)

    Counselor->>Gateway: GET /api/v1/portal/counselor/inquiries?status=NEW
    Gateway->>Svc: ListInquiries(status: NEW)
    Svc->>Repo: Fetch inquiries
    Gateway-->>Counselor: 200 OK (Inquiries List)

    Counselor->>Gateway: POST /api/v1/portal/counselor/inquiries/{id}/assign (CounselorID)
    Gateway->>Svc: AssignCounselor(id, counselorID)
    Svc->>Repo: Update assignedCounselor
    Gateway-->>Counselor: 200 OK

    Counselor->>Gateway: POST /api/v1/portal/counselor/inquiries/{id}/status (Status: CONVERTED, Notes)
    Gateway->>Svc: UpdateInquiryStatus(id, CONVERTED, notes)
    Svc->>Repo: Update status & notes
    Svc->>Audit: Enqueue(portal:inquiry:converted)
    Gateway-->>Counselor: 200 OK
```

---

## 4. REST API Transport Endpoints

| HTTP Method | Path | Description | Access Control |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/portal/landing` | Fetch published institutional landing page content | Public |
| `GET` | `/api/v1/portal/programs` | List academic degree offerings with degree/search filters | Public |
| `POST` | `/api/v1/portal/inquiries` | Submit prospective student admissions inquiry | Public |
| `GET` | `/api/v1/portal/counselor/inquiries` | List & filter prospect inquiry leads | Admissions Staff |
| `POST` | `/api/v1/portal/counselor/inquiries/{id}/assign` | Assign prospect lead to counselor | Admissions Lead |
| `POST` | `/api/v1/portal/counselor/inquiries/{id}/status` | Update inquiry status and counselor notes | Assigned Counselor |

---

## 5. PostgreSQL 18 / Prisma 8 Data Models

- `PortalLandingPage`: Institutional landing page content, admissions cycle status, and contact coordinates.
- `PortalProgramCatalog`: Degree program offerings, accreditation, eligibility, and fee structure.
- `PortalPublicInquiry`: Prospect candidate leads, counseling status history, and conversion tracking.
