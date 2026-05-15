package domain

import "testing"

// -----------------------------------------------------------------------------
// TestAllowedNextStatusesReturnsExpectedTransitions
//
// Verifies that each valid application status can move to every other valid
// status so users can correct or revise tracked application state.
// -----------------------------------------------------------------------------
func TestAllowedNextStatusesReturnsExpectedTransitions(t *testing.T) {
	allStatuses := AllApplicationStatuses()

	for _, status := range allStatuses {
		t.Run(status.String()+" can move to every other valid status", func(t *testing.T) {
			actual := AllowedNextStatuses(status)

			if len(actual) != len(allStatuses)-1 {
				t.Fatalf("expected %d next statuses, got %d", len(allStatuses)-1, len(actual))
			}

			for _, nextStatus := range actual {
				if nextStatus == status {
					t.Fatalf("expected current status %q to be excluded", status)
				}

				if !IsValidApplicationStatus(nextStatus) {
					t.Fatalf("expected next status %q to be valid", nextStatus)
				}
			}
		})
	}
}

// -----------------------------------------------------------------------------
// TestAllowedNextStatusesReturnsEmptyForInvalidStatus
//
// Verifies that invalid statuses do not expose transition options.
// -----------------------------------------------------------------------------
func TestAllowedNextStatusesReturnsEmptyForInvalidStatus(t *testing.T) {
	actual := AllowedNextStatuses(ApplicationStatus("paused"))

	if len(actual) != 0 {
		t.Fatalf("expected invalid status to return no next statuses")
	}
}

// -----------------------------------------------------------------------------
// TestAllowedNextStatusesReturnsCopy
//
// Verifies that callers cannot mutate generated transition results.
// -----------------------------------------------------------------------------
func TestAllowedNextStatusesReturnsCopy(t *testing.T) {
	nextStatuses := AllowedNextStatuses(StatusDraft)
	nextStatuses[0] = StatusDraft

	actual := AllowedNextStatuses(StatusDraft)

	if actual[0] == StatusDraft {
		t.Fatalf("expected lifecycle transition result to be protected from caller mutation")
	}
}

// -----------------------------------------------------------------------------
// TestCanTransitionApplicationStatusAcceptsValidTransition
//
// Verifies that movement between different valid statuses is allowed.
// -----------------------------------------------------------------------------
func TestCanTransitionApplicationStatusAcceptsValidTransition(t *testing.T) {
	if !CanTransitionApplicationStatus(StatusRejected, StatusInterviewing) {
		t.Fatal("expected rejected to interviewing transition to be allowed")
	}
}

// -----------------------------------------------------------------------------
// TestCanTransitionApplicationStatusRejectsSameStatus
//
// Verifies that no-op status updates are rejected.
// -----------------------------------------------------------------------------
func TestCanTransitionApplicationStatusRejectsSameStatus(t *testing.T) {
	if CanTransitionApplicationStatus(StatusApplied, StatusApplied) {
		t.Fatal("expected same-status transition to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestValidateApplicationStatusTransitionAcceptsValidTransition
//
// Verifies that validation accepts movement between different valid statuses.
// -----------------------------------------------------------------------------
func TestValidateApplicationStatusTransitionAcceptsValidTransition(t *testing.T) {
	transition := ApplicationStatusTransition{
		From: StatusRejected,
		To:   StatusInterviewing,
	}

	if err := ValidateApplicationStatusTransition(transition); err != nil {
		t.Fatalf("expected transition to be valid: %v", err)
	}
}

// -----------------------------------------------------------------------------
// TestValidateApplicationStatusTransitionRejectsInvalidSourceStatus
//
// Verifies that validation rejects an unsupported source status.
// -----------------------------------------------------------------------------
func TestValidateApplicationStatusTransitionRejectsInvalidSourceStatus(t *testing.T) {
	transition := ApplicationStatusTransition{
		From: ApplicationStatus("paused"),
		To:   StatusApplied,
	}

	if err := ValidateApplicationStatusTransition(transition); err == nil {
		t.Fatal("expected invalid source status to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestValidateApplicationStatusTransitionRejectsInvalidTargetStatus
//
// Verifies that validation rejects an unsupported target status.
// -----------------------------------------------------------------------------
func TestValidateApplicationStatusTransitionRejectsInvalidTargetStatus(t *testing.T) {
	transition := ApplicationStatusTransition{
		From: StatusApplied,
		To:   ApplicationStatus("paused"),
	}

	if err := ValidateApplicationStatusTransition(transition); err == nil {
		t.Fatal("expected invalid target status to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestValidateApplicationStatusTransitionRejectsSameStatus
//
// Verifies that validation rejects no-op status transitions.
// -----------------------------------------------------------------------------
func TestValidateApplicationStatusTransitionRejectsSameStatus(t *testing.T) {
	transition := ApplicationStatusTransition{
		From: StatusApplied,
		To:   StatusApplied,
	}

	if err := ValidateApplicationStatusTransition(transition); err == nil {
		t.Fatal("expected same-status transition to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestIsTerminalApplicationStatus
//
// Verifies that no valid status is terminal because tracker state is editable.
// -----------------------------------------------------------------------------
func TestIsTerminalApplicationStatus(t *testing.T) {
	for _, status := range AllApplicationStatuses() {
		if IsTerminalApplicationStatus(status) {
			t.Fatalf("expected status %q not to be terminal", status)
		}
	}
}
