/**
 * BLOCK_ONBOARDING_SERVICE_001
 * Subsystem: Rank 3 - Student & Staff Registration System (onboarding)
 * Purpose:   High-level orchestration for applicant KYC, document verification, roll number generation, and enrollment.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package onboarding

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"campus/backend/audit"
)

type CreateApplicantCommand struct {
	TenantID          string           `json:"tenant_id"`
	Type              OnboardingType   `json:"type"`
	FirstName         string           `json:"first_name"`
	LastName          string           `json:"last_name"`
	Email             string           `json:"email"`
	Phone             string           `json:"phone"`
	DateOfBirth       string           `json:"date_of_birth"`
	Gender            Gender           `json:"gender"`
	BloodGroup        string           `json:"blood_group,omitempty"`
	Nationality       string           `json:"nationality"`
	Emergency         EmergencyContact `json:"emergency_contact"`
	Address           Address          `json:"address"`
	ProgramID         string           `json:"program_id,omitempty"`
	DepartmentID      string           `json:"department_id,omitempty"`
	AcademicYear      string           `json:"academic_year"`
	TargetDesignation string           `json:"target_designation,omitempty"`
	Guardian          *Guardian        `json:"guardian,omitempty"`
}

type UpdateApplicantCommand struct {
	TenantID          string            `json:"tenant_id"`
	FirstName         *string           `json:"first_name,omitempty"`
	LastName          *string           `json:"last_name,omitempty"`
	Phone             *string           `json:"phone,omitempty"`
	DateOfBirth       *string           `json:"date_of_birth,omitempty"`
	Gender            *Gender           `json:"gender,omitempty"`
	BloodGroup        *string           `json:"blood_group,omitempty"`
	Nationality       *string           `json:"nationality,omitempty"`
	Emergency         *EmergencyContact `json:"emergency_contact,omitempty"`
	Address           *Address          `json:"address,omitempty"`
	ProgramID         *string           `json:"program_id,omitempty"`
	DepartmentID      *string           `json:"department_id,omitempty"`
	AcademicYear      *string           `json:"academic_year,omitempty"`
	TargetDesignation *string           `json:"target_designation,omitempty"`
	Guardian          *Guardian         `json:"guardian,omitempty"`
}

type AttachDocumentCommand struct {
	TenantID     string       `json:"tenant_id"`
	ApplicantID  string       `json:"applicant_id"`
	DocumentType DocumentType `json:"document_type"`
	FileKey      string       `json:"file_key"`
	FileName     string       `json:"file_name"`
	FileSize     int64        `json:"file_size"`
	MIMEType     string       `json:"mime_type"`
}

type EnrollStudentCommand struct {
	TenantID    string `json:"tenant_id"`
	ApplicantID string `json:"applicant_id"`
	UserID      string `json:"user_id"`
	CohortID    string `json:"cohort_id"`
	AdminID     string `json:"admin_id"`
}

type ProvisionStaffCommand struct {
	TenantID    string `json:"tenant_id"`
	ApplicantID string `json:"applicant_id"`
	UserID      string `json:"user_id"`
	RoleKey     string `json:"role_key"` // e.g. "faculty", "warden", "librarian", "admin"
	AdminID     string `json:"admin_id"`
}

type Service struct {
	applicantRepo ApplicantRepository
	docRepo       DocumentRepository
	academicRepo  AcademicRepository
	profileRepo   ProfileRepository
	sequenceGen   SequenceGenerator
	auditSub      audit.Subscriber
}

func NewService(
	applicantRepo ApplicantRepository,
	docRepo DocumentRepository,
	academicRepo AcademicRepository,
	profileRepo ProfileRepository,
	sequenceGen SequenceGenerator,
	auditSub audit.Subscriber,
) *Service {
	return &Service{
		applicantRepo: applicantRepo,
		docRepo:       docRepo,
		academicRepo:  academicRepo,
		profileRepo:   profileRepo,
		sequenceGen:   sequenceGen,
		auditSub:      auditSub,
	}
}

func (s *Service) CreateDraft(ctx context.Context, cmd CreateApplicantCommand) (*Applicant, error) {
	if strings.TrimSpace(cmd.TenantID) == "" {
		return nil, ErrTenantRequired
	}
	cleanEmail := strings.ToLower(strings.TrimSpace(cmd.Email))
	if !emailRegex.MatchString(cleanEmail) {
		return nil, NewDomainError("INVALID_EMAIL", "invalid email format", ErrInvalidInput)
	}

	existing, _ := s.applicantRepo.GetByEmail(ctx, cmd.TenantID, cleanEmail)
	if existing != nil {
		return nil, ErrDuplicateEmail
	}

	if cmd.Nationality == "" {
		cmd.Nationality = "Indian"
	}
	if cmd.AcademicYear == "" {
		cmd.AcademicYear = fmt.Sprintf("%d-%d", time.Now().UTC().Year(), time.Now().UTC().Year()+1)
	}
	if cmd.Type == "" {
		cmd.Type = TypeStudent
	}

	applicant := &Applicant{
		ID:                generateID("app_"),
		TenantID:          cmd.TenantID,
		Type:              cmd.Type,
		Status:            StatusDraft,
		FirstName:         strings.TrimSpace(cmd.FirstName),
		LastName:          strings.TrimSpace(cmd.LastName),
		Email:             cleanEmail,
		Phone:             strings.TrimSpace(cmd.Phone),
		DateOfBirth:       strings.TrimSpace(cmd.DateOfBirth),
		Gender:            cmd.Gender,
		BloodGroup:        strings.TrimSpace(cmd.BloodGroup),
		Nationality:       cmd.Nationality,
		Emergency:         cmd.Emergency,
		Address:           cmd.Address,
		ProgramID:         strings.TrimSpace(cmd.ProgramID),
		DepartmentID:      strings.TrimSpace(cmd.DepartmentID),
		AcademicYear:      cmd.AcademicYear,
		TargetDesignation: strings.TrimSpace(cmd.TargetDesignation),
		Guardian:          cmd.Guardian,
		Documents:         make([]Document, 0),
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	if err := applicant.ValidateDraft(); err != nil {
		return nil, err
	}

	if err := s.applicantRepo.Create(ctx, applicant); err != nil {
		return nil, err
	}

	s.emitAudit(cmd.TenantID, "", "USER", "onboarding:applicant:draft_created", "onboarding_applicant", applicant.ID, audit.StatusSuccess, nil)
	return applicant, nil
}

func (s *Service) UpdateDraft(ctx context.Context, tenantID, applicantID string, cmd UpdateApplicantCommand) (*Applicant, error) {
	applicant, err := s.applicantRepo.GetByID(ctx, tenantID, applicantID)
	if err != nil {
		return nil, err
	}
	if applicant.Status != StatusDraft {
		return nil, NewDomainError("INVALID_STATUS", "only draft applications can be edited", ErrInvalidStateTransition)
	}

	if cmd.FirstName != nil {
		applicant.FirstName = strings.TrimSpace(*cmd.FirstName)
	}
	if cmd.LastName != nil {
		applicant.LastName = strings.TrimSpace(*cmd.LastName)
	}
	if cmd.Phone != nil {
		applicant.Phone = strings.TrimSpace(*cmd.Phone)
	}
	if cmd.DateOfBirth != nil {
		applicant.DateOfBirth = strings.TrimSpace(*cmd.DateOfBirth)
	}
	if cmd.Gender != nil {
		applicant.Gender = *cmd.Gender
	}
	if cmd.BloodGroup != nil {
		applicant.BloodGroup = strings.TrimSpace(*cmd.BloodGroup)
	}
	if cmd.Nationality != nil {
		applicant.Nationality = strings.TrimSpace(*cmd.Nationality)
	}
	if cmd.Emergency != nil {
		applicant.Emergency = *cmd.Emergency
	}
	if cmd.Address != nil {
		applicant.Address = *cmd.Address
	}
	if cmd.ProgramID != nil {
		applicant.ProgramID = strings.TrimSpace(*cmd.ProgramID)
	}
	if cmd.DepartmentID != nil {
		applicant.DepartmentID = strings.TrimSpace(*cmd.DepartmentID)
	}
	if cmd.AcademicYear != nil {
		applicant.AcademicYear = strings.TrimSpace(*cmd.AcademicYear)
	}
	if cmd.TargetDesignation != nil {
		applicant.TargetDesignation = strings.TrimSpace(*cmd.TargetDesignation)
	}
	if cmd.Guardian != nil {
		applicant.Guardian = cmd.Guardian
	}

	applicant.UpdatedAt = time.Now().UTC()
	if err := s.applicantRepo.Update(ctx, applicant); err != nil {
		return nil, err
	}

	return applicant, nil
}

func (s *Service) AttachDocument(ctx context.Context, cmd AttachDocumentCommand) (*Document, error) {
	if strings.TrimSpace(cmd.TenantID) == "" {
		return nil, ErrTenantRequired
	}
	applicant, err := s.applicantRepo.GetByID(ctx, cmd.TenantID, cmd.ApplicantID)
	if err != nil {
		return nil, err
	}
	if applicant.Status != StatusDraft && applicant.Status != StatusUnderReview {
		return nil, NewDomainError("INVALID_STATUS", "documents can only be attached during draft or under-review states", ErrInvalidStateTransition)
	}

	doc := &Document{
		ID:           generateID("doc_"),
		TenantID:     cmd.TenantID,
		ApplicantID:  cmd.ApplicantID,
		DocumentType: cmd.DocumentType,
		FileKey:      cmd.FileKey,
		FileName:     cmd.FileName,
		FileSize:     cmd.FileSize,
		MIMEType:     cmd.MIMEType,
		Status:       DocPendingReview,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if err := s.docRepo.Attach(ctx, doc); err != nil {
		return nil, err
	}

	s.emitAudit(cmd.TenantID, "", "USER", "onboarding:document:attached", "onboarding_document", doc.ID, audit.StatusSuccess, map[string]interface{}{
		"applicant_id":  cmd.ApplicantID,
		"document_type": cmd.DocumentType,
	})

	return doc, nil
}

func (s *Service) SubmitApplication(ctx context.Context, tenantID, applicantID string) (*Applicant, error) {
	applicant, err := s.applicantRepo.GetByID(ctx, tenantID, applicantID)
	if err != nil {
		return nil, err
	}

	docs, err := s.docRepo.GetByApplicant(ctx, tenantID, applicantID)
	if err != nil {
		return nil, err
	}
	applicant.Documents = docs

	if err := applicant.Submit(); err != nil {
		return nil, err
	}

	if err := s.applicantRepo.Update(ctx, applicant); err != nil {
		return nil, err
	}

	s.emitAudit(tenantID, "", "USER", "onboarding:applicant:submitted", "onboarding_applicant", applicant.ID, audit.StatusSuccess, nil)
	return applicant, nil
}

func (s *Service) AssignReviewer(ctx context.Context, tenantID, applicantID, reviewerID string) (*Applicant, error) {
	applicant, err := s.applicantRepo.GetByID(ctx, tenantID, applicantID)
	if err != nil {
		return nil, err
	}

	if err := applicant.StartReview(reviewerID); err != nil {
		return nil, err
	}

	if err := s.applicantRepo.Update(ctx, applicant); err != nil {
		return nil, err
	}

	s.emitAudit(tenantID, reviewerID, "USER", "onboarding:applicant:review_started", "onboarding_applicant", applicant.ID, audit.StatusSuccess, nil)
	return applicant, nil
}

func (s *Service) VerifyDocument(ctx context.Context, tenantID, applicantID, docID, verifierID string, approved bool, rejectionReason string) (*Document, error) {
	doc, err := s.docRepo.GetByID(ctx, tenantID, docID)
	if err != nil {
		return nil, err
	}
	if doc.ApplicantID != applicantID {
		return nil, ErrDocumentNotFound
	}

	status := DocVerified
	if !approved {
		status = DocRejected
		if strings.TrimSpace(rejectionReason) == "" {
			return nil, ErrRejectionReasonRequired
		}
	}

	now := time.Now().UTC()
	if err := s.docRepo.UpdateStatus(ctx, tenantID, docID, status, rejectionReason, verifierID, &now); err != nil {
		return nil, err
	}

	doc.Status = status
	doc.RejectionReason = rejectionReason
	doc.VerifiedByID = verifierID
	doc.VerifiedAt = &now
	doc.UpdatedAt = now

	s.emitAudit(tenantID, verifierID, "USER", "onboarding:document:verified", "onboarding_document", doc.ID, audit.StatusSuccess, map[string]interface{}{
		"applicant_id": applicantID,
		"status":       status,
		"approved":     approved,
	})

	return doc, nil
}

func (s *Service) VerifyApplication(ctx context.Context, tenantID, applicantID string) (*Applicant, error) {
	applicant, err := s.applicantRepo.GetByID(ctx, tenantID, applicantID)
	if err != nil {
		return nil, err
	}

	docs, err := s.docRepo.GetByApplicant(ctx, tenantID, applicantID)
	if err != nil {
		return nil, err
	}
	applicant.Documents = docs

	if err := applicant.Verify(); err != nil {
		return nil, err
	}

	if err := s.applicantRepo.Update(ctx, applicant); err != nil {
		return nil, err
	}

	s.emitAudit(tenantID, applicant.ReviewerID, "USER", "onboarding:applicant:verified", "onboarding_applicant", applicant.ID, audit.StatusSuccess, nil)
	return applicant, nil
}

func (s *Service) RejectApplication(ctx context.Context, tenantID, applicantID, reviewerID, reason string) (*Applicant, error) {
	applicant, err := s.applicantRepo.GetByID(ctx, tenantID, applicantID)
	if err != nil {
		return nil, err
	}

	if err := applicant.Reject(reviewerID, reason); err != nil {
		return nil, err
	}

	if err := s.applicantRepo.Update(ctx, applicant); err != nil {
		return nil, err
	}

	s.emitAudit(tenantID, reviewerID, "USER", "onboarding:applicant:rejected", "onboarding_applicant", applicant.ID, audit.StatusSuccess, map[string]interface{}{
		"reason": reason,
	})
	return applicant, nil
}

func (s *Service) EnrollStudent(ctx context.Context, cmd EnrollStudentCommand) (*StudentProfile, error) {
	if strings.TrimSpace(cmd.TenantID) == "" {
		return nil, ErrTenantRequired
	}
	applicant, err := s.applicantRepo.GetByID(ctx, cmd.TenantID, cmd.ApplicantID)
	if err != nil {
		return nil, err
	}
	if applicant.Type != TypeStudent {
		return nil, NewDomainError("INVALID_TYPE", "applicant is not a student", ErrInvalidInput)
	}
	if applicant.Status != StatusVerified {
		return nil, NewDomainError("NOT_VERIFIED", "applicant KYC must be verified before enrollment", ErrInvalidStateTransition)
	}

	// Verify program and cohort
	program, err := s.academicRepo.GetProgramByID(ctx, cmd.TenantID, applicant.ProgramID)
	if err != nil {
		return nil, ErrProgramNotFound
	}
	cohort, err := s.academicRepo.GetCohortByID(ctx, cmd.TenantID, cmd.CohortID)
	if err != nil {
		return nil, ErrCohortNotFound
	}
	if cohort.CurrentEnrolled >= cohort.MaxCapacity {
		return nil, ErrCohortFull
	}

	// Generate deterministic Roll Number and Registration Number
	rollNumber, err := s.sequenceGen.GenerateRollNumber(ctx, cmd.TenantID, program.Code, applicant.AcademicYear)
	if err != nil {
		return nil, err
	}
	regNumber, err := s.sequenceGen.GenerateRegistrationNumber(ctx, cmd.TenantID, cohort.StartYear)
	if err != nil {
		return nil, err
	}

	profile := &StudentProfile{
		ID:                 generateID("stu_"),
		TenantID:           cmd.TenantID,
		UserID:             cmd.UserID,
		ApplicantID:        applicant.ID,
		RollNumber:         rollNumber,
		RegistrationNumber: regNumber,
		ProgramID:          program.ID,
		CohortID:           cohort.ID,
		CurrentSemester:    1,
		EnrollmentDate:     time.Now().UTC(),
		CreatedAt:          time.Now().UTC(),
		UpdatedAt:          time.Now().UTC(),
	}

	if err := s.profileRepo.CreateStudentProfile(ctx, profile); err != nil {
		return nil, err
	}

	_ = s.academicRepo.IncrementCohortEnrolled(ctx, cmd.TenantID, cohort.ID)

	_ = applicant.MarkEnrolled()
	_ = s.applicantRepo.Update(ctx, applicant)

	s.emitAudit(cmd.TenantID, cmd.AdminID, "USER", "onboarding:student:enrolled", "student_profile", profile.ID, audit.StatusSuccess, map[string]interface{}{
		"applicant_id":        applicant.ID,
		"roll_number":         rollNumber,
		"registration_number": regNumber,
		"cohort_id":           cohort.ID,
	})

	return profile, nil
}

func (s *Service) ProvisionStaff(ctx context.Context, cmd ProvisionStaffCommand) (*StaffProfile, error) {
	if strings.TrimSpace(cmd.TenantID) == "" {
		return nil, ErrTenantRequired
	}
	applicant, err := s.applicantRepo.GetByID(ctx, cmd.TenantID, cmd.ApplicantID)
	if err != nil {
		return nil, err
	}
	if applicant.Type == TypeStudent {
		return nil, NewDomainError("INVALID_TYPE", "applicant is a student, not staff", ErrInvalidInput)
	}
	if applicant.Status != StatusVerified {
		return nil, NewDomainError("NOT_VERIFIED", "applicant KYC must be verified before staff provisioning", ErrInvalidStateTransition)
	}

	dept, err := s.academicRepo.GetDepartmentByID(ctx, cmd.TenantID, applicant.DepartmentID)
	if err != nil {
		return nil, ErrDepartmentNotFound
	}

	employeeID, err := s.sequenceGen.GenerateEmployeeID(ctx, cmd.TenantID, dept.Code)
	if err != nil {
		return nil, err
	}

	roleKey := cmd.RoleKey
	if roleKey == "" {
		roleKey = strings.ToLower(string(applicant.Type))
	}

	profile := &StaffProfile{
		ID:           generateID("stf_"),
		TenantID:     cmd.TenantID,
		UserID:       cmd.UserID,
		ApplicantID:  applicant.ID,
		EmployeeID:   employeeID,
		DepartmentID: dept.ID,
		Designation:  applicant.TargetDesignation,
		RoleKey:      roleKey,
		JoiningDate:  time.Now().UTC(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if err := s.profileRepo.CreateStaffProfile(ctx, profile); err != nil {
		return nil, err
	}

	_ = applicant.MarkEnrolled()
	_ = s.applicantRepo.Update(ctx, applicant)

	s.emitAudit(cmd.TenantID, cmd.AdminID, "USER", "onboarding:staff:provisioned", "staff_profile", profile.ID, audit.StatusSuccess, map[string]interface{}{
		"applicant_id": applicant.ID,
		"employee_id":  employeeID,
		"department":   dept.Code,
		"role_key":     roleKey,
	})

	return profile, nil
}

func (s *Service) GetApplicant(ctx context.Context, tenantID, applicantID string) (*Applicant, error) {
	applicant, err := s.applicantRepo.GetByID(ctx, tenantID, applicantID)
	if err != nil {
		return nil, err
	}
	docs, _ := s.docRepo.GetByApplicant(ctx, tenantID, applicantID)
	applicant.Documents = docs
	return applicant, nil
}

func (s *Service) ListApplicants(ctx context.Context, filter ApplicantFilter) ([]*Applicant, int, error) {
	return s.applicantRepo.List(ctx, filter)
}

func (s *Service) ListDepartments(ctx context.Context, tenantID string) ([]*AcademicDepartment, error) {
	return s.academicRepo.ListDepartments(ctx, tenantID)
}

func (s *Service) ListPrograms(ctx context.Context, tenantID, departmentID string) ([]*AcademicProgram, error) {
	return s.academicRepo.ListPrograms(ctx, tenantID, departmentID)
}

func (s *Service) ListCohorts(ctx context.Context, tenantID, programID, academicYear string) ([]*CohortBatch, error) {
	return s.academicRepo.ListCohorts(ctx, tenantID, programID, academicYear)
}

func (s *Service) emitAudit(tenantID, actorID, actorType, action, resType, resID string, status audit.Status, meta map[string]interface{}) {
	if s.auditSub == nil {
		return
	}
	var metaRaw json.RawMessage
	if meta != nil {
		if b, err := json.Marshal(meta); err == nil {
			metaRaw = b
		}
	}
	_ = s.auditSub.Enqueue(audit.RecordAuditRequest{
		TenantID:     tenantID,
		ActorID:      &actorID,
		ActorType:    audit.ActorType(actorType),
		Action:       action,
		ResourceType: resType,
		ResourceID:   &resID,
		Status:       status,
		Metadata:     metaRaw,
	})
}

func generateID(prefix string) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return prefix + hex.EncodeToString(b)
}
