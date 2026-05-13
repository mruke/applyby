package application

import (
	"context"
	"errors"
	"testing"
)

// -----------------------------------------------------------------------------
// TestAddContactServiceSavesValidContact
//
// Verifies that the add contact workflow validates and saves a contact.
// -----------------------------------------------------------------------------
func TestAddContactServiceSavesValidContact(t *testing.T) {
	repository := newFakeContactRepository()
	service := NewAddContactService(repository, &fakeApplicationHistoryRepository{})

	contact, err := service.Execute(context.Background(), AddContactInput{
		ID:            "contact-001",
		ApplicationID: "app-001",
		Name:          "Sam Recruiter",
		Email:         "sam@example.com",
		Role:          "Recruiter",
	})

	if err != nil {
		t.Fatalf("expected add contact workflow to succeed: %v", err)
	}

	if contact.ID != "contact-001" {
		t.Fatalf("expected contact id to be preserved")
	}

	if repository.saveCalls != 1 {
		t.Fatalf("expected repository save to be called once")
	}
}

// -----------------------------------------------------------------------------
// TestAddContactServiceRejectsInvalidContact
//
// Verifies that invalid contact data is rejected before saving.
// -----------------------------------------------------------------------------
func TestAddContactServiceRejectsInvalidContact(t *testing.T) {
	repository := newFakeContactRepository()
	service := NewAddContactService(repository, &fakeApplicationHistoryRepository{})

	_, err := service.Execute(context.Background(), AddContactInput{
		ID:            "",
		ApplicationID: "app-001",
		Name:          "Sam Recruiter",
		Email:         "sam@example.com",
		Role:          "Recruiter",
	})

	if err == nil {
		t.Fatal("expected invalid contact to be rejected")
	}

	if repository.saveCalls != 0 {
		t.Fatal("expected invalid contact not to be saved")
	}
}

// -----------------------------------------------------------------------------
// TestAddContactServiceReturnsRepositoryError
//
// Verifies that repository save errors are returned to the caller.
// -----------------------------------------------------------------------------
func TestAddContactServiceReturnsRepositoryError(t *testing.T) {
	repository := newFakeContactRepository()
	repository.saveErr = errors.New("save failed")
	service := NewAddContactService(repository, &fakeApplicationHistoryRepository{})

	_, err := service.Execute(context.Background(), AddContactInput{
		ID:            "contact-001",
		ApplicationID: "app-001",
		Name:          "Sam Recruiter",
		Email:         "sam@example.com",
		Role:          "Recruiter",
	})

	if err == nil {
		t.Fatal("expected repository error to be returned")
	}
}

// -----------------------------------------------------------------------------
// TestAddContactServiceRequiresRepository
//
// Verifies that the add contact workflow requires a repository boundary.
// -----------------------------------------------------------------------------
func TestAddContactServiceRequiresRepository(t *testing.T) {
	service := NewAddContactService(nil, &fakeApplicationHistoryRepository{})

	_, err := service.Execute(context.Background(), AddContactInput{
		ID:            "contact-001",
		ApplicationID: "app-001",
		Name:          "Sam Recruiter",
		Email:         "sam@example.com",
		Role:          "Recruiter",
	})

	if err == nil {
		t.Fatal("expected missing repository to be rejected")
	}
}
