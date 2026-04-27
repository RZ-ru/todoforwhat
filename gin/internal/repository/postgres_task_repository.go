package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
	"todo/gin/internal/models"
	"todo/gin/internal/outbox"

	"github.com/google/uuid"
)

type PostgresTaskRepo struct {
	db *sql.DB
}

func NewPostgresTaskRepository(db *sql.DB) TaskRepository {
	return &PostgresTaskRepo{db: db}
}

func (r *PostgresTaskRepo) CreateTask(ctx context.Context, task models.Task) (models.Task, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Task{}, err
	}

	query := `
	INSERT INTO tasks (description, done, created_at)
	VALUES ($1,$2,$3)
	RETURNING id
	`

	err = tx.QueryRowContext(
		ctx,
		query,
		task.Description,
		task.Done,
		task.CreatedAt,
	).Scan(&task.Id)

	if err != nil {
		tx.Rollback()
		return models.Task{}, err
	}

	payload := outbox.TaskCreatedPayload{
		TaskID:      string(task.Id),
		Description: task.Description,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		tx.Rollback()
		return models.Task{}, err
	}

	event := outbox.Event{
		ID:        uuid.New().String(),
		EventType: "task_created",
		Payload:   data,
		CreatedAt: time.Now().UTC(),
		Processed: false,
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO outbox (id, event_type, payload, created_at, processed)
		 VALUES ($1,$2,$3,$4,$5)`,
		event.ID,
		event.EventType,
		event.Payload,
		event.CreatedAt,
		event.Processed,
	)

	if err != nil {
		tx.Rollback()
		return models.Task{}, err
	}

	err = tx.Commit()
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}
func (r *PostgresTaskRepo) GetTasks(ctx context.Context) ([]models.Task, error) {
	// слайс для хранения тасок
	tasks := []models.Task{}
	// параметр для формирования результатов
	query := `
	SELECT id, description, done, created_at
	FROM tasks
	`
	// ставим курсор в начало результирующей таблицы
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return []models.Task{}, err
	}
	// закрывает лишнее подключение
	defer rows.Close()
	// идем по циклу, пока есть строки
	for rows.Next() {
		var task models.Task
		// сканирование строки результата по переменным
		err = rows.Scan(&task.Id, &task.Description, &task.Done, &task.CreatedAt)
		if err != nil {
			return []models.Task{}, err
		}
		tasks = append(tasks, task)
	}
	err = rows.Err()
	if err != nil {
		return []models.Task{}, err
	}
	return tasks, nil
}
func (r *PostgresTaskRepo) GetTaskByID(ctx context.Context, id models.ID) (models.Task, error) {
	task := models.Task{}
	query := `SELECT id, description, done, created_at
	FROM tasks
	WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&task.Id, &task.Description, &task.Done, &task.CreatedAt)
	if err == sql.ErrNoRows {
		return task, err
	}
	if err != nil {
		return task, err
	}
	return task, nil
}
func (r *PostgresTaskRepo) UpdateTask(ctx context.Context, task models.Task) (models.Task, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Task{}, err
	}

	query := `
	UPDATE tasks 
	SET description = $1, done = $2 
	WHERE id = $3`
	result, err := tx.ExecContext(ctx, query, task.Description, task.Done, task.Id)
	if err != nil {
		tx.Rollback()
		return models.Task{}, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return models.Task{}, sql.ErrNoRows
	}
	if rowsAffected == 0 {
		tx.Rollback()
		return models.Task{}, sql.ErrNoRows
	}

	payload := outbox.TaskUpdatedPayload{
		TaskID:      string(task.Id),
		Description: task.Description,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		tx.Rollback()
		return models.Task{}, err
	}
	event := outbox.Event{
		ID:        uuid.New().String(),
		EventType: "task_updated",
		Payload:   data,
		CreatedAt: time.Now().UTC(),
		Processed: false,
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO outbox (id, event_type, payload, created_at, processed) 
		VALUES($1,$2,$3,$4,$5)`,
		event.ID,
		event.EventType,
		event.Payload,
		event.CreatedAt,
		event.Processed,
	)
	if err != nil {
		tx.Rollback()
		return models.Task{}, err
	}
	err = tx.Commit()
	if err != nil {
		//tx.Rollback()
		return models.Task{}, err
	}

	return task, nil
}
func (r *PostgresTaskRepo) DeleteTask(ctx context.Context, id models.ID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `DELETE FROM tasks WHERE id = $1`
	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if rowsAffected == 0 {
		tx.Rollback()
		return sql.ErrNoRows
	}

	payload := outbox.TaskDeletedPayload{
		TaskID: string(id),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		tx.Rollback()
		return err
	}

	event := outbox.Event{
		ID:        uuid.New().String(),
		EventType: "task_deleted",
		Payload:   data,
		CreatedAt: time.Now().UTC(),
		Processed: false,
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO outbox (id, event_type, payload, created_at, processed) 
		VALUES ($1, $2, $3, $4, $5)`,
		event.ID,
		event.EventType,
		event.Payload,
		event.CreatedAt,
		event.Processed,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresTaskRepo) GetUnprocessedEvents(ctx context.Context) ([]outbox.Event, error) {
	var events []outbox.Event

	query := `
	SELECT id, event_type, payload, created_at, processed, attempts, last_error, next_retry_at
	FROM outbox 
	WHERE processed = false
	AND attempts < 5
	AND (next_retry_at IS NULL OR next_retry_at <= NOW())`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var event outbox.Event

		err := rows.Scan(
			&event.ID,
			&event.EventType,
			&event.Payload,
			&event.CreatedAt,
			&event.Processed,
			&event.Attempts,
			&event.LastError,
			&event.NextRetryAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
func (r *PostgresTaskRepo) MarkEventProcessed(ctx context.Context, id string) error {
	query := `
	UPDATE outbox
	SET processed = true
	WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
func (r *PostgresTaskRepo) MarkEventFailed(ctx context.Context, id string, errorMsg string) error {
	query := `
	UPDATE outbox 
	SET
		attemps = attemps + 1,
		last_error = $2,
		next_retry_at = NOW() + INTERVAL '10 seconds'
	WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, errorMsg)
	return err
}
