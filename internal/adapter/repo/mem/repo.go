package mem

import (
	"errors"
	"simple-reminder/internal/core"
)

type ReminderRepo struct {
	store map[string]*core.Reminder
}

func NewReminderRepo() *ReminderRepo {
	return &ReminderRepo{store: make(map[string]*core.Reminder)}
}

func (r *ReminderRepo) Save(rem *core.Reminder) error {
	r.store[rem.ID] = rem
	return nil
}

func (r *ReminderRepo) List() ([]*core.Reminder, error) {
	reminders := make([]*core.Reminder, 0, len(r.store))
	for _, rem := range r.store {
		reminders = append(reminders, rem)
	}
	return reminders, nil
}

func (r *ReminderRepo) Delete(id string) error {
	if _, ok := r.store[id]; !ok {
		return errors.New("not found")
	}
	delete(r.store, id)
	return nil
}
