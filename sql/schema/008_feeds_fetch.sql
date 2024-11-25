-- +goose Up
ALTER TABLE gator.feeds
ADD COLUMN last_fetched_at TIMESTAMP;

-- +goose Down
ALTER TABLE gator.feeds
DROP COLUMN last_fetched_at;

