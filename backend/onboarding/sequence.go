/**
 * BLOCK_ONBOARDING_SEQUENCE_001
 * Subsystem: Rank 3 - Student & Staff Registration System (onboarding)
 * Purpose:   Deterministic, concurrency-safe roll number and employee ID sequence generator.
 * Standard:  Proprietary & Enterprise Strict Modular Monolith
 */

package onboarding

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type SequenceGenerator interface {
	GenerateRollNumber(ctx context.Context, tenantID, programCode, academicYear string) (string, error)
	GenerateEmployeeID(ctx context.Context, tenantID, departmentCode string) (string, error)
	GenerateRegistrationNumber(ctx context.Context, tenantID string, year int) (string, error)
}

type SequenceEngine struct {
	repo SequenceRepository
}

func NewSequenceEngine(repo SequenceRepository) *SequenceEngine {
	return &SequenceEngine{repo: repo}
}

// GenerateRollNumber produces deterministic roll numbers formatted as {Year}-{ProgramCode}-{Seq:04d}
// e.g. 2026-BTECH_CSE-0001
func (s *SequenceEngine) GenerateRollNumber(ctx context.Context, tenantID, programCode, academicYear string) (string, error) {
	if strings.TrimSpace(tenantID) == "" {
		return "", ErrTenantRequired
	}
	cleanProgram := strings.ToUpper(strings.TrimSpace(programCode))
	if cleanProgram == "" {
		return "", NewDomainError("INVALID_PROGRAM", "program code required for roll number", ErrInvalidInput)
	}

	yearPrefix := academicYear
	if len(academicYear) >= 4 {
		yearPrefix = academicYear[:4]
	} else {
		yearPrefix = fmt.Sprintf("%d", time.Now().UTC().Year())
	}

	prefix := fmt.Sprintf("%s-%s", yearPrefix, cleanProgram)
	seq, err := s.repo.NextSequence(ctx, tenantID, "ROLL_NUMBER", prefix)
	if err != nil {
		return "", fmt.Errorf("failed to reserve roll number sequence: %w", err)
	}

	return fmt.Sprintf("%s-%04d", prefix, seq), nil
}

// GenerateEmployeeID produces deterministic employee IDs formatted as EMP-{DeptCode}-{Seq:04d}
// e.g. EMP-CSE-0001
func (s *SequenceEngine) GenerateEmployeeID(ctx context.Context, tenantID, departmentCode string) (string, error) {
	if strings.TrimSpace(tenantID) == "" {
		return "", ErrTenantRequired
	}
	cleanDept := strings.ToUpper(strings.TrimSpace(departmentCode))
	if cleanDept == "" {
		return "", NewDomainError("INVALID_DEPT", "department code required for employee ID", ErrInvalidInput)
	}

	prefix := fmt.Sprintf("EMP-%s", cleanDept)
	seq, err := s.repo.NextSequence(ctx, tenantID, "EMPLOYEE_ID", prefix)
	if err != nil {
		return "", fmt.Errorf("failed to reserve employee ID sequence: %w", err)
	}

	return fmt.Sprintf("%s-%04d", prefix, seq), nil
}

// GenerateRegistrationNumber produces institutional registration numbers formatted as REG-{Year}-{Seq:05d}
func (s *SequenceEngine) GenerateRegistrationNumber(ctx context.Context, tenantID string, year int) (string, error) {
	if strings.TrimSpace(tenantID) == "" {
		return "", ErrTenantRequired
	}
	if year == 0 {
		year = time.Now().UTC().Year()
	}

	prefix := fmt.Sprintf("REG-%d", year)
	seq, err := s.repo.NextSequence(ctx, tenantID, "REGISTRATION_NUMBER", prefix)
	if err != nil {
		return "", fmt.Errorf("failed to reserve registration sequence: %w", err)
	}

	return fmt.Sprintf("%s-%05d", prefix, seq), nil
}
