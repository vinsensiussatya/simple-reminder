package postgres

import (
	"context"
	"errors"
	"os"
	"simple-reminder/internal/core"

	"github.com/jackc/pgx/v5"
)

type ReminderRepo struct {
	conn *pgx.Conn
}

func NewReminderRepo(ctx context.Context) (*ReminderRepo, error) {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		return nil, errors.New("DB_URL env var not set")
	}
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &ReminderRepo{conn: conn}, nil
}

func (r *ReminderRepo) Ping(ctx context.Context) error {
	return r.conn.Ping(ctx)
}

func (r *ReminderRepo) Save(rem *core.Reminder) error {
	_, err := r.conn.Exec(context.Background(),
		"INSERT INTO reminders (id, message, remind_at) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET message=$2, remind_at=$3",
		rem.ID, rem.Message, rem.RemindAt)
	return err
}

func (r *ReminderRepo) List() ([]*core.Reminder, error) {
	rows, err := r.conn.Query(context.Background(), "SELECT id, message, remind_at FROM reminders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reminders := []*core.Reminder{}
	for rows.Next() {
		var rem core.Reminder
		if err := rows.Scan(&rem.ID, &rem.Message, &rem.RemindAt); err != nil {
			return nil, err
		}
		reminders = append(reminders, &rem)
	}
	return reminders, nil
}

func (r *ReminderRepo) Delete(id string) error {
	ct, err := r.conn.Exec(context.Background(), "DELETE FROM reminders WHERE id = $1", id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}
