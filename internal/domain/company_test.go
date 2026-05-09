package domain

import "testing"

// -----------------------------------------------------------------------------
// TestNewCompanyAcceptsValidCompany
//
// Verifies that a company with a name can be created.
// -----------------------------------------------------------------------------
func TestNewCompanyAcceptsValidCompany(t *testing.T) {
	company, err := NewCompany("Example Studio", "https://example.com")

	if err != nil {
		t.Fatalf("expected company to be valid: %v", err)
	}

	if company.Name != "Example Studio" {
		t.Fatalf("expected company name to be preserved")
	}
}

// -----------------------------------------------------------------------------
// TestNewCompanyRejectsMissingName
//
// Verifies that a company without a name is rejected.
// -----------------------------------------------------------------------------
func TestNewCompanyRejectsMissingName(t *testing.T) {
	_, err := NewCompany(" ", "https://example.com")

	if err == nil {
		t.Fatal("expected company without a name to be invalid")
	}
}
