/**
 * BLOCK_API_STUDYHUB_HANDLER_TEST_001
 * Subsystem: Rank 10 - Study Hub System (studyhub)
 * Purpose:   HTTP integration tests verifying REST endpoints for courses, materials, assignments, submissions, and peer reviews.
 */

package studyhub_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	apiStudyHub "campus/api/http/studyhub"
	backendStudyHub "campus/backend/studyhub"
)

func setupTestRouter() (chi.Router, backendStudyHub.Service) {
	r := chi.NewRouter()
	repo := backendStudyHub.NewMockRepository()
	svc := backendStudyHub.NewService(repo, nil)
	handler := apiStudyHub.NewHandler(svc)
	handler.RegisterRoutes(r)
	return r, svc
}

func TestHTTP_StudyHubCourseAndMaterials(t *testing.T) {
	router, _ := setupTestRouter()

	// 1. Create Course
	coursePayload := map[string]interface{}{
		"tenant_id":     "tenant-test",
		"code":          "CS-401",
		"title":         "Distributed Systems",
		"description":   "Consensus algorithms, Raft, Paxos, and distributed transactions.",
		"department":    "Computer Science",
		"credits":       4,
		"semester":      7,
		"syllabus_text": "Unit 1: Logical Clocks, Unit 2: Consensus, Unit 3: Distributed Storage",
	}
	body, _ := json.Marshal(coursePayload)
	req := httptest.NewRequest("POST", "/api/v1/studyhub/courses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for course, got %d: %s", rec.Code, rec.Body.String())
	}

	var course backendStudyHub.StudyCourse
	_ = json.Unmarshal(rec.Body.Bytes(), &course)

	// 2. Publish Material
	matPayload := map[string]interface{}{
		"tenant_id":       "tenant-test",
		"title":           "Raft Protocol Paper",
		"description":     "Original Ongaro & Ousterhout paper on In Search of an Understandable Consensus Algorithm.",
		"unit_number":     2,
		"material_type":   "NOTE",
		"file_url":        "https://storage.campus.internal/study/raft.pdf",
		"file_size_bytes": 524288,
	}
	matBody, _ := json.Marshal(matPayload)
	reqMat := httptest.NewRequest("POST", "/api/v1/studyhub/courses/"+course.ID+"/materials", bytes.NewReader(matBody))
	reqMat.Header.Set("Content-Type", "application/json")
	recMat := httptest.NewRecorder()
	router.ServeHTTP(recMat, reqMat)

	if recMat.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for material, got %d: %s", recMat.Code, recMat.Body.String())
	}

	// 3. List Courses
	reqList := httptest.NewRequest("GET", "/api/v1/studyhub/courses?tenant_id=tenant-test&department=Computer+Science", nil)
	recList := httptest.NewRecorder()
	router.ServeHTTP(recList, reqList)

	if recList.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for list courses, got %d: %s", recList.Code, recList.Body.String())
	}
}

func TestHTTP_StudyHubAssignmentAndSubmissionFlow(t *testing.T) {
	router, svc := setupTestRouter()
	ctx := context.Background()

	course, _ := svc.CreateCourse(ctx, backendStudyHub.CreateCourseRequest{
		TenantID:   "tenant-test",
		Code:       "CS-402",
		Title:      "Cloud Computing",
		Department: "Computer Science",
	})

	dueDate := time.Now().UTC().Add(48 * time.Hour).Format(time.RFC3339)

	// 1. Create Assignment
	asgnPayload := map[string]interface{}{
		"tenant_id":                   "tenant-test",
		"title":                       "Kubernetes Cluster Setup",
		"description":                 "Deploy a multi-node cluster using k3s and ingress controller.",
		"max_marks":                   100,
		"due_date":                    dueDate,
		"allow_late_submission":       true,
		"late_penalty_percent_per_day": 5,
	}
	asgnBody, _ := json.Marshal(asgnPayload)
	reqAsgn := httptest.NewRequest("POST", "/api/v1/studyhub/courses/"+course.ID+"/assignments", bytes.NewReader(asgnBody))
	reqAsgn.Header.Set("Content-Type", "application/json")
	recAsgn := httptest.NewRecorder()
	router.ServeHTTP(recAsgn, reqAsgn)

	if recAsgn.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for assignment, got %d: %s", recAsgn.Code, recAsgn.Body.String())
	}

	var asgn backendStudyHub.StudyAssignment
	_ = json.Unmarshal(recAsgn.Body.Bytes(), &asgn)

	// 2. Submit Assignment
	subPayload := map[string]interface{}{
		"tenant_id":    "tenant-test",
		"student_id":   "student-001",
		"file_url":     "https://storage.campus.internal/subs/k8s-cluster.yaml",
		"content_text": "K8s manifest submitted with helm chart.",
	}
	subBody, _ := json.Marshal(subPayload)
	reqSub := httptest.NewRequest("POST", "/api/v1/studyhub/assignments/"+asgn.ID+"/submissions", bytes.NewReader(subBody))
	reqSub.Header.Set("Content-Type", "application/json")
	recSub := httptest.NewRecorder()
	router.ServeHTTP(recSub, reqSub)

	if recSub.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for submission, got %d: %s", recSub.Code, recSub.Body.String())
	}

	var sub backendStudyHub.StudySubmission
	_ = json.Unmarshal(recSub.Body.Bytes(), &sub)

	// 3. Grade Submission
	gradePayload := map[string]interface{}{
		"tenant_id":    "tenant-test",
		"raw_marks":    94.5,
		"feedback":     "Flawless ingress rules and healthy pod replicas.",
		"graded_by_id": "faculty-prof-cloud",
	}
	gradeBody, _ := json.Marshal(gradePayload)
	reqGrade := httptest.NewRequest("POST", "/api/v1/studyhub/submissions/"+sub.ID+"/grade", bytes.NewReader(gradeBody))
	reqGrade.Header.Set("Content-Type", "application/json")
	recGrade := httptest.NewRecorder()
	router.ServeHTTP(recGrade, reqGrade)

	if recGrade.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for grade, got %d: %s", recGrade.Code, recGrade.Body.String())
	}

	// 4. Peer Review
	reviewPayload := map[string]interface{}{
		"tenant_id":            "tenant-test",
		"reviewer_student_id":  "student-002",
		"score":                92.0,
		"comments":             "Well structured ingress controllers and clean documentation.",
	}
	reviewBody, _ := json.Marshal(reviewPayload)
	reqReview := httptest.NewRequest("POST", "/api/v1/studyhub/submissions/"+sub.ID+"/reviews", bytes.NewReader(reviewBody))
	reqReview.Header.Set("Content-Type", "application/json")
	recReview := httptest.NewRecorder()
	router.ServeHTTP(recReview, reqReview)

	if recReview.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for peer review, got %d: %s", recReview.Code, recReview.Body.String())
	}
}
