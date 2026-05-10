package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mruke/applyby/internal/application"
	"github.com/mruke/applyby/internal/domain"
	"github.com/mruke/applyby/internal/search"
)

// -----------------------------------------------------------------------------
// fakeSearchApplicationsExecutor
//
// Provides a fake search workflow for API handler tests.
// -----------------------------------------------------------------------------
type fakeSearchApplicationsExecutor struct {
	applications []domain.Application
	err          error
	called       bool
}

// -----------------------------------------------------------------------------
// Execute
//
// Records search workflow execution and returns the configured fake result.
// -----------------------------------------------------------------------------
func (executor *fakeSearchApplicationsExecutor) Execute(ctx context.Context, criteria search.ApplicationCriteria) ([]domain.Application, error) {
	executor.called = true

	if executor.err != nil {
		return nil, executor.err
	}

	return executor.applications, nil
}

// -----------------------------------------------------------------------------
// fakeListActivityEventsExecutor
//
// Provides a fake activity timeline workflow for API handler tests.
// -----------------------------------------------------------------------------
type fakeListActivityEventsExecutor struct {
	events []domain.ActivityEvent
	err    error
	called bool
}

// -----------------------------------------------------------------------------
// Execute
//
// Records activity timeline workflow execution and returns the configured fake result.
// -----------------------------------------------------------------------------
func (executor *fakeListActivityEventsExecutor) Execute(ctx context.Context, input application.ListActivityEventsInput) ([]domain.ActivityEvent, error) {
	executor.called = true

	if executor.err != nil {
		return nil, executor.err
	}

	return executor.events, nil
}

// -----------------------------------------------------------------------------
// fakeScheduleReminderExecutor
//
// Provides a fake schedule reminder workflow for API handler tests.
// -----------------------------------------------------------------------------
type fakeScheduleReminderExecutor struct {
	reminder domain.Reminder
	err      error
	called   bool
}

// -----------------------------------------------------------------------------
// Execute
//
// Records schedule reminder workflow execution and returns the configured fake result.
// -----------------------------------------------------------------------------
func (executor *fakeScheduleReminderExecutor) Execute(ctx context.Context, input application.ScheduleReminderInput) (domain.Reminder, error) {
	executor.called = true

	if executor.err != nil {
		return domain.Reminder{}, executor.err
	}

	return executor.reminder, nil
}

// -----------------------------------------------------------------------------
// fakeListRemindersExecutor
//
// Provides a fake list reminders workflow for API handler tests.
// -----------------------------------------------------------------------------
type fakeListRemindersExecutor struct {
	reminders []domain.Reminder
	err       error
	called    bool
}

// -----------------------------------------------------------------------------
// Execute
//
// Records list reminders workflow execution and returns the configured fake result.
// -----------------------------------------------------------------------------
func (executor *fakeListRemindersExecutor) Execute(ctx context.Context, input application.ListRemindersInput) ([]domain.Reminder, error) {
	executor.called = true

	if executor.err != nil {
		return nil, executor.err
	}

	return executor.reminders, nil
}

// -----------------------------------------------------------------------------
// fakeCompleteReminderExecutor
//
// Provides a fake complete reminder workflow for API handler tests.
// -----------------------------------------------------------------------------
type fakeCompleteReminderExecutor struct {
	reminder domain.Reminder
	err      error
	called   bool
}

// -----------------------------------------------------------------------------
// Execute
//
// Records complete reminder workflow execution and returns the configured fake result.
// -----------------------------------------------------------------------------
func (executor *fakeCompleteReminderExecutor) Execute(ctx context.Context, input application.CompleteReminderInput) (domain.Reminder, error) {
	executor.called = true

	if executor.err != nil {
		return domain.Reminder{}, executor.err
	}

	return executor.reminder, nil
}

// -----------------------------------------------------------------------------
// fakeAddContactExecutor
//
// Provides a fake add contact workflow for API handler tests.
// -----------------------------------------------------------------------------
type fakeAddContactExecutor struct {
	contact domain.Contact
	err     error
	called  bool
}

