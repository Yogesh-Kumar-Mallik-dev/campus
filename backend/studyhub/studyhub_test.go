/**
 * BLOCK_STUDYHUB_TEST_001
 * Subsystem: Rank 10 - Study Hub System (studyhub)
 * Purpose:   Unit test suite covering courses, materials, assignments, late penalties, grading, and peer reviews.
 */

package studyhub_test

import (
	"context"
	"testing"
	"time"

	"campus/backend/studyhub"
)

func setupStudyHubService() studyhub.Service {
	repo := studyhub.NewMockRepository()
	return studyhub.NewService(repo, nil)
}

func TestStudyHub_CourseAndMaterialManagement(t *testing.T) {
	ctx := context.Background()
	svc := setupStudyHubService()

	tenantID := "tenant-alpha"

	// 1. Create Course
	course, err := svc.CreateCourse(ctx, studyhub.CreateCourseRequest{
		TenantID:     tenantID,
		Code:         "cs-301",
		Title:        "Advanced Operating Systems",
		Description:  "Kernel design, memory paging, and concurrency primitives.",
		Department:   "Computer Science",
		Credits:      4,
		Semester:     5,
		SyllabusText: "Unit 1: Virtual Memory, Unit 2: File Systems, Unit 3: Concurrency",
	})
	if err != nil {
		t.Fatalf("unexpected error creating course: %v", err)
	}
	if course.Code != "CS-301" {
		t.Fatalf("expected course code to be normalized to CS-301, got %s", course.Code)
	}

	// 2. Duplicate Course Code Rejection
	_, err = svc.CreateCourse(ctx, studyhub.CreateCourseRequest{
		TenantID:   tenantID,
		Code:       "CS-301",
		Title:      "Duplicate Course",
		Department: "Computer Science",
	})
	if err != studyhub.ErrCourseCodeExists {
		t.Fatalf("expected ErrCourseCodeExists, got %v", err)
	}

	// 3. Publish Study Materials
	mat1, err := svc.PublishMaterial(ctx, studyhub.PublishMaterialRequest{
		TenantID:      tenantID,
		CourseID:      course.ID,
		Title:         "Lecture 01: Kernel Architecture",
		UnitNumber:    1,
		MaterialType:  studyhub.MaterialTypeNote,
		FileURL:       "https://storage.campus.internal/notes/cs301-01.pdf",
		FileSizeBytes: 1048576,
	})
	if err != nil {
		t.Fatalf("unexpected error publishing material 1: %v", err)
	}

	mat2, err := svc.PublishMaterial(ctx, studyhub.PublishMaterialRequest{
		TenantID:      tenantID,
		CourseID:      course.ID,
		Title:         "Lab Manual: Process Scheduling in C",
		UnitNumber:    1,
		MaterialType:  studyhub.MaterialTypeLabManual,
		FileURL:       "https://storage.campus.internal/labs/cs301-lab1.pdf",
		FileSizeBytes: 2097152,
	})
	if err != nil {
		t.Fatalf("unexpected error publishing material 2: %v", err)
	}

	// 4. List Materials
	materials, err := svc.ListMaterials(ctx, tenantID, course.ID, nil)
	if err != nil {
		t.Fatalf("unexpected error listing materials: %v", err)
	}
	if len(materials) != 2 {
		t.Fatalf("expected 2 materials, got %d", len(materials))
	}
	_ = mat1
	_ = mat2
}

func TestStudyHub_AssignmentCreationAndListing(t *testing.T) {
	ctx := context.Background()
	svc := setupStudyHubService()

	tenantID := "tenant-alpha"

	course, err := svc.CreateCourse(ctx, studyhub.CreateCourseRequest{
		TenantID:   tenantID,
		Code:       "CS-302",
		Title:      "Database Systems",
		Department: "Computer Science",
		Credits:    4,
		Semester:   5,
	})
	if err != nil {
		t.Fatalf("failed to create course: %v", err)
	}

	dueDate := time.Now().UTC().Add(7 * 24 * time.Hour)
	asgn, err := svc.CreateAssignment(ctx, studyhub.CreateAssignmentRequest{
		TenantID:                 tenantID,
		CourseID:                 course.ID,
		Title:                    "Assignment 1: B-Tree Index Implementation",
		Description:              "Implement an in-memory B-Tree index with insert and search operations.",
		MaxMarks:                 100,
		DueDate:                  dueDate,
		AllowLateSubmission:      true,
		LatePenaltyPercentPerDay: 5,
	})
	if err != nil {
		t.Fatalf("failed to create assignment: %v", err)
	}

	if asgn.MaxMarks != 100 || asgn.Status != studyhub.AssignmentStatusPublished {
		t.Fatalf("assignment attributes mismatch: %+v", asgn)
	}

	list, err := svc.ListAssignments(ctx, tenantID, course.ID)
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 assignment in list, got %d, err: %v", len(list), err)
	}
}

