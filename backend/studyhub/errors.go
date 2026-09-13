/**
 * BLOCK_STUDYHUB_ERRORS_001
 * Subsystem: Rank 10 - Study Hub System (studyhub)
 * Purpose:   Domain error sentinels and RFC 7807 problem details mapping for study hub, courses, materials, and submissions.
 */

package studyhub

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrCourseNotFound             = errors.New("study course not found")
	ErrCourseCodeExists           = errors.New("course with this code already exists in tenant")
	ErrMaterialNotFound           = errors.New("study material not found")
	ErrAssignmentNotFound         = errors.New("study assignment not found")
	ErrAssignmentClosed           = errors.New("assignment is closed and no longer accepting submissions")
	ErrLateSubmissionDisallowed   = errors.New("assignment due date has passed and late submissions are not allowed")
	ErrSubmissionNotFound         = errors.New("assignment submission not found")
	ErrDuplicateSubmission        = errors.New("student has already submitted for this assignment")
	ErrMarksExceedMaxMarks        = errors.New("obtained marks cannot exceed assignment maximum marks")
	ErrSelfPeerReviewNotAllowed   = errors.New("students cannot peer review their own submissions")
	ErrDuplicatePeerReview        = errors.New("reviewer has already submitted peer evaluation for this submission")
	ErrInvalidScoreRange          = errors.New("peer review score must be between 0 and 100")
)

type DomainError struct {
	Err        error
	HTTPStatus int
	Title      string
	Detail     string
	Type       string
}

func (e *DomainError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s", e.Title, e.Detail)
	}
	return e.Title
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

func NewDomainError(err error, status int, title, detail, problemType string) *DomainError {
	return &DomainError{
		Err:        err,
		HTTPStatus: status,
		Title:      title,
		Detail:     detail,
		Type:       problemType,
	}
}

func MapErrorToProblem(err error) *DomainError {
	if err == nil {
		return nil
	}

	var domErr *DomainError
	if errors.As(err, &domErr) {
		return domErr
	}

	switch {
	case errors.Is(err, ErrCourseNotFound), errors.Is(err, ErrMaterialNotFound),
		errors.Is(err, ErrAssignmentNotFound), errors.Is(err, ErrSubmissionNotFound):
		return NewDomainError(err, http.StatusNotFound, "Resource Not Found", err.Error(), "https://campus.internal/errors/not-found")

	case errors.Is(err, ErrCourseCodeExists), errors.Is(err, ErrDuplicateSubmission),
		errors.Is(err, ErrAssignmentClosed), errors.Is(err, ErrDuplicatePeerReview):
		return NewDomainError(err, http.StatusConflict, "Resource Conflict / Policy Guard", err.Error(), "https://campus.internal/errors/conflict")

	case errors.Is(err, ErrLateSubmissionDisallowed), errors.Is(err, ErrMarksExceedMaxMarks),
		errors.Is(err, ErrSelfPeerReviewNotAllowed), errors.Is(err, ErrInvalidScoreRange):
		return NewDomainError(err, http.StatusUnprocessableEntity, "Unprocessable Entity", err.Error(), "https://campus.internal/errors/unprocessable")

	default:
		return NewDomainError(err, http.StatusInternalServerError, "Internal Server Error", "An unexpected study hub error occurred", "https://campus.internal/errors/internal")
	}
}
