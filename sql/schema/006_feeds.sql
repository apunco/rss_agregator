-- +goose Up
CREATE TABLE gator.feeds (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT UNIQUE NOT NULL,
    url TEXT UNIQUE NOT NULL,
    added_by UUID NOT NULL,
    CONSTRAINT fk_users
    FOREIGN KEY (added_by)
    REFERENCES gator.users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE gator.feeds;

