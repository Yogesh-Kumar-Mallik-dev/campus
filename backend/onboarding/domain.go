/**
 * BLOCK_ONBOARDING_DOMAIN_001
 * Subsystem: Rank 3 - Student & Staff Registration System (onboarding)
 * Purpose:   Domain entities, state machine transitions, KYC validation, and cohort aggregates.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package onboarding

import (
	"regexp"
	"strings"
	"time"
)

type OnboardingType string

const (
	TypeStudent   OnboardingType = "STUDENT"
	TypeFaculty   OnboardingType = "FACULTY"
	TypeStaff     OnboardingType = "STAFF"
	TypeWarden    OnboardingType = "WARDEN"
	TypeLibrarian OnboardingType = "LIBRARIAN"
)

type OnboardingStatus string

const (
	StatusDraft       OnboardingStatus = "DRAFT"
	StatusSubmitted   OnboardingStatus = "SUBMITTED"
	StatusUnderReview OnboardingStatus = "UNDER_REVIEW"
	StatusVerified    OnboardingStatus = "VERIFIED"
	StatusRejected    OnboardingStatus = "REJECTED"
	StatusEnrolled    OnboardingStatus = "ENROLLED"
)

type DocumentType string

const (
	DocNationalID        DocumentType = "NATIONAL_ID"
	DocPassport          DocumentType = "PASSPORT"
	DocBirthCert         DocumentType = "BIRTH_CERTIFICATE"
	DocTranscript        DocumentType = "ACADEMIC_TRANSCRIPT"
	DocDegreeCert        DocumentType = "DEGREE_CERTIFICATE"
	DocTransferCert      DocumentType = "TRANSFER_CERTIFICATE"
	DocMedicalFitness    DocumentType = "MEDICAL_FITNESS"
	DocPassportPhoto     DocumentType = "PASSPORT_PHOTO"
)

type DocumentStatus string

const (
	DocPendingReview DocumentStatus = "PENDING_REVIEW"
	DocVerified      DocumentStatus = "VERIFIED"
	DocRejected      DocumentStatus = "REJECTED"
)

type Gender string

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
	GenderOther  Gender = "OTHER"
)

type DegreeType string

const (
	DegreeUG       DegreeType = "UG"
	DegreePG       DegreeType = "PG"
	DegreeDiploma  DegreeType = "DIPLOMA"
	DegreeDoctoral DegreeType = "DOCTORAL"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type EmergencyContact struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Relation string `json:"relation"`
}

type Address struct {
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

type Guardian struct {
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone"`
	Relation string `json:"relation"`
}

type Document struct {
	ID              string         `json:"id"`
	TenantID        string         `json:"tenant_id"`
	ApplicantID     string         `json:"applicant_id"`
	DocumentType    DocumentType   `json:"document_type"`
	FileKey         string         `json:"file_key"`
	FileName        string         `json:"file_name"`
	FileSize        int64          `json:"file_size"`
	MIMEType        string         `json:"mime_type"`
	Status          DocumentStatus `json:"status"`
	RejectionReason string         `json:"rejection_reason,omitempty"`
	VerifiedByID    string         `json:"verified_by_id,omitempty"`
	VerifiedAt      *time.Time     `json:"verified_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type Applicant struct {
	ID                string           `json:"id"`
	TenantID          string           `json:"tenant_id"`
	Type              OnboardingType   `json:"type"`
	Status            OnboardingStatus `json:"status"`
	FirstName         string           `json:"first_name"`
	LastName          string           `json:"last_name"`
	Email             string           `json:"email"`
	Phone             string           `json:"phone"`
	DateOfBirth       string           `json:"date_of_birth"` // YYYY-MM-DD
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
	ReviewerID        string           `json:"reviewer_id,omitempty"`
	ReviewedAt        *time.Time       `json:"reviewed_at,omitempty"`
	RejectionReason   string           `json:"rejection_reason,omitempty"`
	Documents         []Document       `json:"documents"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

type AcademicDepartment struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AcademicProgram struct {
	ID             string     `json:"id"`
	TenantID       string     `json:"tenant_id"`
	DepartmentID   string     `json:"department_id"`
	Code           string     `json:"code"`
	Name           string     `json:"name"`
	DegreeType     DegreeType `json:"degree_type"`
	DurationYears  int        `json:"duration_years"`
	TotalSemesters int        `json:"total_semesters"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CohortBatch struct {
	ID              string    `json:"id"`
	TenantID        string    `json:"tenant_id"`
	ProgramID       string    `json:"program_id"`
	AcademicYear    string    `json:"academic_year"`
	StartYear       int       `json:"start_year"`
	EndYear         int       `json:"end_year"`
	Section         string    `json:"section"`
	MaxCapacity     int       `json:"max_capacity"`
	CurrentEnrolled int       `json:"current_enrolled"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type StudentProfile struct {
	ID                 string    `json:"id"`
	TenantID           string    `json:"tenant_id"`
	UserID             string    `json:"user_id"`
	ApplicantID        string    `json:"applicant_id"`
	RollNumber         string    `json:"roll_number"`
	RegistrationNumber string    `json:"registration_number"`
	ProgramID          string    `json:"program_id"`
	CohortID           string    `json:"cohort_id"`
	CurrentSemester    int       `json:"current_semester"`
	EnrollmentDate     time.Time `json:"enrollment_date"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type StaffProfile struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	UserID       string    `json:"user_id"`
	ApplicantID  string    `json:"applicant_id"`
	EmployeeID   string    `json:"employee_id"`
	DepartmentID string    `json:"department_id"`
	Designation  string    `json:"designation"`
	RoleKey      string    `json:"role_key"`
	JoiningDate  time.Time `json:"joining_date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// MandatoryDocuments returns required documents according to applicant type.
func (a *Applicant) MandatoryDocuments() []DocumentType {
	if a.Type == TypeStudent {
		return []DocumentType{DocNationalID, DocTranscript, DocPassportPhoto}
	}
	// Staff / Faculty
	return []DocumentType{DocNationalID, DocDegreeCert, DocPassportPhoto}
}

// ValidateDraft ensures minimum fields for draft creation.
func (a *Applicant) ValidateDraft() error {
	if strings.TrimSpace(a.TenantID) == "" {
		return ErrTenantRequired
	}
	if strings.TrimSpace(a.FirstName) == "" || strings.TrimSpace(a.LastName) == "" {
		return NewDomainError("INVALID_NAME", "first and last name are required", ErrInvalidInput)
	}
	if !emailRegex.MatchString(a.Email) {
		return NewDomainError("INVALID_EMAIL", "a valid email is required", ErrInvalidInput)
	}
	return nil
}

// ValidateSubmission ensures all required KYC fields & mandatory documents are attached.
func (a *Applicant) ValidateSubmission() error {
	if err := a.ValidateDraft(); err != nil {
		return err
	}
	if strings.TrimSpace(a.Phone) == "" {
		return NewDomainError("INVALID_PHONE", "contact phone number is required", ErrInvalidInput)
	}
	if strings.TrimSpace(a.DateOfBirth) == "" {
		return NewDomainError("INVALID_DOB", "date of birth is required", ErrInvalidInput)
	}
	if strings.TrimSpace(a.Emergency.Name) == "" || strings.TrimSpace(a.Emergency.Phone) == "" {
		return NewDomainError("INVALID_EMERGENCY", "emergency contact name and phone are required", ErrInvalidInput)
	}
	if strings.TrimSpace(a.Address.Line1) == "" || strings.TrimSpace(a.Address.City) == "" || strings.TrimSpace(a.Address.PostalCode) == "" {
		return NewDomainError("INVALID_ADDRESS", "address line 1, city, and postal code are required", ErrInvalidInput)
	}

	if a.Type == TypeStudent {
		if strings.TrimSpace(a.ProgramID) == "" {
			return NewDomainError("INVALID_PROGRAM", "program ID is required for student registration", ErrInvalidInput)
		}
		if a.Guardian == nil || strings.TrimSpace(a.Guardian.Name) == "" || strings.TrimSpace(a.Guardian.Phone) == "" {
			return NewDomainError("INVALID_GUARDIAN", "guardian details are required for student registration", ErrInvalidInput)
		}
	} else {
		if strings.TrimSpace(a.DepartmentID) == "" {
			return NewDomainError("INVALID_DEPARTMENT", "department ID is required for staff registration", ErrInvalidInput)
		}
		if strings.TrimSpace(a.TargetDesignation) == "" {
			return NewDomainError("INVALID_DESIGNATION", "target designation is required for staff registration", ErrInvalidInput)
		}
	}

	// Check mandatory documents attached
	mand := a.MandatoryDocuments()
	attached := make(map[DocumentType]bool)
	for _, doc := range a.Documents {
		attached[doc.DocumentType] = true
	}
	for _, m := range mand {
		if !attached[m] {
			return ErrMissingMandatoryDocs
		}
	}

	return nil
}

// Submit transitions DRAFT -> SUBMITTED.
func (a *Applicant) Submit() error {
	if a.Status != StatusDraft {
		return NewDomainError("INVALID_STATUS", "only draft applications can be submitted", ErrInvalidStateTransition)
	}
	if err := a.ValidateSubmission(); err != nil {
		return err
	}
	a.Status = StatusSubmitted
	a.UpdatedAt = time.Now().UTC()
	return nil
}

// StartReview transitions SUBMITTED -> UNDER_REVIEW.
func (a *Applicant) StartReview(reviewerID string) error {
	if strings.TrimSpace(reviewerID) == "" {
		return ErrUnauthorizedReviewer
	}
	if a.Status != StatusSubmitted && a.Status != StatusUnderReview {
		return NewDomainError("INVALID_STATUS", "only submitted applications can be placed under review", ErrInvalidStateTransition)
	}
	a.Status = StatusUnderReview
	a.ReviewerID = reviewerID
	now := time.Now().UTC()
	a.ReviewedAt = &now
	a.UpdatedAt = now
	return nil
}

// Verify transitions UNDER_REVIEW -> VERIFIED once all mandatory docs are verified.
func (a *Applicant) Verify() error {
	if a.Status != StatusUnderReview {
		return NewDomainError("INVALID_STATUS", "application must be under review before verification", ErrInvalidStateTransition)
	}

	mand := a.MandatoryDocuments()
	verifiedMap := make(map[DocumentType]bool)
	for _, doc := range a.Documents {
		if doc.Status == DocRejected {
			return ErrUnverifiedDocuments
		}
		if doc.Status == DocVerified {
			verifiedMap[doc.DocumentType] = true
		}
	}

	for _, m := range mand {
		if !verifiedMap[m] {
			return ErrMissingMandatoryDocs
		}
	}

	a.Status = StatusVerified
	a.UpdatedAt = time.Now().UTC()
	return nil
}

// Reject transitions to REJECTED with a mandatory reason.
func (a *Applicant) Reject(reviewerID, reason string) error {
	if strings.TrimSpace(reason) == "" {
		return ErrRejectionReasonRequired
	}
	if a.Status == StatusEnrolled {
		return NewDomainError("ALREADY_ENROLLED", "cannot reject an already enrolled applicant", ErrInvalidStateTransition)
	}
	a.Status = StatusRejected
	a.ReviewerID = reviewerID
	a.RejectionReason = reason
	now := time.Now().UTC()
	a.ReviewedAt = &now
	a.UpdatedAt = now
	return nil
}

// MarkEnrolled marks applicant as ENROLLED.
func (a *Applicant) MarkEnrolled() error {
	if a.Status != StatusVerified {
		return NewDomainError("NOT_VERIFIED", "applicant must be verified before enrollment", ErrInvalidStateTransition)
	}
	a.Status = StatusEnrolled
	a.UpdatedAt = time.Now().UTC()
	return nil
}
