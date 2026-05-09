package domain

import (
	"fmt"
	"strings"
	"time"
)

// -----------------------------------------------------------------------------
// requireNonEmptyField
//
// Validates that a required string field contains non-whitespace text.
// -----------------------------------------------------------------------------
func requireNonEmptyField(fieldName string, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}

	return nil
}

// -----------------------------------------------------------------------------
// requireNonZeroTime
//
// Validates that a required time field has been set.
// -----------------------------------------------------------------------------
func requireNonZeroTime(fieldName string, value time.Time) error {
	if value.IsZero() {
		return fmt.Errorf("%s is required", fieldName)
	}

	return nil
}
