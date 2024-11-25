-- name: DeleteUsers :exec
DELETE FROM gator.users;

-- name: GetUserByName :one
SELECT * FROM gator.users 
WHERE name = $1;

-- name: CreateUser :one
INSERT INTO gator.users (created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetUsers :many
SELECT * FROM gator.users;

-- name: GetUserName :one
SELECT name from gator.users
WHERE id = $1;