-- name: GetUrl :one
SELECT * FROM urls
WHERE original_url = ? LIMIT 1;

-- name: ListUrls :many
SELECT * FROM urls
ORDER BY original_url ASC;

-- name: CreateUrl :one
INSERT INTO urls (
  original_url, shortened_url
) VALUES (
  ?, ?
)
RETURNING *;

-- name: UpdateUrl :one
UPDATE urls
  set shortened_url = ?
WHERE original_url = ?
RETURNING *;

-- name: DeleteUrl :exec
DELETE FROM urls
WHERE original_url = ?;
