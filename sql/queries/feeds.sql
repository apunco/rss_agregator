-- name: AddFeed :one
INSERT INTO gator.feeds (created_at, updated_at, name, url, added_by)
VALUES (NOW(),
        NOW(),
        $1,
        $2,
        $3
)
RETURNING *;