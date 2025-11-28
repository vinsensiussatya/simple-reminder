package usecase

import (
	"simple-reminder/internal/core"
)

type ReminderUsecase struct {
	repo core.ReminderRepo
}

func NewReminderUsecase(r core.ReminderRepo) *ReminderUsecase {
	return &ReminderUsecase{repo: r}
}

func (u *ReminderUsecase) Add(reminder *core.Reminder) error {
	return u.repo.Save(reminder)
}

func (u *ReminderUsecase) List() ([]*core.Reminder, error) {
	return u.repo.List()
}

func (u *ReminderUsecase) Delete(id string) error {
	return u.repo.Delete(id)
}
