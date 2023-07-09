CREATE TABLE IF NOT EXISTS composers (
   id BIGSERIAL PRIMARY KEY,
   first_name TEXT NOT NULL,
   last_name TEXT NOT NULL,
   image_link TEXT,
   description TEXT,
   version BIGINT NOT NULL DEFAULT 1
);

-- name: SelectComposer :one
SELECT * FROM composers
WHERE id = $1 LIMIT 1;

-- name: SelectComposers :many
SELECT * FROM composers
ORDER BY composers.id ASC
OFFSET $1 LIMIT $2;

-- name: InsertComposer :one
INSERT INTO composers (
   first_name, last_name, image_link, description
) VALUES (
   $1, $2, $3, $4
)
RETURNING *;

-- name: DeleteComposer :exec
DELETE FROM composers
WHERE id = $1;

-- name: UpdateComposer :one
UPDATE composers
SET first_name = $1, last_name = $2, image_link = $3, description = $4, version = version + 1
WHERE id = $5
RETURNING *;
