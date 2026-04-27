package outbox

import (
	"database/sql"
	"time"
)

type Event struct {
	ID          string
	EventType   string
	Payload     []byte
	CreatedAt   time.Time
	Processed   bool
	Attempts    int
	LastError   sql.NullString
	NextRetryAt sql.NullTime
}

type TaskCreatedPayload struct {
	TaskID      string `json:"task_id"`
	Description string `json:"description"`
}

type TaskUpdatedPayload struct {
	TaskID      string `json:"task_id"`
	Description string `json:"description"`
}

type TaskDeletedPayload struct {
	TaskID string `json:"task_id"`
}
