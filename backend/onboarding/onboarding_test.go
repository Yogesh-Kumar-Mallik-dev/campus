/**
 * BLOCK_ONBOARDING_TEST_001
 * Subsystem: Rank 3 - Student & Staff Registration System (onboarding)
 * Purpose:   Exhaustive unit test suite verifying KYC validation, document verifications, sequence generation, state transitions, and audit integration.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package onboarding

import (
	"context"
	"sync"
	"testing"
	"time"

	"campus/backend/audit"
)

type mockAuditSubscriber struct {
	mu     sync.Mutex
	events []audit.RecordAuditRequest
}

func (m *mockAuditSubscriber) Enqueue(req audit.RecordAuditRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, req)
	return nil
}

func (m *mockAuditSubscriber) Start(ctx context.Context) {}
func (m *mockAuditSubscriber) Stop()                      {}

func setupTestEnvironment() (*Service, *MockAcademicRepository, *mockAuditSubscriber) {
	appRepo := NewMockApplicantRepository()
	docRepo := NewMockDocumentRepository()
	acadRepo := NewMockAcademicRepository()
	profRepo := NewMockProfileRepository()
	seqRepo := NewMockSequenceRepository()
	seqEngine := NewSequenceEngine(seqRepo)
	auditSub := &mockAuditSubscriber{}

	svc := NewService(appRepo, docRepo, acadRepo, profRepo, seqEngine, auditSub)
	return svc, acadRepo, auditSub
}

func TestApplicant_DraftCreationAndValidation(t *testing.T) {
	ctx := context.Background()
	svc, _, _ := setupTestEnvironment()

	// 1. Missing tenant ID
	_, err := svc.CreateDraft(ctx, CreateApplicantCommand{
		Email: "student@campus.edu",
	})
	if err != ErrTenantRequired {
		t.Fatalf("expected ErrTenantRequired, got %v", err)
	}

	// 2. Invalid Email
	_, err = svc.CreateDraft(ctx, CreateApplicantCommand{
		TenantID:  "tenant-alpha",
		Email:     "invalid-email",
		FirstName: "Aarav",
		LastName:  "Sharma",
	})
	if err == nil {
		t.Fatalf("expected error for invalid email")
	}

	// 3. Nominal Draft Creation
	cmd := CreateApplicantCommand{
		TenantID:    "tenant-alpha",
		Type:        TypeStudent,
		FirstName:   "Aarav",
		LastName:    "Sharma",
		Email:       "aarav.sharma@example.com",
		Phone:       "+919876543210",
		DateOfBirth: "2005-08-15",
		Gender:      GenderMale,
		Nationality: "Indian",
		Emergency: EmergencyContact{
			Name:     "Rajesh Sharma",
			Phone:    "+919876543211",
			Relation: "Father",
		},
		Address: Address{
			Line1:      "42 Campus Road",
			City:       "Bengaluru",
			State:      "Karnataka",
			PostalCode: "560001",
			Country:    "India",
		},
		ProgramID:    "prog-cse",
		AcademicYear: "2026-2027",
		Guardian: &Guardian{
			Name:     "Rajesh Sharma",
			Phone:    "+919876543211",
			Relation: "Father",
		},
	}

	app, err := svc.CreateDraft(ctx, cmd)
	if err != nil {
		t.Fatalf("unexpected error creating draft: %v", err)
	}
	if app.ID == "" || app.Status != StatusDraft {
		t.Fatalf("expected valid draft applicant, got %+v", app)
	}

	// 4. Duplicate Email
	_, err = svc.CreateDraft(ctx, cmd)
	if err != ErrDuplicateEmail {
		t.Fatalf("expected ErrDuplicateEmail, got %v", err)
	}
}

func TestApplicant_DraftUpdate(t *testing.T) {
	ctx := context.Background()
	svc, _, _ := setupTestEnvironment()

	app, err := svc.CreateDraft(ctx, CreateApplicantCommand{
		TenantID:  "tenant-alpha",
		Type:      TypeStudent,
		FirstName: "Priya",
		LastName:  "Nair",
		Email:     "priya.nair@example.com",
	})
	if err != nil {
		t.Fatalf("failed to create draft: %v", err)
	}

	newPhone := "+919876500000"
	newBloodGroup := "O+"
	updated, err := svc.UpdateDraft(ctx, "tenant-alpha", app.ID, UpdateApplicantCommand{
		Phone:      &newPhone,
		BloodGroup: &newBloodGroup,
	})
	if err != nil {
		t.Fatalf("failed to update draft: %v", err)
	}
	if updated.Phone != newPhone || updated.BloodGroup != "O+" {
		t.Fatalf("update not applied: phone=%s, bloodGroup=%s", updated.Phone, updated.BloodGroup)
	}
}

func TestApplicant_FullStudentLifecycle_Enrollment(t *testing.T) {
	ctx := context.Background()
	svc, acadRepo, auditSub := setupTestEnvironment()

	tenantID := "tenant-alpha"

	// Seed Department, Program, Cohort
	dept := &AcademicDepartment{
		ID:        "dept-cse",
		TenantID:  tenantID,
		Code:      "CSE",
		Name:      "Computer Science & Engineering",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_ = acadRepo.CreateDepartment(ctx, dept)

	prog := &AcademicProgram{
		ID:             "prog-btech-cse",
		TenantID:       tenantID,
		DepartmentID:   dept.ID,
		Code:           "BTECH_CSE",
		Name:           "B.Tech in Computer Science",
		DegreeType:     DegreeUG,
		DurationYears:  4,
		TotalSemesters: 8,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	_ = acadRepo.CreateProgram(ctx, prog)

	cohort := &CohortBatch{
		ID:              "cohort-2026-cse-a",
		TenantID:        tenantID,
		ProgramID:       prog.ID,
		AcademicYear:    "2026-2027",
		StartYear:       2026,
		EndYear:         2030,
		Section:         "A",
		MaxCapacity:     60,
		CurrentEnrolled: 0,
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	_ = acadRepo.CreateCohort(ctx, cohort)

	// Step 1: Create Draft
	app, err := svc.CreateDraft(ctx, CreateApplicantCommand{
		TenantID:    tenantID,
		Type:        TypeStudent,
		FirstName:   "Rohan",
		LastName:    "Verma",
		Email:       "rohan.verma@example.com",
		Phone:       "+919876543222",
		DateOfBirth: "2006-03-20",
		Gender:      GenderMale,
		Nationality: "Indian",
		Emergency: EmergencyContact{
			Name:     "Suresh Verma",
			Phone:    "+919876543223",
			Relation: "Father",
		},
		Address: Address{
			Line1:      "12 Tech Enclave",
			City:       "Hyderabad",
			State:      "Telangana",
			PostalCode: "500081",
			Country:    "India",
		},
		ProgramID:    prog.ID,
		AcademicYear: "2026-2027",
		Guardian: &Guardian{
			Name:     "Suresh Verma",
			Phone:    "+919876543223",
			Relation: "Father",
		},
	})
	if err != nil {
		t.Fatalf("failed to create draft: %v", err)
	}

	// Step 2: Try submitting without documents -> Must fail
	_, err = svc.SubmitApplication(ctx, tenantID, app.ID)
	if err != ErrMissingMandatoryDocs {
		t.Fatalf("expected ErrMissingMandatoryDocs, got %v", err)
	}

	// Step 3: Attach Mandatory Documents (National ID, Transcript, Photo)
	docID1, err := svc.AttachDocument(ctx, AttachDocumentCommand{
		TenantID:     tenantID,
		ApplicantID:  app.ID,
		DocumentType: DocNationalID,
		FileKey:      "uploads/aadhaar.pdf",
		FileName:     "aadhaar.pdf",
		FileSize:     204800,
		MIMEType:     "application/pdf",
	})
	if err != nil {
		t.Fatalf("failed to attach national ID: %v", err)
	}

	docID2, err := svc.AttachDocument(ctx, AttachDocumentCommand{
		TenantID:     tenantID,
		ApplicantID:  app.ID,
		DocumentType: DocTranscript,
		FileKey:      "uploads/marksheet.pdf",
		FileName:     "marksheet.pdf",
		FileSize:     512000,
		MIMEType:     "application/pdf",
	})
	if err != nil {
		t.Fatalf("failed to attach transcript: %v", err)
	}

	docID3, err := svc.AttachDocument(ctx, AttachDocumentCommand{
		TenantID:     tenantID,
		ApplicantID:  app.ID,
		DocumentType: DocPassportPhoto,
		FileKey:      "uploads/photo.jpg",
		FileName:     "photo.jpg",
		FileSize:     102400,
		MIMEType:     "image/jpeg",
	})
	if err != nil {
		t.Fatalf("failed to attach photo: %v", err)
	}

	// Step 4: Submit Application -> Now succeeds
	app, err = svc.SubmitApplication(ctx, tenantID, app.ID)
	if err != nil {
		t.Fatalf("failed to submit application: %v", err)
	}
	if app.Status != StatusSubmitted {
		t.Fatalf("expected status SUBMITTED, got %s", app.Status)
	}

	// Step 5: Assign Reviewer (moves to UNDER_REVIEW)
	reviewerID := "usr_admissions_admin"
	app, err = svc.AssignReviewer(ctx, tenantID, app.ID, reviewerID)
	if err != nil {
		t.Fatalf("failed to assign reviewer: %v", err)
	}
	if app.Status != StatusUnderReview || app.ReviewerID != reviewerID {
		t.Fatalf("expected status UNDER_REVIEW with reviewer %s, got %+v", reviewerID, app)
	}

	// Step 6: Try verifying before documents are verified -> Must fail
	_, err = svc.VerifyApplication(ctx, tenantID, app.ID)
	if err != ErrMissingMandatoryDocs {
		t.Fatalf("expected ErrMissingMandatoryDocs on unverified docs, got %v", err)
	}

	// Step 7: Verify all 3 documents
	_, err = svc.VerifyDocument(ctx, tenantID, app.ID, docID1.ID, reviewerID, true, "")
	if err != nil {
		t.Fatalf("failed to verify doc 1: %v", err)
	}
	_, err = svc.VerifyDocument(ctx, tenantID, app.ID, docID2.ID, reviewerID, true, "")
	if err != nil {
		t.Fatalf("failed to verify doc 2: %v", err)
	}
	_, err = svc.VerifyDocument(ctx, tenantID, app.ID, docID3.ID, reviewerID, true, "")
	if err != nil {
		t.Fatalf("failed to verify doc 3: %v", err)
	}

	// Step 8: Verify Application -> Now succeeds
	app, err = svc.VerifyApplication(ctx, tenantID, app.ID)
	if err != nil {
		t.Fatalf("failed to verify application: %v", err)
	}
	if app.Status != StatusVerified {
		t.Fatalf("expected status VERIFIED, got %s", app.Status)
	}

	// Step 9: Enroll Student -> Generates deterministic Roll Number
	studentUser := "usr_student_rohan"
	studentProfile, err := svc.EnrollStudent(ctx, EnrollStudentCommand{
		TenantID:    tenantID,
		ApplicantID: app.ID,
		UserID:      studentUser,
		CohortID:    cohort.ID,
		AdminID:     reviewerID,
	})
	if err != nil {
		t.Fatalf("failed to enroll student: %v", err)
	}

	if studentProfile.RollNumber != "2026-BTECH_CSE-0001" {
		t.Fatalf("expected roll number '2026-BTECH_CSE-0001', got '%s'", studentProfile.RollNumber)
	}
	if studentProfile.CurrentSemester != 1 {
		t.Fatalf("expected semester 1, got %d", studentProfile.CurrentSemester)
	}

	// Verify applicant marked ENROLLED
	finalApp, err := svc.GetApplicant(ctx, tenantID, app.ID)
	if err != nil {
		t.Fatalf("failed to fetch final applicant: %v", err)
	}
	if finalApp.Status != StatusEnrolled {
		t.Fatalf("expected status ENROLLED, got %s", finalApp.Status)
	}

	// Verify Audit Events emitted
	auditSub.mu.Lock()
	defer auditSub.mu.Unlock()
	if len(auditSub.events) < 5 {
		t.Fatalf("expected at least 5 audit events, got %d", len(auditSub.events))
	}
}

func TestApplicant_FullStaffLifecycle_Provisioning(t *testing.T) {
	ctx := context.Background()
	svc, acadRepo, _ := setupTestEnvironment()

	tenantID := "tenant-alpha"

	// Seed Department
	dept := &AcademicDepartment{
		ID:        "dept-ece",
		TenantID:  tenantID,
		Code:      "ECE",
		Name:      "Electronics & Communication Engineering",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_ = acadRepo.CreateDepartment(ctx, dept)

	// Step 1: Create Faculty Draft
	app, err := svc.CreateDraft(ctx, CreateApplicantCommand{
		TenantID:          tenantID,
		Type:              TypeFaculty,
		FirstName:         "Dr. Ananya",
		LastName:          "Iyer",
		Email:             "ananya.iyer@example.com",
		Phone:             "+919876543333",
		DateOfBirth:       "1985-11-22",
		Gender:            GenderFemale,
		Nationality:       "Indian",
		DepartmentID:      dept.ID,
		AcademicYear:      "2026-2027",
		TargetDesignation: "Associate Professor",
		Emergency: EmergencyContact{
			Name:     "Karthik Iyer",
			Phone:    "+919876543334",
			Relation: "Spouse",
		},
		Address: Address{
			Line1:      "88 Faculty Quarters",
			City:       "Chennai",
			State:      "Tamil Nadu",
			PostalCode: "600025",
			Country:    "India",
		},
	})
	if err != nil {
		t.Fatalf("failed to create faculty draft: %v", err)
	}

	// Step 2: Attach Mandatory Staff Documents (National ID, Degree Certificate, Passport Photo)
	doc1, err := svc.AttachDocument(ctx, AttachDocumentCommand{
		TenantID:     tenantID,
		ApplicantID:  app.ID,
		DocumentType: DocNationalID,
		FileKey:      "uploads/passport.pdf",
		FileName:     "passport.pdf",
		FileSize:     102400,
		MIMEType:     "application/pdf",
	})
	if err != nil {
		t.Fatalf("failed to attach national ID: %v", err)
	}

	doc2, err := svc.AttachDocument(ctx, AttachDocumentCommand{
		TenantID:     tenantID,
		ApplicantID:  app.ID,
		DocumentType: DocDegreeCert,
		FileKey:      "uploads/phd_degree.pdf",
		FileName:     "phd_degree.pdf",
		FileSize:     300000,
		MIMEType:     "application/pdf",
	})
	if err != nil {
		t.Fatalf("failed to attach degree: %v", err)
	}

	doc3, err := svc.AttachDocument(ctx, AttachDocumentCommand{
		TenantID:     tenantID,
		ApplicantID:  app.ID,
		DocumentType: DocPassportPhoto,
		FileKey:      "uploads/ananya.jpg",
		FileName:     "ananya.jpg",
		FileSize:     80000,
		MIMEType:     "image/jpeg",
	})
	if err != nil {
		t.Fatalf("failed to attach photo: %v", err)
	}

	// Step 3: Submit -> Under Review -> Verify Docs -> Verify Application
	app, err = svc.SubmitApplication(ctx, tenantID, app.ID)
	if err != nil {
		t.Fatalf("failed to submit faculty application: %v", err)
	}

	adminID := "usr_hr_lead"
	app, err = svc.AssignReviewer(ctx, tenantID, app.ID, adminID)
	if err != nil {
		t.Fatalf("failed to assign reviewer: %v", err)
	}

	_, _ = svc.VerifyDocument(ctx, tenantID, app.ID, doc1.ID, adminID, true, "")
	_, _ = svc.VerifyDocument(ctx, tenantID, app.ID, doc2.ID, adminID, true, "")
	_, _ = svc.VerifyDocument(ctx, tenantID, app.ID, doc3.ID, adminID, true, "")

	app, err = svc.VerifyApplication(ctx, tenantID, app.ID)
	if err != nil {
		t.Fatalf("failed to verify faculty application: %v", err)
	}

	// Step 4: Provision Staff -> Generates EMP-ECE-0001
	staffProfile, err := svc.ProvisionStaff(ctx, ProvisionStaffCommand{
		TenantID:    tenantID,
		ApplicantID: app.ID,
		UserID:      "usr_faculty_ananya",
		RoleKey:     "faculty",
		AdminID:     adminID,
	})
	if err != nil {
		t.Fatalf("failed to provision staff: %v", err)
	}

	if staffProfile.EmployeeID != "EMP-ECE-0001" {
		t.Fatalf("expected employee ID 'EMP-ECE-0001', got '%s'", staffProfile.EmployeeID)
	}
	if staffProfile.Designation != "Associate Professor" {
		t.Fatalf("expected designation 'Associate Professor', got '%s'", staffProfile.Designation)
	}
}

func TestApplicant_DocumentRejection_And_AppRejection(t *testing.T) {
	ctx := context.Background()
	svc, _, _ := setupTestEnvironment()

	tenantID := "tenant-beta"

	app, err := svc.CreateDraft(ctx, CreateApplicantCommand{
		TenantID:    tenantID,
		Type:        TypeStudent,
		FirstName:   "Karan",
		LastName:    "Singh",
		Email:       "karan.singh@example.com",
		Phone:       "+919876543444",
		DateOfBirth: "2005-01-10",
		Gender:      GenderMale,
		Nationality: "Indian",
		Emergency: EmergencyContact{
			Name:     "Mohan Singh",
			Phone:    "+919876543445",
			Relation: "Father",
		},
		Address: Address{
			Line1:      "45 Lake View",
			City:       "Pune",
			State:      "Maharashtra",
			PostalCode: "411001",
			Country:    "India",
		},
		ProgramID:    "prog-mech",
		AcademicYear: "2026-2027",
		Guardian: &Guardian{
			Name:     "Mohan Singh",
			Phone:    "+919876543445",
			Relation: "Father",
		},
	})
	if err != nil {
		t.Fatalf("failed to create draft: %v", err)
	}

	doc, err := svc.AttachDocument(ctx, AttachDocumentCommand{
		TenantID:     tenantID,
		ApplicantID:  app.ID,
		DocumentType: DocNationalID,
		FileKey:      "uploads/blurry_id.pdf",
		FileName:     "blurry_id.pdf",
		FileSize:     50000,
		MIMEType:     "application/pdf",
	})
	if err != nil {
		t.Fatalf("failed to attach doc: %v", err)
	}

	_, _ = svc.AttachDocument(ctx, AttachDocumentCommand{
		TenantID:     tenantID,
		ApplicantID:  app.ID,
		DocumentType: DocTranscript,
		FileKey:      "uploads/transcript.pdf",
		FileName:     "transcript.pdf",
		FileSize:     50000,
		MIMEType:     "application/pdf",
	})

	_, _ = svc.AttachDocument(ctx, AttachDocumentCommand{
		TenantID:     tenantID,
		ApplicantID:  app.ID,
		DocumentType: DocPassportPhoto,
		FileKey:      "uploads/photo.jpg",
		FileName:     "photo.jpg",
		FileSize:     50000,
		MIMEType:     "image/jpeg",
	})

	app, err = svc.SubmitApplication(ctx, tenantID, app.ID)
	if err != nil {
		t.Fatalf("failed to submit: %v", err)
	}

	reviewerID := "usr_admissions_officer"
	app, _ = svc.AssignReviewer(ctx, tenantID, app.ID, reviewerID)

	// Reject document with reason
	_, err = svc.VerifyDocument(ctx, tenantID, app.ID, doc.ID, reviewerID, false, "Document is unreadable/blurry")
	if err != nil {
		t.Fatalf("failed to reject doc: %v", err)
	}

	// Verify application must fail with unverified/rejected documents
	_, err = svc.VerifyApplication(ctx, tenantID, app.ID)
	if err != ErrUnverifiedDocuments {
		t.Fatalf("expected ErrUnverifiedDocuments, got %v", err)
	}

	// Reject applicant
	app, err = svc.RejectApplication(ctx, tenantID, app.ID, reviewerID, "Document verification failed repeatedly")
	if err != nil {
		t.Fatalf("failed to reject application: %v", err)
	}
	if app.Status != StatusRejected || app.RejectionReason == "" {
		t.Fatalf("expected status REJECTED with reason, got %+v", app)
	}
}

func TestSequenceEngine_ConcurrentSequenceGeneration(t *testing.T) {
	seqRepo := NewMockSequenceRepository()
	engine := NewSequenceEngine(seqRepo)

	ctx := context.Background()
	tenantID := "tenant-concurrency"
	programCode := "BTECH_AI"
	academicYear := "2026-2027"

	const concurrency = 50
	var wg sync.WaitGroup
	rollNumbers := make(chan string, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rn, err := engine.GenerateRollNumber(ctx, tenantID, programCode, academicYear)
			if err != nil {
				t.Errorf("concurrent roll number gen failed: %v", err)
				return
			}
			rollNumbers <- rn
		}()
	}

	wg.Wait()
	close(rollNumbers)

	seen := make(map[string]bool)
	for rn := range rollNumbers {
		if seen[rn] {
			t.Fatalf("collision detected in concurrent sequence generation: %s", rn)
		}
		seen[rn] = true
	}

	if len(seen) != concurrency {
		t.Fatalf("expected %d distinct roll numbers, got %d", concurrency, len(seen))
	}
}

func TestOnboarding_MultiTenantIsolation(t *testing.T) {
	ctx := context.Background()
	svc, _, _ := setupTestEnvironment()

	// Create applicant in Tenant A
	appA, err := svc.CreateDraft(ctx, CreateApplicantCommand{
		TenantID:  "tenant-A",
		Type:      TypeStudent,
		FirstName: "Tenant",
		LastName:  "A",
		Email:     "user@tenant-a.com",
	})
	if err != nil {
		t.Fatalf("failed to create applicant in tenant A: %v", err)
	}

	// Try fetching from Tenant B -> Must return ErrApplicantNotFound
	_, err = svc.GetApplicant(ctx, "tenant-B", appA.ID)
	if err != ErrApplicantNotFound {
		t.Fatalf("expected ErrApplicantNotFound across tenant boundary, got %v", err)
	}
}

func TestSequenceEngine_Formats_And_Validation(t *testing.T) {
	ctx := context.Background()
	repo := NewMockSequenceRepository()
	engine := NewSequenceEngine(repo)

	// Test Roll Number formats
	rn1, err := engine.GenerateRollNumber(ctx, "tenant-1", "CSE", "2026-2027")
	if err != nil || rn1 != "2026-CSE-0001" {
		t.Fatalf("expected 2026-CSE-0001, got %s (err: %v)", rn1, err)
	}

	rn2, err := engine.GenerateRollNumber(ctx, "tenant-1", "CSE", "2026-2027")
	if err != nil || rn2 != "2026-CSE-0002" {
		t.Fatalf("expected 2026-CSE-0002, got %s (err: %v)", rn2, err)
	}

	// Test Employee ID format
	emp1, err := engine.GenerateEmployeeID(ctx, "tenant-1", "CSE")
	if err != nil || emp1 != "EMP-CSE-0001" {
		t.Fatalf("expected EMP-CSE-0001, got %s (err: %v)", emp1, err)
	}

	// Test Registration Number format
	reg1, err := engine.GenerateRegistrationNumber(ctx, "tenant-1", 2026)
	if err != nil || reg1 != "REG-2026-00001" {
		t.Fatalf("expected REG-2026-00001, got %s (err: %v)", reg1, err)
	}

	// Validation errors
	_, err = engine.GenerateRollNumber(ctx, "", "CSE", "2026")
	if err != ErrTenantRequired {
		t.Fatalf("expected ErrTenantRequired, got %v", err)
	}
	_, err = engine.GenerateRollNumber(ctx, "tenant-1", "", "2026")
	if err == nil {
		t.Fatalf("expected error for empty program code")
	}
	_, err = engine.GenerateEmployeeID(ctx, "", "CSE")
	if err != ErrTenantRequired {
		t.Fatalf("expected ErrTenantRequired, got %v", err)
	}
	_, err = engine.GenerateEmployeeID(ctx, "tenant-1", "")
	if err == nil {
		t.Fatalf("expected error for empty dept code")
	}
}

func TestOnboarding_ListApplicants_And_Filters(t *testing.T) {
	ctx := context.Background()
	svc, acadRepo, _ := setupTestEnvironment()
	tenantID := "tenant-filter"

	dept := &AcademicDepartment{ID: "dept-1", TenantID: tenantID, Code: "CSE", Name: "Computer Science", IsActive: true}
	_ = acadRepo.CreateDepartment(ctx, dept)
	prog := &AcademicProgram{ID: "prog-1", TenantID: tenantID, DepartmentID: dept.ID, Code: "BTECH_CSE", Name: "B.Tech CSE", DegreeType: DegreeUG, IsActive: true}
	_ = acadRepo.CreateProgram(ctx, prog)
	cohort := &CohortBatch{ID: "coh-1", TenantID: tenantID, ProgramID: prog.ID, AcademicYear: "2026-2027", StartYear: 2026, EndYear: 2030, Section: "A", MaxCapacity: 60, IsActive: true}
	_ = acadRepo.CreateCohort(ctx, cohort)

	// Create 3 applicants
	_, _ = svc.CreateDraft(ctx, CreateApplicantCommand{TenantID: tenantID, Type: TypeStudent, FirstName: "Alice", LastName: "Smith", Email: "alice@test.com", ProgramID: prog.ID, AcademicYear: "2026-2027"})
	_, _ = svc.CreateDraft(ctx, CreateApplicantCommand{TenantID: tenantID, Type: TypeStudent, FirstName: "Bob", LastName: "Jones", Email: "bob@test.com", ProgramID: prog.ID, AcademicYear: "2026-2027"})
	_, _ = svc.CreateDraft(ctx, CreateApplicantCommand{TenantID: tenantID, Type: TypeFaculty, FirstName: "Charlie", LastName: "Brown", Email: "charlie@test.com", DepartmentID: dept.ID, AcademicYear: "2026-2027"})

	// 1. List All in Tenant
	list, total, err := svc.ListApplicants(ctx, ApplicantFilter{TenantID: tenantID})
	if err != nil || total != 3 || len(list) != 3 {
		t.Fatalf("expected 3 applicants, got total=%d, count=%d (err: %v)", total, len(list), err)
	}

	// 2. Filter by Type Student
	stuType := TypeStudent
	list, total, err = svc.ListApplicants(ctx, ApplicantFilter{TenantID: tenantID, Type: &stuType})
	if err != nil || total != 2 {
		t.Fatalf("expected 2 students, got %d", total)
	}

	// 3. Filter by Search Query
	list, total, err = svc.ListApplicants(ctx, ApplicantFilter{TenantID: tenantID, Search: "charlie"})
	if err != nil || total != 1 || list[0].FirstName != "Charlie" {
		t.Fatalf("expected 1 search result for 'charlie', got %d", total)
	}

	// 4. Test List Departments, Programs, Cohorts
	depts, err := svc.ListDepartments(ctx, tenantID)
	if err != nil || len(depts) != 1 {
		t.Fatalf("expected 1 department, got %d (err: %v)", len(depts), err)
	}
	progs, err := svc.ListPrograms(ctx, tenantID, dept.ID)
	if err != nil || len(progs) != 1 {
		t.Fatalf("expected 1 program, got %d (err: %v)", len(progs), err)
	}
	cohorts, err := svc.ListCohorts(ctx, tenantID, prog.ID, "2026-2027")
	if err != nil || len(cohorts) != 1 {
		t.Fatalf("expected 1 cohort, got %d (err: %v)", len(cohorts), err)
	}
}

func TestOnboarding_InvalidStateTransitions_EdgeCases(t *testing.T) {
	ctx := context.Background()
	svc, _, _ := setupTestEnvironment()
	tenantID := "tenant-transitions"

	app, _ := svc.CreateDraft(ctx, CreateApplicantCommand{
		TenantID:  tenantID,
		Type:      TypeStudent,
		FirstName: "Zack",
		LastName:  "Taylor",
		Email:     "zack@test.com",
	})

	// Try verifying draft directly -> Must fail
	_, err := svc.VerifyApplication(ctx, tenantID, app.ID)
	if err == nil {
		t.Fatalf("expected error verifying draft")
	}

	// Try enrolling draft directly -> Must fail
	_, err = svc.EnrollStudent(ctx, EnrollStudentCommand{
		TenantID:    tenantID,
		ApplicantID: app.ID,
		UserID:      "usr-1",
		CohortID:    "coh-1",
	})
	if err == nil {
		t.Fatalf("expected error enrolling draft")
	}

	// Try rejecting with empty reason -> Must fail
	_, err = svc.RejectApplication(ctx, tenantID, app.ID, "reviewer-1", "")
	if err != ErrRejectionReasonRequired {
		t.Fatalf("expected ErrRejectionReasonRequired, got %v", err)
	}

	// Test DomainError Formatting
	domErr := NewDomainError("TEST_CODE", "test error message", ErrInvalidInput)
	if domErr.Error() == "" || domErr.Unwrap() != ErrInvalidInput {
		t.Fatalf("unexpected domain error format: %s", domErr.Error())
	}
}
