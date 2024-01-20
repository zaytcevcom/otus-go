package schedulerstorage

import (
	"context"
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

func (s *Storage) GetNeedNotify(ctx context.Context, t time.Time) (events []storage.Event) {
	query := `
		SELECT
    		*
		FROM
			events
		WHERE
		    DATE(time_from - INTERVAL '1 SECOND' * notification_time) < $1
	`
	_ = s.db.SelectContext(ctx, &events, query, t)
	return events
}

func (s *Storage) DeleteOldest(ctx context.Context, t time.Time) (err error) {
	query := `DELETE FROM events WHERE time_from <= $1`
	_, err = s.db.ExecContext(ctx, query, t)
	return err
}
