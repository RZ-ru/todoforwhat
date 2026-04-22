package models

import "time"

type ID string

type Task struct {
	Id          ID
	Description string
	Done        bool
	CreatedAt   time.Time
}
