package domain

import "testing"

// -----------------------------------------------------------------------------
// TestAllowedNextStatusesReturnsExpectedTransitions
//
// Verifies that each application status exposes its expected next statuses.
// -----------------------------------------------------------------------------
func TestAllowedNextStatusesReturnsExpectedTransitions(t *testing.T) {
	tests := []struct {
		name     string
		status   ApplicationStatus
		expected []ApplicationStatus
	}{
		{
			name:     "draft can move to early active or closed states",
			status:   StatusDraft,
			expected: []ApplicationStatus{StatusInterested, StatusApplied, StatusWithdrawn, StatusArchived},
		},
		{
			name:     "interested can move to applied or closed states",
			status:   StatusInterested,
			expected: []ApplicationStatus{StatusApplied, StatusWithdrawn, StatusArchived},
		},
		{
			name:     "applied can move through the active pipeline or closed states",
			status:   StatusApplied,
			expected: []ApplicationStatus{StatusInterviewing, StatusOffer, StatusRejected, StatusWithdrawn, StatusArchived},
		},
		{
			name:     "interviewing can move to outcome or closed states",
			status:   StatusInterviewing,
			expected: []ApplicationStatus{StatusOffer, StatusRejected, StatusWithdrawn, StatusArchived},
		},
		{
			name:     "offer can only be archived",
			status:   StatusOffer,
			expected: []ApplicationStatus{StatusArchived},
		},
		{
			name:     "rejected can only be archived",
			status:   StatusRejected,
			expected: []ApplicationStatus{StatusArchived},
		},
		{
			name:     "withdrawn can only be archived",
			status:   StatusWithdrawn,
			expected: []ApplicationStatus{StatusArchived},
		},
		{
			name:     "archived has no next statuses",
			status:   StatusArchived,
			expected: []ApplicationStatus{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := AllowedNextStatuses(test.status)

			if len(actual) != len(test.expected) {
				t.Fatalf("expected %d next statuses, got %d", len(test.expected), len(actual))
			}

			for index, expectedStatus := range test.expected {
				if actual[index] != expectedStatus {
					t.Fatalf("expected status at index %d to be %q, got %q", index, expectedStatus, actual[index])
				}
			}
		})
	}
}

// -----------------------------------------------------------------------------
// TestAllowedNextStatusesReturnsCopy
//
// Verifies that callers cannot mutate the stored lifecycle transition table.
// -----------------------------------------------------------------------------
func TestAllowedNextStatusesReturnsCopy(t *testing.T) {
	nextStatuses := AllowedNextStatuses(StatusDraft)
	nextStatuses[0] = StatusRejected

	actual := AllowedNextStatuses(StatusDraft)

	if actual[0] != StatusInterested {
		t.Fatalf("expected lifecycle transition table to be protected from caller mutation")
	}
}

// -----------------------------------------------------------------------------
// TestCanTransitionApplicationStatusAcceptsAllowedTransition
//
// Verifies that an allowed lifecycle transition returns true.
// -----------------------------------------------------------------------------
func TestCanTransitionApplicationStatusAcceptsAllowedTransition(t *testing.T) {
	if !CanTransitionApplicationStatus(StatusApplied, StatusInterviewing) {
		t.Fatal("expected applied to interviewing transition to be allowed")
	}
}

// -----------------------------------------------------------------------------
// TestCanTransitionApplicationStatusRejectsInvalidTransition
//
// Verifies that an invalid lifecycle transition returns false.
// -----------------------------------------------------------------------------
func TestCanTransitionApplicationStatusRejectsInvalidTransition(t *testing.T) {
	if CanTransitionApplicationStatus(StatusRejected, StatusInterviewing) {
		t.Fatal("expected rejected to interviewing transition to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestValidateApplicationStatusTransitionAcceptsAllowedTransition
//
// Verifies that validation accepts a supported lifecycle transition.
// -----------------------------------------------------------------------------
func TestValidateApplicationStatusTransitionAcceptsAllowedTransition(t *testing.T) {
	transition := ApplicationStatusTransition{
		From: StatusApplied,
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
// TestValidateApplicationStatusTransitionRejectsUnsupportedTransition
//
// Verifies that validation rejects unsupported movement between valid statuses.
// -----------------------------------------------------------------------------
func TestValidateApplicationStatusTransitionRejectsUnsupportedTransition(t *testing.T) {
	transition := ApplicationStatusTransition{
		From: StatusRejected,
		To:   StatusInterviewing,
	}

	if err := ValidateApplicationStatusTransition(transition); err == nil {
		t.Fatal("expected unsupported transition to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestIsTerminalApplicationStatus
//
// Verifies that only statuses with no forward lifecycle movement are terminal.
// -----------------------------------------------------------------------------
func TestIsTerminalApplicationStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   ApplicationStatus
		expected bool
	}{
		{
			name:     "draft is not terminal",
			status:   StatusDraft,
			expected: false,
		},
		{
			name:     "offer is not terminal because it can be archived",
			status:   StatusOffer,
			expected: false,
		},
		{
			name:     "archived is terminal",
			status:   StatusArchived,
			expected: true,
		},
		{
			name:     "unknown status is not treated as terminal",
			status:   ApplicationStatus("paused"),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := IsTerminalApplicationStatus(test.status)

			if actual != test.expected {
				t.Fatalf("expected terminal result %v, got %v", test.expected, actual)
			}
		})
	}
}
