-- +goose Up
CREATE TABLE gator.posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    description TEXT,
    published_at TIMESTAMP,
    feed_id UUID NOT NULL,
    CONSTRAINT fk_post_feed
    FOREIGN KEY (feed_id)
    REFERENCES gator.feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE gator.posts;

