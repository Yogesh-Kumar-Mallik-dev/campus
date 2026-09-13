/**
 * BLOCK_API_HOSTEL_HANDLER_001
 * Subsystem: Rank 7 - Hostel Management System (hostel)
 * Purpose:   HTTP REST transport handlers with RFC 7807 problem details for residential blocks, bed allocations, and gate passes.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package hostel

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"campus/api/problem"
	"campus/backend/hostel"
)

type Handler struct {
	service *hostel.Service
}

func NewHandler(service *hostel.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/hostel", func(r chi.Router) {
		// Blocks
		r.Post("/blocks", h.CreateBlock)
		r.Get("/blocks", h.ListBlocks)

		// Rooms
		r.Post("/rooms", h.CreateRoom)
		r.Get("/rooms", h.ListRooms)

		// Beds
		r.Post("/beds", h.CreateBed)
		r.Get("/beds", h.ListBeds)

		// Allocations
		r.Post("/allocations", h.AllocateBed)
		r.Post("/allocations/{id}/vacate", h.VacateBed)

		// Gate Passes
		r.Post("/gate-passes", h.ApplyGatePass)
		r.Get("/gate-passes", h.ListGatePasses)
		r.Post("/gate-passes/{id}/review", h.ReviewGatePass)
		r.Post("/gate-passes/{id}/exit", h.RecordGatePassExit)
		r.Post("/gate-passes/{id}/return", h.RecordGatePassReturn)

		// Incidents
		r.Post("/incidents", h.LogIncident)
		r.Get("/incidents", h.ListIncidents)
	})
}

// CreateBlock handles POST /api/v1/hostel/blocks
func (h *Handler) CreateBlock(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID    string        `json:"tenant_id"`
		Name        string        `json:"name"`
		Code        string        `json:"code"`
		Gender      hostel.Gender `json:"gender"`
		TotalFloors int           `json:"total_floors"`
		TotalRooms  int           `json:"total_rooms"`
		Capacity    int           `json:"capacity"`
		WardenID    *string       `json:"warden_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	block, err := h.service.CreateBlock(r.Context(), req.TenantID, req.Name, req.Code, req.Gender, req.TotalFloors, req.TotalRooms, req.Capacity, req.WardenID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"block": block})
}

// ListBlocks handles GET /api/v1/hostel/blocks
func (h *Handler) ListBlocks(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	if tenantID == "" {
		problem.BadRequest(w, r, "tenant_id query parameter is required", "MISSING_TENANT_ID")
		return
	}

	blocks, err := h.service.ListBlocks(r.Context(), tenantID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"blocks": blocks})
}

// CreateRoom handles POST /api/v1/hostel/rooms
func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID    string          `json:"tenant_id"`
		BlockID     string          `json:"block_id"`
		RoomNumber  string          `json:"room_number"`
		FloorNumber int             `json:"floor_number"`
		RoomType    hostel.RoomType `json:"room_type"`
		IsAC        bool            `json:"is_ac"`
		BaseFee     float64         `json:"base_fee_per_semester"`
		MaxBeds     int             `json:"max_beds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	room, err := h.service.CreateRoom(r.Context(), req.TenantID, req.BlockID, req.RoomNumber, req.FloorNumber, req.RoomType, req.IsAC, req.BaseFee, req.MaxBeds)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"room": room})
}

// ListRooms handles GET /api/v1/hostel/rooms
func (h *Handler) ListRooms(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	blockID := r.URL.Query().Get("block_id")
	if tenantID == "" || blockID == "" {
		problem.BadRequest(w, r, "tenant_id and block_id are required", "MISSING_PARAMS")
		return
	}

	rooms, err := h.service.ListRooms(r.Context(), tenantID, blockID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"rooms": rooms})
}

// CreateBed handles POST /api/v1/hostel/beds
func (h *Handler) CreateBed(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID  string `json:"tenant_id"`
		RoomID    string `json:"room_id"`
		BedNumber string `json:"bed_number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	bed, err := h.service.CreateBed(r.Context(), req.TenantID, req.RoomID, req.BedNumber)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"bed": bed})
}

// ListBeds handles GET /api/v1/hostel/beds
func (h *Handler) ListBeds(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	roomID := r.URL.Query().Get("room_id")
	if tenantID == "" || roomID == "" {
		problem.BadRequest(w, r, "tenant_id and room_id are required", "MISSING_PARAMS")
		return
	}

	beds, err := h.service.ListBeds(r.Context(), tenantID, roomID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"beds": beds})
}

// AllocateBed handles POST /api/v1/hostel/allocations
func (h *Handler) AllocateBed(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID     string `json:"tenant_id"`
		BedID        string `json:"bed_id"`
		StudentID    string `json:"student_id"`
		AcademicYear string `json:"academic_year"`
		Semester     int    `json:"semester"`
		Remarks      string `json:"remarks"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	alloc, err := h.service.AllocateBed(r.Context(), req.TenantID, req.BedID, req.StudentID, req.AcademicYear, req.Semester, req.Remarks)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"allocation": alloc})
}

// VacateBed handles POST /api/v1/hostel/allocations/{id}/vacate
func (h *Handler) VacateBed(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		TenantID string `json:"tenant_id"`
		Remarks  string `json:"remarks"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	alloc, err := h.service.VacateBed(r.Context(), req.TenantID, id, req.Remarks)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"allocation": alloc})
}

// ApplyGatePass handles POST /api/v1/hostel/gate-passes
func (h *Handler) ApplyGatePass(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID         string    `json:"tenant_id"`
		StudentID        string    `json:"student_id"`
		BlockID          string    `json:"block_id"`
		Reason           string    `json:"reason"`
		Destination      string    `json:"destination"`
		EmergencyContact string    `json:"emergency_contact"`
		ExpectedOutAt    time.Time `json:"expected_out_at"`
		ExpectedInAt     time.Time `json:"expected_in_at"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	gp, err := h.service.ApplyGatePass(r.Context(), req.TenantID, req.StudentID, req.BlockID, req.Reason, req.Destination, req.EmergencyContact, req.ExpectedOutAt, req.ExpectedInAt)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"gate_pass": gp})
}

// ListGatePasses handles GET /api/v1/hostel/gate-passes
func (h *Handler) ListGatePasses(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	studentID := r.URL.Query().Get("student_id")
	blockID := r.URL.Query().Get("block_id")

	var status *hostel.GatePassStatus
	if s := r.URL.Query().Get("status"); s != "" {
		parsed := hostel.GatePassStatus(s)
		status = &parsed
	}

	passes, err := h.service.ListGatePasses(r.Context(), tenantID, studentID, blockID, status)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"gate_passes": passes})
}

// ReviewGatePass handles POST /api/v1/hostel/gate-passes/{id}/review
func (h *Handler) ReviewGatePass(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		TenantID        string `json:"tenant_id"`
		WardenID        string `json:"warden_id"`
		Approve         bool   `json:"approve"`
		RejectionReason string `json:"rejection_reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	gp, err := h.service.ReviewGatePass(r.Context(), req.TenantID, id, req.WardenID, req.Approve, req.RejectionReason)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"gate_pass": gp})
}

// RecordGatePassExit handles POST /api/v1/hostel/gate-passes/{id}/exit
func (h *Handler) RecordGatePassExit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		TenantID    string    `json:"tenant_id"`
		ActualOutAt time.Time `json:"actual_out_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}
	if req.ActualOutAt.IsZero() {
		req.ActualOutAt = time.Now().UTC()
	}

	gp, err := h.service.RecordGatePassExit(r.Context(), req.TenantID, id, req.ActualOutAt)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"gate_pass": gp})
}

// RecordGatePassReturn handles POST /api/v1/hostel/gate-passes/{id}/return
func (h *Handler) RecordGatePassReturn(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		TenantID   string    `json:"tenant_id"`
		WardenID   string    `json:"warden_id"`
		ActualInAt time.Time `json:"actual_in_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}
	if req.ActualInAt.IsZero() {
		req.ActualInAt = time.Now().UTC()
	}

	gp, incident, err := h.service.RecordGatePassReturn(r.Context(), req.TenantID, id, req.WardenID, req.ActualInAt)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	resp := map[string]interface{}{"gate_pass": gp}
	if incident != nil {
		resp["curfew_incident"] = incident
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// LogIncident handles POST /api/v1/hostel/incidents
func (h *Handler) LogIncident(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID     string                 `json:"tenant_id"`
		StudentID    string                 `json:"student_id"`
		BlockID      string                 `json:"block_id"`
		WardenID     string                 `json:"warden_id"`
		IncidentType hostel.IncidentType    `json:"incident_type"`
		Severity     hostel.IncidentSeverity `json:"severity"`
		Title        string                 `json:"title"`
		Description  string                 `json:"description"`
		ActionTaken  string                 `json:"action_taken"`
		FineAmount   float64                `json:"fine_amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem.BadRequest(w, r, "malformed request JSON payload", "INVALID_PAYLOAD")
		return
	}
	if req.TenantID == "" {
		req.TenantID = r.Header.Get("X-Tenant-ID")
	}

	inc, err := h.service.LogIncident(r.Context(), req.TenantID, req.StudentID, req.BlockID, req.WardenID, req.IncidentType, req.Severity, req.Title, req.Description, req.ActionTaken, req.FineAmount)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"incident": inc})
}

// ListIncidents handles GET /api/v1/hostel/incidents
func (h *Handler) ListIncidents(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = r.Header.Get("X-Tenant-ID")
	}
	studentID := r.URL.Query().Get("student_id")
	blockID := r.URL.Query().Get("block_id")

	incidents, err := h.service.ListIncidents(r.Context(), tenantID, studentID, blockID)
	if err != nil {
		h.handleDomainError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"incidents": incidents})
}

func (h *Handler) handleDomainError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, hostel.ErrBlockNotFound),
		errors.Is(err, hostel.ErrRoomNotFound),
		errors.Is(err, hostel.ErrBedNotFound),
		errors.Is(err, hostel.ErrAllocationNotFound),
		errors.Is(err, hostel.ErrGatePassNotFound),
		errors.Is(err, hostel.ErrIncidentNotFound):
		problem.NotFound(w, r, err.Error(), "RESOURCE_NOT_FOUND")

	case errors.Is(err, hostel.ErrBlockCodeExists),
		errors.Is(err, hostel.ErrRoomNumberExists),
		errors.Is(err, hostel.ErrBedNumberExists),
		errors.Is(err, hostel.ErrBedUnavailable),
		errors.Is(err, hostel.ErrStudentAlreadyAllocated):
		problem.Conflict(w, r, err.Error(), "RESOURCE_CONFLICT")

	case errors.Is(err, hostel.ErrInvalidGatePassState),
		errors.Is(err, hostel.ErrInvalidTimeRange):
		problem.UnprocessableEntity(w, r, err.Error(), "UNPROCESSABLE_ENTITY", nil)

	case errors.Is(err, hostel.ErrUnauthorizedWarden):
		problem.Forbidden(w, r, err.Error(), "FORBIDDEN_ACCESS")

	default:
		var domErr *hostel.DomainError
		if errors.As(err, &domErr) {
			problem.BadRequest(w, r, domErr.Error(), "BAD_REQUEST")
			return
		}
		problem.InternalServerError(w, r, "an unexpected hostel error occurred")
	}
}
