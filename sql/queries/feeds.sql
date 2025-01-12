-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds
ORDER BY created_at DESC;

-- name: GetFeedURL :one 
SELECT * FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $2, updated_at = $3
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
WHERE user_id = $1
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;

-- name: UpdateFeedUser :exec
DELETE FROM feeds
WHERE user_id = $1 AND url = $2;
