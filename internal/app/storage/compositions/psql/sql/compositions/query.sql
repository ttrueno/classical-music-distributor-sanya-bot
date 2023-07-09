CREATE TABLE IF NOT EXISTS compositions (
   id BIGSERIAL PRIMARY KEY,
   composer_id BIGINT NOT NULL,
   name TEXT NOT NULL,
   version BIGINT NOT NULL DEFAULT 1,
   CONSTRAINT fk_composer
      FOREIGN KEY(composer_id)
         REFERENCES composers(id)
);

-- name: SelectComposition :one
SELECT * FROM compositions
WHERE id = $1 LIMIT 1;

-- name: SelectAllCompositionsByComposerID :many
SELECT * FROM compositions
WHERE composer_id = $1
OFFSET $2 LIMIT $3;

-- name: InsertComposition :one
INSERT INTO compositions (
   composer_id,
   name
) VALUES (
   $1, $2
)
RETURNING *;

-- name: DeleteComposition :exec
DELETE FROM compositions
WHERE id = $1;

-- name: UpdateComposition :one
UPDATE compositions
SET composer_id = $1, name = $2, version = version + 1
WHERE id = $3
RETURNING *;
