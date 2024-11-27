-- name: CreatePost :exec
INSERT INTO gator.posts (created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (NOW(),
        NOW(),
        $1,
        $2,
        $3,
        $4,
        $5
);

-- name: GetUserPosts :many
SELECT p.title, p.url, p.description, p.published_at FROM gator.posts p
JOIN gator.feeds f on p.feed_id = f.id
WHERE f.added_by = $1
ORDER BY p.published_at DESC
LIMIT $2;
