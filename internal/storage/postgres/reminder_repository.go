package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// SaveReminder
//
// Inserts or updates a reminder for an application.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) SaveReminder(ctx context.Context, reminder domain.Reminder) error {
	if repository.db == nil {
		return fmt.Errorf("database connection is required")
	}

	if err := reminder.Validate(); err != nil {
		return err
	}

	_, err := repository.db.ExecContext(
		ctx,
		`
        INSERT INTO reminders (
            id,
            application_id,
            title,
            due_at,
            completed
        )
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (id)
        DO UPDATE SET
            application_id = EXCLUDED.application_id,
            title = EXCLUDED.title,
            due_at = EXCLUDED.due_at,
            completed = EXCLUDED.completed,
            updated_at = NOW()
        `,
		reminder.ID,
		reminder.ApplicationID,
		reminder.Title,
		reminder.DueAt,
		reminder.Completed,
	)

	return err
}

// -----------------------------------------------------------------------------
// FindReminderByID
//
// Retrieves one reminder by its stable domain identity.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) FindReminderByID(ctx context.Context, id domain.ReminderID) (domain.Reminder, error) {
	if repository.db == nil {
		return domain.Reminder{}, fmt.Errorf("database connection is required")
	}

	if err := id.Validate(); err != nil {
		return domain.Reminder{}, err
	}

	row := repository.db.QueryRowContext(
		ctx,
		`
        SELECT
            id,
            application_id,
            title,
            due_at,
            completed
        FROM reminders
        WHERE id = $1
        `,
		id,
	)

	reminder, err := scanReminder(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Reminder{}, fmt.Errorf("reminder not found: %s", id)
		}

		return domain.Reminder{}, err
	}

	return reminder, nil
}

// -----------------------------------------------------------------------------
// UpdateReminder
//
// Updates an existing reminder.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) UpdateReminder(ctx context.Context, reminder domain.Reminder) error {
	if repository.db == nil {
		return fmt.Errorf("database connection is required")
	}

	if err := reminder.Validate(); err != nil {
		return err
	}

	result, err := repository.db.ExecContext(
		ctx,
		`
        UPDATE reminders
        SET
            title = $2,
            due_at = $3,
            completed = $4,
            updated_at = NOW()
        WHERE id = $1
        `,
		reminder.ID,
		reminder.Title,
		reminder.DueAt,
		reminder.Completed,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("reminder not found: %s", reminder.ID)
	}

	return nil
}

// -----------------------------------------------------------------------------
// RemoveReminder
//
// Removes one reminder.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) RemoveReminder(ctx context.Context, id domain.ReminderID) error {
	if repository.db == nil {
		return fmt.Errorf("database connection is required")
	}

	if err := id.Validate(); err != nil {
		return err
	}

	result, err := repository.db.ExecContext(
		ctx,
		`
        DELETE FROM reminders
        WHERE id = $1
        `,
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("reminder not found: %s", id)
	}

	return nil
}

// -----------------------------------------------------------------------------
// ListRemindersForApplication
//
// Retrieves reminders for one application in database priority order.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) ListRemindersForApplication(ctx context.Context, applicationID domain.ApplicationID) ([]domain.Reminder, error) {
	if repository.db == nil {
		return nil, fmt.Errorf("database connection is required")
	}

	if err := applicationID.Validate(); err != nil {
		return nil, err
	}

	rows, err := repository.db.QueryContext(
		ctx,
		`
        SELECT
            id,
            application_id,
            title,
            due_at,
            completed
        FROM reminders
        WHERE application_id = $1
        ORDER BY completed ASC, due_at ASC, title ASC
        `,
		applicationID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applicationReminders := []domain.Reminder{}

	for rows.Next() {
		reminder, err := scanReminder(rows)
		if err != nil {
			return nil, err
		}

		applicationReminders = append(applicationReminders, reminder)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return applicationReminders, nil
}

// -----------------------------------------------------------------------------
// reminderScanner
//
// Defines the scan behavior shared by SQL reminder rows and row results.
// -----------------------------------------------------------------------------
type reminderScanner interface {
	Scan(dest ...any) error
}

// -----------------------------------------------------------------------------
// scanReminder
//
// Converts a database row into a domain reminder.
// -----------------------------------------------------------------------------
func scanReminder(scanner reminderScanner) (domain.Reminder, error) {
	var reminder domain.Reminder
	var completed sql.NullBool

	err := scanner.Scan(
		&reminder.ID,
		&reminder.ApplicationID,
		&reminder.Title,
		&reminder.DueAt,
		&completed,
	)
	if err != nil {
		return domain.Reminder{}, err
	}

	reminder.Completed = completed.Valid && completed.Bool

	if err := reminder.Validate(); err != nil {
		return domain.Reminder{}, err
	}

	return reminder, nil
}
