-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
   id VARCHAR(255) NOT NULL PRIMARY KEY,
   title VARCHAR(255) NOT NULL,
   time_from TIMESTAMP NOT NULL,
   time_to TIMESTAMP NOT NULL,
   description TEXT,
   user_id VARCHAR(255) NOT NULL,
   notification_time INTERVAL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table events
-- +goose StatementEnd