// -----------------------------------------------------------------------------
// Execute
//
// Records add contact workflow execution and returns the configured fake result.
// -----------------------------------------------------------------------------
func (executor *fakeAddContactExecutor) Execute(ctx context.Context, input application.AddContactInput) (domain.Contact, error) {
	executor.called = true

	if executor.err != nil {
		return domain.Contact{}, executor.err
	}

	return executor.contact, nil
}

// -----------------------------------------------------------------------------
// fakeListContactsExecutor
//
// Provides a fake list contacts workflow for API handler tests.
// -----------------------------------------------------------------------------
type fakeListContactsExecutor struct {
	contacts []domain.Contact
	err      error
	called   bool
}

// -----------------------------------------------------------------------------
// Execute
//
// Records list contacts workflow execution and returns the configured fake result.
// -----------------------------------------------------------------------------
func (executor *fakeListContactsExecutor) Execute(ctx context.Context, input application.ListContactsInput) ([]domain.Contact, error) {
	executor.called = true

	if executor.err != nil {
		return nil, executor.err
	}

	return executor.contacts, nil
}

// -----------------------------------------------------------------------------
// fakeAddDocumentExecutor
//
// Provides a fake add document workflow for API handler tests.
// -----------------------------------------------------------------------------
type fakeAddDocumentExecutor struct {
	document domain.Document
	err      error
	called   bool
}

// -----------------------------------------------------------------------------
// Execute
//
// Records add document workflow execution and returns the configured fake result.
// -----------------------------------------------------------------------------
func (executor *fakeAddDocumentExecutor) Execute(ctx context.Context, input application.AddDocumentInput) (domain.Document, error) {
	executor.called = true

	if executor.err != nil {
		return domain.Document{}, executor.err
	}

	return executor.document, nil
}

// -----------------------------------------------------------------------------
// fakeListDocumentsExecutor
//
// Provides a fake list documents workflow for API handler tests.
// -----------------------------------------------------------------------------
type fakeListDocumentsExecutor struct {
	documents []domain.Document
	err       error
	called    bool
}

// -----------------------------------------------------------------------------
// Execute
//
// Records list documents workflow execution and returns the configured fake result.
// -----------------------------------------------------------------------------
func (executor *fakeListDocumentsExecutor) Execute(ctx context.Context, input application.ListDocumentsInput) ([]domain.Document, error) {
	executor.called = true

	if executor.err != nil {
		return nil, executor.err
	}

	return executor.documents, nil
}

// -----------------------------------------------------------------------------
// TestHandleApplicationSearch
//
// Verifies that the search endpoint executes the search workflow.
// -----------------------------------------------------------------------------
func TestHandleApplicationSearch(t *testing.T) {
	application := newAPIHandlerTestApplication(t, "app-001", domain.StatusApplied)
	searchExecutor := &fakeSearchApplicationsExecutor{applications: []domain.Application{application}}
	handlers := NewWorkflowHandlers(WorkflowHandlerDependencies{SearchApplications: searchExecutor})

	request := httptest.NewRequest(http.MethodGet, "/applications/search?status=applied", nil)
	response := httptest.NewRecorder()

	handlers.HandleApplicationSearch(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if !searchExecutor.called {
		t.Fatal("expected search workflow to be called")
	}
}

// -----------------------------------------------------------------------------
// TestHandleApplicationWorkflowListsActivityEvents
//
// Verifies that the activity endpoint executes the activity timeline workflow.
// -----------------------------------------------------------------------------
func TestHandleApplicationWorkflowListsActivityEvents(t *testing.T) {
	event, err := domain.NewActivityEvent(
		"app-001",
		domain.ActivityStatusChanged,
		time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC),
		"Status changed from applied to interviewing.",
	)
	if err != nil {
		t.Fatalf("failed to create activity event: %v", err)
	}

	listExecutor := &fakeListActivityEventsExecutor{events: []domain.ActivityEvent{event}}
	handlers := NewWorkflowHandlers(WorkflowHandlerDependencies{ListActivityEvents: listExecutor})

	request := httptest.NewRequest(http.MethodGet, "/applications/app-001/activity", nil)
	response := httptest.NewRecorder()

	if !handlers.HandleApplicationWorkflow(response, request) {
		t.Fatal("expected workflow handler to route activity request")
	}

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if !listExecutor.called {
		t.Fatal("expected activity workflow to be called")
	}
}

