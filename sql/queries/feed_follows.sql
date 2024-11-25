-- name: AddFeedFollow :one
WITH feedfollow AS (
INSERT INTO gator.feed_follows (created_at, updated_at, user_id, feed_id)
VALUES (
    NOW(),
    NOW(),
    $1,
    $2
)

RETURNING *
) 

SELECT feedfollow.*, u.name AS user_name, f.name as feed_name
FROM feedfollow
JOIN gator.users u on u.id = feedfollow.user_id
JOIN gator.feeds f on f.id = feedfollow.feed_id;

-- name: GetFeedsForUser :many
SELECT f.*, u.name AS user_name
FROM gator.feed_follows ff
JOIN gator.users u on ff.user_id = u.id
JOIN gator.feeds f on ff.feed_id = f.id
WHERE ff.user_id = $1;