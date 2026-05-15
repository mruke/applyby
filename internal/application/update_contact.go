package application

import (
	"context"
	"fmt"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// UpdateContactRepository
//
// Defines repository behavior required to update one contact.
// -----------------------------------------------------------------------------
type UpdateContactRepository interface {
	ContactFinder
	ContactUpdater
}

// -----------------------------------------------------------------------------
// UpdateContactInput
//
// Contains fields required to update an application contact.
// -----------------------------------------------------------------------------
type UpdateContactInput struct {
	ApplicationID domain.ApplicationID
	ContactID     domain.ContactID
	Name          string
	Email         string
	Role          string
}

// -----------------------------------------------------------------------------
// UpdateContactService
//
// Coordinates the workflow for editing an application contact.
// -----------------------------------------------------------------------------
type UpdateContactService struct {
	repository       UpdateContactRepository
	activityRecorder ActivityEventRecorder
}

// -----------------------------------------------------------------------------
// NewUpdateContactService
//
// Creates a service for the update contact workflow.
// -----------------------------------------------------------------------------
func NewUpdateContactService(repository UpdateContactRepository, activityRecorder ActivityEventRecorder) UpdateContactService {
	return UpdateContactService{
		repository:       repository,
		activityRecorder: activityRecorder,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Updates one contact and records an activity event.
// -----------------------------------------------------------------------------
func (service UpdateContactService) Execute(ctx context.Context, input UpdateContactInput) (domain.Contact, error) {
	if service.repository == nil {
		return domain.Contact{}, fmt.Errorf("contact updater repository is required")
	}

	if service.activityRecorder == nil {
		return domain.Contact{}, fmt.Errorf("activity recorder is required")
	}

	if err := input.ApplicationID.Validate(); err != nil {
		return domain.Contact{}, err
	}

	if err := input.ContactID.Validate(); err != nil {
		return domain.Contact{}, err
	}

	existingContact, err := service.repository.FindContactByID(ctx, input.ApplicationID, input.ContactID)
	if err != nil {
		return domain.Contact{}, err
	}

	updatedContact, err := domain.NewContact(existingContact.ID, existingContact.ApplicationID, input.Name, input.Email, input.Role)
	if err != nil {
		return domain.Contact{}, err
	}

	if err := service.repository.UpdateContact(ctx, updatedContact); err != nil {
		return domain.Contact{}, err
	}

	event, err := domain.NewActivityEvent(
		updatedContact.ApplicationID,
		domain.ActivityContactUpdated,
		time.Now().UTC(),
		fmt.Sprintf("Contact updated: %s.", updatedContact.Name),
	)
	if err != nil {
		return domain.Contact{}, err
	}

	if err := service.activityRecorder.RecordActivityEvent(ctx, event); err != nil {
		return domain.Contact{}, err
	}

	return updatedContact, nil
}