// -----------------------------------------------------------------------------
// TestHandleApplicationWorkflowSchedulesReminder
//
// Verifies that POST /applications/{id}/reminders executes the schedule workflow.
// -----------------------------------------------------------------------------
func TestHandleApplicationWorkflowSchedulesReminder(t *testing.T) {
	reminder := newWorkflowHandlerTestReminder(t, "rem-001", false)
	scheduleExecutor := &fakeScheduleReminderExecutor{reminder: reminder}
	handlers := NewWorkflowHandlers(WorkflowHandlerDependencies{ScheduleReminder: scheduleExecutor})

	request := httptest.NewRequest(
		http.MethodPost,
		"/applications/app-001/reminders",
		strings.NewReader(`{"id":"rem-001","title":"Follow up","due_at":"2026-05-10T09:00:00Z"}`),
	)
	response := httptest.NewRecorder()

	if !handlers.HandleApplicationWorkflow(response, request) {
		t.Fatal("expected workflow handler to route reminder request")
	}

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
	}

	if !scheduleExecutor.called {
		t.Fatal("expected schedule reminder workflow to be called")
	}
}

// -----------------------------------------------------------------------------
// TestHandleApplicationWorkflowListsReminders
//
// Verifies that GET /applications/{id}/reminders executes the list reminders workflow.
// -----------------------------------------------------------------------------
func TestHandleApplicationWorkflowListsReminders(t *testing.T) {
	reminder := newWorkflowHandlerTestReminder(t, "rem-001", false)
	listExecutor := &fakeListRemindersExecutor{reminders: []domain.Reminder{reminder}}
	handlers := NewWorkflowHandlers(WorkflowHandlerDependencies{ListReminders: listExecutor})

	request := httptest.NewRequest(http.MethodGet, "/applications/app-001/reminders", nil)
	response := httptest.NewRecorder()

	if !handlers.HandleApplicationWorkflow(response, request) {
		t.Fatal("expected workflow handler to route reminder request")
	}

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if !listExecutor.called {
		t.Fatal("expected list reminders workflow to be called")
	}
}

// -----------------------------------------------------------------------------
// TestHandleReminderComplete
//
// Verifies that PATCH /reminders/{id}/complete executes the complete reminder workflow.
// -----------------------------------------------------------------------------
func TestHandleReminderComplete(t *testing.T) {
	reminder := newWorkflowHandlerTestReminder(t, "rem-001", true)
	completeExecutor := &fakeCompleteReminderExecutor{reminder: reminder}
	handlers := NewWorkflowHandlers(WorkflowHandlerDependencies{CompleteReminder: completeExecutor})

	request := httptest.NewRequest(http.MethodPatch, "/reminders/rem-001/complete", nil)
	response := httptest.NewRecorder()

	handlers.HandleReminderComplete(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if !completeExecutor.called {
		t.Fatal("expected complete reminder workflow to be called")
	}
}

// -----------------------------------------------------------------------------
// TestHandleApplicationWorkflowAddsAndListsContacts
//
// Verifies that contact collection routes execute contact workflows.
// -----------------------------------------------------------------------------
func TestHandleApplicationWorkflowAddsAndListsContacts(t *testing.T) {
	contact, err := domain.NewContact("contact-001", "app-001", "Sam Recruiter", "sam@example.com", "Recruiter")
	if err != nil {
		t.Fatalf("failed to create contact: %v", err)
	}

	addExecutor := &fakeAddContactExecutor{contact: contact}
	listExecutor := &fakeListContactsExecutor{contacts: []domain.Contact{contact}}
	handlers := NewWorkflowHandlers(WorkflowHandlerDependencies{
		AddContact:   addExecutor,
		ListContacts: listExecutor,
	})

	postRequest := httptest.NewRequest(
		http.MethodPost,
		"/applications/app-001/contacts",
		strings.NewReader(`{"id":"contact-001","name":"Sam Recruiter","email":"sam@example.com","role":"Recruiter"}`),
	)
	postResponse := httptest.NewRecorder()

	if !handlers.HandleApplicationWorkflow(postResponse, postRequest) {
		t.Fatal("expected workflow handler to route contact request")
	}

	if postResponse.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, postResponse.Code)
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/applications/app-001/contacts", nil)
	getResponse := httptest.NewRecorder()

	if !handlers.HandleApplicationWorkflow(getResponse, getRequest) {
		t.Fatal("expected workflow handler to route contact request")
	}

	if getResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getResponse.Code)
	}

	if !addExecutor.called || !listExecutor.called {
		t.Fatal("expected contact workflows to be called")
	}
}