func TestStudyHub_OnTimeSubmissionAndGrading(t *testing.T) {
	ctx := context.Background()
	svc := setupStudyHubService()

	tenantID := "tenant-alpha"

	course, _ := svc.CreateCourse(ctx, studyhub.CreateCourseRequest{
		TenantID:   tenantID,
		Code:       "CS-303",
		Title:      "Computer Networks",
		Department: "Computer Science",
	})

	dueDate := time.Now().UTC().Add(3 * 24 * time.Hour)
	asgn, _ := svc.CreateAssignment(ctx, studyhub.CreateAssignmentRequest{
		TenantID:                 tenantID,
		CourseID:                 course.ID,
		Title:                    "Socket Programming in Go",
		MaxMarks:                 50,
		DueDate:                  dueDate,
		AllowLateSubmission:      true,
		LatePenaltyPercentPerDay: 5,
	})

	studentID := "student-001"
	subTime := time.Now().UTC()
	sub, err := svc.SubmitAssignment(ctx, studyhub.SubmitAssignmentRequest{
		TenantID:       tenantID,
		AssignmentID:   asgn.ID,
		StudentID:      studentID,
		FileURL:        "https://storage.campus.internal/subs/sub-01.zip",
		ContentText:    "Socket client and server implemented with TLS support.",
		SubmissionTime: &subTime,
	})
	if err != nil {
		t.Fatalf("failed to submit assignment: %v", err)
	}
	if sub.Status != studyhub.SubmissionStatusSubmitted {
		t.Fatalf("expected status SUBMITTED, got %s", sub.Status)
	}

	// Duplicate Submission Rejection
	_, err = svc.SubmitAssignment(ctx, studyhub.SubmitAssignmentRequest{
		TenantID:     tenantID,
		AssignmentID: asgn.ID,
		StudentID:    studentID,
	})
	if err != studyhub.ErrDuplicateSubmission {
		t.Fatalf("expected ErrDuplicateSubmission, got %v", err)
	}

	// Grade submission
	gradedSub, err := svc.GradeSubmission(ctx, studyhub.GradeSubmissionRequest{
		TenantID:   tenantID,
		ID:         sub.ID,
		RawMarks:   48,
		Feedback:   "Excellent socket error handling and TLS configuration.",
		GradedByID: "faculty-prof-sharma",
	})
	if err != nil {
		t.Fatalf("failed to grade submission: %v", err)
	}
	if gradedSub.Status != studyhub.SubmissionStatusGraded {
		t.Fatalf("expected status GRADED, got %s", gradedSub.Status)
	}
	if *gradedSub.MarksObtained != 48.0 {
		t.Fatalf("expected marks 48.0, got %f", *gradedSub.MarksObtained)
	}
}

func TestStudyHub_LateSubmissionAndPenaltyDeduction(t *testing.T) {
	ctx := context.Background()
	svc := setupStudyHubService()

	tenantID := "tenant-alpha"

	course, _ := svc.CreateCourse(ctx, studyhub.CreateCourseRequest{
		TenantID:   tenantID,
		Code:       "CS-304",
		Title:      "Compiler Design",
		Department: "Computer Science",
	})

	dueDate := time.Now().UTC().Add(-48 * time.Hour) // Due 2 days ago
	asgn, _ := svc.CreateAssignment(ctx, studyhub.CreateAssignmentRequest{
		TenantID:                 tenantID,
		CourseID:                 course.ID,
		Title:                    "AST Parser Construction",
		MaxMarks:                 100,
		DueDate:                  dueDate,
		AllowLateSubmission:      true,
		LatePenaltyPercentPerDay: 5, // 5% per day = 10% penalty for 2 days late = 10 marks deduction on 100 maxMarks
	})

	studentID := "student-002"
	now := time.Now().UTC()
	sub, err := svc.SubmitAssignment(ctx, studyhub.SubmitAssignmentRequest{
		TenantID:       tenantID,
		AssignmentID:   asgn.ID,
		StudentID:      studentID,
		FileURL:        "https://storage.campus.internal/subs/sub-02.tar.gz",
		SubmissionTime: &now,
	})
	if err != nil {
		t.Fatalf("failed to submit late assignment: %v", err)
	}
	if sub.Status != studyhub.SubmissionStatusLate {
		t.Fatalf("expected status LATE, got %s", sub.Status)
	}

	// Grade submission with raw 95
	gradedSub, err := svc.GradeSubmission(ctx, studyhub.GradeSubmissionRequest{
		TenantID:   tenantID,
		ID:         sub.ID,
		RawMarks:   95,
		Feedback:   "Great AST generator. 10 marks deducted for 2 days late submission.",
		GradedByID: "faculty-prof-sharma",
	})
	if err != nil {
		t.Fatalf("failed to grade late submission: %v", err)
	}

	// Expected: 95 - 10 = 85
	if *gradedSub.MarksObtained != 85.0 {
		t.Fatalf("expected adjusted marks to be 85.0, got %f", *gradedSub.MarksObtained)
	}
}

