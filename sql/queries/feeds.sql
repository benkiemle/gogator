-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id, last_fetched_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    NULL
)
RETURNING *;

-- name: GetFeedsView :many
SELECT 
    f.id, 
    f.name, 
    f.url, 
    u.name AS user_name
FROM feeds f
INNER JOIN users u
    ON f.user_id = u.id;

-- name: GetFeedByUrl :one
SELECT
    id,
    name,
    created_at,
    updated_at,
    url,
    user_id
FROM feeds
WHERE url = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;

-- name: MarkFeedFetched :exec 
UPDATE feeds
SET updated_at = NOW(),
    last_fetched_at = NOW()
WHERE id = $1;
