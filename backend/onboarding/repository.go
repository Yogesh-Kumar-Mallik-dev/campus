/**
 * BLOCK_ONBOARDING_REPOSITORY_001
 * Subsystem: Rank 3 - Student & Staff Registration System (onboarding)
 * Purpose:   Data access interface definitions for applicants, documents, academic units, and sequences.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package onboarding

import (
	"context"
	"time"
)

type ApplicantFilter struct {
	TenantID     string
	Type         *OnboardingType
	Status       *OnboardingStatus
	ProgramID    string
	DepartmentID string
	AcademicYear string
	Search       string
	Limit        int
	Offset       int
}

type ApplicantRepository interface {
	Create(ctx context.Context, applicant *Applicant) error
	GetByID(ctx context.Context, tenantID, id string) (*Applicant, error)
	GetByEmail(ctx context.Context, tenantID, email string) (*Applicant, error)
	Update(ctx context.Context, applicant *Applicant) error
	List(ctx context.Context, filter ApplicantFilter) ([]*Applicant, int, error)
}

type DocumentRepository interface {
	Attach(ctx context.Context, doc *Document) error
	GetByID(ctx context.Context, tenantID, id string) (*Document, error)
	GetByApplicant(ctx context.Context, tenantID, applicantID string) ([]Document, error)
	UpdateStatus(ctx context.Context, tenantID, docID string, status DocumentStatus, reason, verifierID string, verifiedAt *time.Time) error
}

type AcademicRepository interface {
	CreateDepartment(ctx context.Context, dept *AcademicDepartment) error
	GetDepartmentByID(ctx context.Context, tenantID, id string) (*AcademicDepartment, error)
	GetDepartmentByCode(ctx context.Context, tenantID, code string) (*AcademicDepartment, error)
	ListDepartments(ctx context.Context, tenantID string) ([]*AcademicDepartment, error)

	CreateProgram(ctx context.Context, program *AcademicProgram) error
	GetProgramByID(ctx context.Context, tenantID, id string) (*AcademicProgram, error)
	GetProgramByCode(ctx context.Context, tenantID, code string) (*AcademicProgram, error)
	ListPrograms(ctx context.Context, tenantID string, departmentID string) ([]*AcademicProgram, error)

	CreateCohort(ctx context.Context, cohort *CohortBatch) error
	GetCohortByID(ctx context.Context, tenantID, id string) (*CohortBatch, error)
	ListCohorts(ctx context.Context, tenantID, programID, academicYear string) ([]*CohortBatch, error)
	IncrementCohortEnrolled(ctx context.Context, tenantID, cohortID string) error
}

type ProfileRepository interface {
	CreateStudentProfile(ctx context.Context, profile *StudentProfile) error
	GetStudentProfileByRollNumber(ctx context.Context, tenantID, rollNumber string) (*StudentProfile, error)
	GetStudentProfileByApplicantID(ctx context.Context, tenantID, applicantID string) (*StudentProfile, error)

	CreateStaffProfile(ctx context.Context, profile *StaffProfile) error
	GetStaffProfileByEmployeeID(ctx context.Context, tenantID, employeeID string) (*StaffProfile, error)
	GetStaffProfileByApplicantID(ctx context.Context, tenantID, applicantID string) (*StaffProfile, error)
}

type SequenceRepository interface {
	NextSequence(ctx context.Context, tenantID, sequenceType, prefix string) (int, error)
}
