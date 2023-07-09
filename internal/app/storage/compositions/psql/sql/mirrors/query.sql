CREATE TABLE IF NOT EXISTS composition_mirrors (
   id BIGSERIAL PRIMARY KEY,
   composition_id BIGINT NOT NULL,
   link TEXT NOT NULL,
   version BIGINT NOT NULL DEFAULT 1,
   CONSTRAINT fk_composition
      FOREIGN KEY(composition_id) 
         REFERENCES compositions(id)
);

-- name: GetCompositionMirror :one
SELECT * FROM composition_mirrors
WHERE id = $1 LIMIT 1;

-- name: GetCompositionsMirrorsByCompositionID :many
SELECT * FROM composition_mirrors
WHERE composition_id = $1;

-- name: InsertCompositionMirror :one
INSERT INTO composition_mirrors (
   composition_id, link
) VALUES (
   $1, $2
)
RETURNING *;

-- name: DeleteCompositionMirror :exec
DELETE FROM composition_mirrors
WHERE id = $1;

-- name: UpdateCompositionMirror :one
UPDATE composition_mirrors
SET composition_id = $1, link = $2, version = version + 1
WHERE id = $3
RETURNING *;
