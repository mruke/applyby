package application

import (
	"context"
	"fmt"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// RemoveContactRepository
//
// Defines repository behavior required to remove one contact.
// -----------------------------------------------------------------------------
type RemoveContactRepository interface {
	ContactFinder
	ContactRemover
}

// -----------------------------------------------------------------------------
// RemoveContactInput
//
// Contains identifiers required to remove an application contact.
// -----------------------------------------------------------------------------
type RemoveContactInput struct {
	ApplicationID domain.ApplicationID
	ContactID     domain.ContactID
}

// -----------------------------------------------------------------------------
// RemoveContactService
//
// Coordinates the workflow for removing an application contact.
// -----------------------------------------------------------------------------
type RemoveContactService struct {
	repository       RemoveContactRepository
	activityRecorder ActivityEventRecorder
}

// -----------------------------------------------------------------------------
// NewRemoveContactService
//
// Creates a service for the remove contact workflow.
// -----------------------------------------------------------------------------
func NewRemoveContactService(repository RemoveContactRepository, activityRecorder ActivityEventRecorder) RemoveContactService {
	return RemoveContactService{
		repository:       repository,
		activityRecorder: activityRecorder,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Removes one contact and records an activity event.
// -----------------------------------------------------------------------------
func (service RemoveContactService) Execute(ctx context.Context, input RemoveContactInput) error {
	if service.repository == nil {
		return fmt.Errorf("contact remover repository is required")
	}

	if service.activityRecorder == nil {
		return fmt.Errorf("activity recorder is required")
	}

	if err := input.ApplicationID.Validate(); err != nil {
		return err
	}

	if err := input.ContactID.Validate(); err != nil {
		return err
	}

	contact, err := service.repository.FindContactByID(ctx, input.ApplicationID, input.ContactID)
	if err != nil {
		return err
	}

	if err := service.repository.RemoveContact(ctx, input.ApplicationID, input.ContactID); err != nil {
		return err
	}

	event, err := domain.NewActivityEvent(
		contact.ApplicationID,
		domain.ActivityContactRemoved,
		time.Now().UTC(),
		fmt.Sprintf("Contact removed: %s.", contact.Name),
	)
	if err != nil {
		return err
	}

	return service.activityRecorder.RecordActivityEvent(ctx, event)
}