// -----------------------------------------------------------------------------
// TestHandleApplicationWorkflowAddsAndListsDocuments
//
// Verifies that document collection routes execute document workflows.
// -----------------------------------------------------------------------------
func TestHandleApplicationWorkflowAddsAndListsDocuments(t *testing.T) {
	document, err := domain.NewDocument("doc-001", "app-001", "Backend Resume", "resume", "documents/backend-resume.pdf")
	if err != nil {
		t.Fatalf("failed to create document: %v", err)
	}

	addExecutor := &fakeAddDocumentExecutor{document: document}
	listExecutor := &fakeListDocumentsExecutor{documents: []domain.Document{document}}
	handlers := NewWorkflowHandlers(WorkflowHandlerDependencies{
		AddDocument:   addExecutor,
		ListDocuments: listExecutor,
	})

	postRequest := httptest.NewRequest(
		http.MethodPost,
		"/applications/app-001/documents",
		strings.NewReader(`{"id":"doc-001","name":"Backend Resume","kind":"resume","path":"documents/backend-resume.pdf"}`),
	)
	postResponse := httptest.NewRecorder()

	if !handlers.HandleApplicationWorkflow(postResponse, postRequest) {
		t.Fatal("expected workflow handler to route document request")
	}

	if postResponse.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, postResponse.Code)
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/applications/app-001/documents", nil)
	getResponse := httptest.NewRecorder()

	if !handlers.HandleApplicationWorkflow(getResponse, getRequest) {
		t.Fatal("expected workflow handler to route document request")
	}

	if getResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getResponse.Code)
	}

	if !addExecutor.called || !listExecutor.called {
		t.Fatal("expected document workflows to be called")
	}
}

// -----------------------------------------------------------------------------
// TestHandleApplicationWorkflowRejectsUnknownResource
//
// Verifies that unknown application subresources are not handled by workflow routes.
// -----------------------------------------------------------------------------
func TestHandleApplicationWorkflowRejectsUnknownResource(t *testing.T) {
	handlers := NewWorkflowHandlers(WorkflowHandlerDependencies{})

	request := httptest.NewRequest(http.MethodGet, "/applications/app-001/unknown", nil)
	response := httptest.NewRecorder()

	if handlers.HandleApplicationWorkflow(response, request) {
		t.Fatal("expected unknown workflow route not to be handled")
	}
}

// -----------------------------------------------------------------------------
// TestHandleReminderCompleteRejectsWorkflowError
//
// Verifies that complete reminder workflow errors return bad request.
// -----------------------------------------------------------------------------
func TestHandleReminderCompleteRejectsWorkflowError(t *testing.T) {
	completeExecutor := &fakeCompleteReminderExecutor{err: errors.New("complete failed")}
	handlers := NewWorkflowHandlers(WorkflowHandlerDependencies{CompleteReminder: completeExecutor})

	request := httptest.NewRequest(http.MethodPatch, "/reminders/rem-001/complete", nil)
	response := httptest.NewRecorder()

	handlers.HandleReminderComplete(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

// -----------------------------------------------------------------------------
// newWorkflowHandlerTestReminder
//
// Creates a valid reminder for expanded workflow handler tests.
// -----------------------------------------------------------------------------
func newWorkflowHandlerTestReminder(t *testing.T, id domain.ReminderID, completed bool) domain.Reminder {
	t.Helper()

	reminder, err := domain.NewReminder(
		id,
		"app-001",
		"Follow up",
		time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("failed to create reminder: %v", err)
	}

	reminder.Completed = completed

	return reminder
}
