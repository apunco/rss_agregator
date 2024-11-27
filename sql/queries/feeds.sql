-- name: AddFeed :one
INSERT INTO gator.feeds (created_at, updated_at, name, url, added_by)
VALUES (NOW(),
        NOW(),
        $1,
        $2,
        $3
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM gator.feeds;

-- name: GetFeedByUrl :one
SELECT * FROM gator.feeds
WHERE url = $1;

-- name: MarkedFeedFetched :exec
UPDATE gator.feeds
SET updated_at = NOW(),
last_fetched_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * from gator.feeds
ORDER BY last_fetched_at asc
NULLS FIRST;