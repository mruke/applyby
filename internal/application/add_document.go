package application

import (
	"context"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
)

type AddDocumentInput struct {
	ID            domain.DocumentID
	ApplicationID domain.ApplicationID
	Name          string
	Kind          string
	Path          string
}

type AddDocumentService struct {
	repository       DocumentSaver
	activityRecorder ActivityEventRecorder
}

func NewAddDocumentService(repository DocumentSaver, activityRecorder ActivityEventRecorder) AddDocumentService {
	return AddDocumentService{
		repository:       repository,
		activityRecorder: activityRecorder,
	}
}

func (service AddDocumentService) Execute(ctx context.Context, input AddDocumentInput) (domain.Document, error) {
	if service.repository == nil {
		return domain.Document{}, fmt.Errorf("document saver is required")
	}

	if service.activityRecorder == nil {
		return domain.Document{}, fmt.Errorf("activity recorder is required")
	}

	document, err := domain.NewDocument(input.ID, input.ApplicationID, input.Name, input.Kind, input.Path)
	if err != nil {
		return domain.Document{}, err
	}

	if err := service.repository.SaveDocument(ctx, document); err != nil {
		return domain.Document{}, err
	}
	if err := recordActivityEvent(
		ctx,
		service.activityRecorder,
		document.ApplicationID,
		domain.ActivityDocumentAdded,
		fmt.Sprintf("Document metadata added: %s.", document.Name),
	); err != nil {
		return domain.Document{}, err
	}

	return document, nil
}
