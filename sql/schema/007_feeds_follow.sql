-- +goose Up
CREATE TABLE gator.feed_follows (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT fk_users_follow
    FOREIGN KEY (user_id)
    REFERENCES gator.users(id) ON DELETE CASCADE,
    feed_id UUID NOT NULL,
    CONSTRAINT fk_feed_follow
    FOREIGN KEY (feed_id)
    REFERENCES gator.feeds(id) ON DELETE CASCADE,    
    UNIQUE(user_id, feed_id)

);

-- +goose Down
DROP TABLE gator.feed_follows;

