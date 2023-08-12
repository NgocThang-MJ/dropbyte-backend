-- name: CreateFile :one
INSERT INTO files (
  file_id,
  bucket_id,
  owner,
  name,
  size,
  file_type
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetFile :one
SELECT * FROM files
WHERE id = $1 LIMIT 1;

-- name: ListFiles :many
SELECT * FROM files
WHERE owner = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: UpdateFile :one
UPDATE files
  set name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteFile :exec
DELETE FROM files
WHERE file_id = $1;