func TestStudyHub_ClosedAssignmentAndNoLatePolicy(t *testing.T) {
	ctx := context.Background()
	svc := setupStudyHubService()

	tenantID := "tenant-alpha"

	course, _ := svc.CreateCourse(ctx, studyhub.CreateCourseRequest{
		TenantID:   tenantID,
		Code:       "CS-305",
		Title:      "Software Architecture",
		Department: "Computer Science",
	})

	dueDate := time.Now().UTC().Add(-24 * time.Hour) // Past due
	asgn, _ := svc.CreateAssignment(ctx, studyhub.CreateAssignmentRequest{
		TenantID:                 tenantID,
		CourseID:                 course.ID,
		Title:                    "Microservices Architecture Plan",
		MaxMarks:                 100,
		DueDate:                  dueDate,
		AllowLateSubmission:      false, // Late not allowed
		LatePenaltyPercentPerDay: 0,
	})

	now := time.Now().UTC()
	_, err := svc.SubmitAssignment(ctx, studyhub.SubmitAssignmentRequest{
		TenantID:       tenantID,
		AssignmentID:   asgn.ID,
		StudentID:      "student-003",
		SubmissionTime: &now,
	})
	if err != studyhub.ErrLateSubmissionDisallowed {
		t.Fatalf("expected ErrLateSubmissionDisallowed, got %v", err)
	}
}

func TestStudyHub_PeerReviewLifecycleAndGuards(t *testing.T) {
	ctx := context.Background()
	svc := setupStudyHubService()

	tenantID := "tenant-alpha"

	course, _ := svc.CreateCourse(ctx, studyhub.CreateCourseRequest{
		TenantID:   tenantID,
		Code:       "CS-306",
		Title:      "Human Computer Interaction",
		Department: "Computer Science",
	})

	dueDate := time.Now().UTC().Add(48 * time.Hour)
	asgn, _ := svc.CreateAssignment(ctx, studyhub.CreateAssignmentRequest{
		TenantID:            tenantID,
		CourseID:            course.ID,
		Title:               "Wireframe Prototype Analysis",
		MaxMarks:            100,
		DueDate:             dueDate,
		AllowLateSubmission: true,
	})

	submitterID := "student-001"
	sub, _ := svc.SubmitAssignment(ctx, studyhub.SubmitAssignmentRequest{
		TenantID:     tenantID,
		AssignmentID: asgn.ID,
		StudentID:    submitterID,
		ContentText:  "Figma prototype submission link.",
	})

	// 1. Self Review Guard
	_, err := svc.SubmitPeerReview(ctx, studyhub.SubmitPeerReviewRequest{
		TenantID:          tenantID,
		SubmissionID:      sub.ID,
		ReviewerStudentID: submitterID,
		Score:             90,
		Comments:          "Self reviewing my own work",
	})
	if err != studyhub.ErrSelfPeerReviewNotAllowed {
		t.Fatalf("expected ErrSelfPeerReviewNotAllowed, got %v", err)
	}

	// 2. Score Range Guard
	_, err = svc.SubmitPeerReview(ctx, studyhub.SubmitPeerReviewRequest{
		TenantID:          tenantID,
		SubmissionID:      sub.ID,
		ReviewerStudentID: "student-002",
		Score:             110,
		Comments:          "Score exceeds 100",
	})
	if err != studyhub.ErrInvalidScoreRange {
		t.Fatalf("expected ErrInvalidScoreRange, got %v", err)
	}

	// 3. Valid Peer Review
	reviewerID := "student-002"
	review, err := svc.SubmitPeerReview(ctx, studyhub.SubmitPeerReviewRequest{
		TenantID:          tenantID,
		SubmissionID:      sub.ID,
		ReviewerStudentID: reviewerID,
		Score:             88.5,
		Comments:          "Clear navigation hierarchy and great visual contrast.",
	})
	if err != nil {
		t.Fatalf("failed to submit peer review: %v", err)
	}
	if review.Score != 88.5 {
		t.Fatalf("expected score 88.5, got %f", review.Score)
	}

	// 4. Duplicate Peer Review Guard
	_, err = svc.SubmitPeerReview(ctx, studyhub.SubmitPeerReviewRequest{
		TenantID:          tenantID,
		SubmissionID:      sub.ID,
		ReviewerStudentID: reviewerID,
		Score:             90,
	})
	if err != studyhub.ErrDuplicatePeerReview {
		t.Fatalf("expected ErrDuplicatePeerReview, got %v", err)
	}

	// 5. List Reviews
	reviews, err := svc.ListPeerReviews(ctx, sub.ID)
	if err != nil || len(reviews) != 1 {
		t.Fatalf("expected 1 review, got %d, err: %v", len(reviews), err)
	}
}
