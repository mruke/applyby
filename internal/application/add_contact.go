package application

import (
	"context"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
)

type AddContactInput struct {
	ID            domain.ContactID
	ApplicationID domain.ApplicationID
	Name          string
	Email         string
	Role          string
}

type AddContactService struct {
	repository       ContactSaver
	activityRecorder ActivityEventRecorder
}

func NewAddContactService(repository ContactSaver, activityRecorder ActivityEventRecorder) AddContactService {
	return AddContactService{
		repository:       repository,
		activityRecorder: activityRecorder,
	}
}

func (service AddContactService) Execute(ctx context.Context, input AddContactInput) (domain.Contact, error) {
	if service.repository == nil {
		return domain.Contact{}, fmt.Errorf("contact saver is required")
	}

	if service.activityRecorder == nil {
		return domain.Contact{}, fmt.Errorf("activity recorder is required")
	}

	contact, err := domain.NewContact(input.ID, input.ApplicationID, input.Name, input.Email, input.Role)
	if err != nil {
		return domain.Contact{}, err
	}

	if err := service.repository.SaveContact(ctx, contact); err != nil {
		return domain.Contact{}, err
	}
	if err := recordActivityEvent(
		ctx,
		service.activityRecorder,
		contact.ApplicationID,
		domain.ActivityContactAdded,
		fmt.Sprintf("Contact added: %s.", contact.Name),
	); err != nil {
		return domain.Contact{}, err
	}

	return contact, nil
}
