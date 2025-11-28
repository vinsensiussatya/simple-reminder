package core

import "time"

type Reminder struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	RemindAt  time.Time `json:"remind_at"`
}

type ReminderRepo interface {
	Save(r *Reminder) error
	List() ([]*Reminder, error)
	Delete(id string) error
}
