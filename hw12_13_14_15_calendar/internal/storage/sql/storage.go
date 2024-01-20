package sqlstorage

import (
	"context"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/internal/storage"
	"time"
)

type Storage struct {
	dsn string
	db  *sqlx.DB
}

func New(dsn string) *Storage {
	return &Storage{
		dsn: dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) (err error) {
	s.db, err = sqlx.Open("pgx", s.dsn)
	if err != nil {
		return err
	}

	return s.db.PingContext(ctx)
}

func (s *Storage) Close(_ context.Context) error {
	return s.db.Close()
}

func (s *Storage) GetEventsByDay(ctx context.Context, t time.Time) (events []storage.Event) {
	query := `
		SELECT
    		*
		FROM
			events
		WHERE
		    DATE(time_from) = DATE($1)`
	err := s.db.SelectContext(ctx, &events, query, t)
	if err != nil {
		return nil
	}
	return events
}

func (s *Storage) GetEventsByWeek(ctx context.Context, t time.Time) (events []storage.Event) {
	query := `
		SELECT
		    *
		FROM
		    events
		WHERE
		    DATE_PART('week', time_from) = DATE_PART('week', $1)
	`
	err := s.db.SelectContext(ctx, &events, query, t)

	if err != nil {
		return nil
	}
	return events
}

func (s *Storage) GetEventsByMonth(ctx context.Context, t time.Time) (events []storage.Event) {
	query := `
		SELECT
			*
		FROM
		    events
		WHERE
		    DATE_PART('month', time_from) = DATE_PART('month', $1)
	`
	err := s.db.SelectContext(ctx, &events, query, t)
	if err != nil {
		return nil
	}
	return events
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (err error) {
	query := `
		INSERT INTO events (
			id, title, time_from, time_to, description, user_id, notification_time
		) VALUES (
		  :id, :title, :time_from, :time_to, :description, :user_id, :notification_time
	  	)
	`
	event.ID = uuid.NewString()
	_, err = s.db.NamedExecContext(ctx, query, event)
	return err
}

func (s *Storage) UpdateEvent(ctx context.Context, id string, event storage.Event) (err error) {
	query := `
		UPDATE
    		events
		SET
		    title=:title,
		    time_from=:time_from,
		    time_to=:time_to,
		    description=:description,
		    user_id=:user_id,
		    notification_time=:notification_time
		WHERE
		    id=:id
	`
	event.ID = id
	_, err = s.db.NamedExecContext(ctx, query, event)
	return err
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) (err error) {
	query := `DELETE FROM events WHERE id=$1`
	_, err = s.db.ExecContext(ctx, query, id)
	return err
}
