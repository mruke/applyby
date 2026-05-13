package application

import (
	"context"
	"fmt"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

type CreateApplicationInput struct {
	ID        domain.ApplicationID
	Title     string
	Company   domain.Company
	Status    domain.ApplicationStatus
	Source    string
	Notes     string
	CreatedAt time.Time
}

type CreateApplicationService struct {
	repository       ApplicationSaver
	activityRecorder ActivityEventRecorder
}

func NewCreateApplicationService(repository ApplicationSaver, activityRecorder ActivityEventRecorder) CreateApplicationService {
	return CreateApplicationService{
		repository:       repository,
		activityRecorder: activityRecorder,
	}
}

func (service CreateApplicationService) Execute(ctx context.Context, input CreateApplicationInput) (domain.Application, error) {
	if service.repository == nil {
		return domain.Application{}, fmt.Errorf("application saver is required")
	}

	if service.activityRecorder == nil {
		return domain.Application{}, fmt.Errorf("activity recorder is required")
	}

	application, err := domain.NewApplication(input.ID, input.Title, input.Company, input.Status, input.CreatedAt)
	if err != nil {
		return domain.Application{}, err
	}

	application.Source = input.Source
	application.Notes = input.Notes

	if err := service.repository.SaveApplication(ctx, application); err != nil {
		return domain.Application{}, err
	}

	event, err := domain.NewActivityEvent(
		application.ID,
		domain.ActivityApplicationCreated,
		application.CreatedAt,
		fmt.Sprintf("Application created: %s at %s.", application.Title, application.Company.Name),
	)
	if err != nil {
		return domain.Application{}, err
	}

	if err := service.activityRecorder.RecordActivityEvent(ctx, event); err != nil {
		return domain.Application{}, err
	}

	return application, nil
}
