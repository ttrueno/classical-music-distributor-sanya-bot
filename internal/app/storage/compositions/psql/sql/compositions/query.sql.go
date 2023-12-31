// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: query.sql

package compositions

import (
	"context"
)

const deleteComposition = `-- name: DeleteComposition :exec
DELETE FROM compositions
WHERE id = $1
`

func (q *Queries) DeleteComposition(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteComposition, id)
	return err
}

const insertComposition = `-- name: InsertComposition :one
INSERT INTO compositions (
   composer_id,
   name
) VALUES (
   $1, $2
)
RETURNING id, composer_id, name, version
`

type InsertCompositionParams struct {
	ComposerID int64
	Name       string
}

func (q *Queries) InsertComposition(ctx context.Context, arg InsertCompositionParams) (Composition, error) {
	row := q.db.QueryRow(ctx, insertComposition, arg.ComposerID, arg.Name)
	var i Composition
	err := row.Scan(
		&i.ID,
		&i.ComposerID,
		&i.Name,
		&i.Version,
	)
	return i, err
}

const selectAllCompositionsByComposerID = `-- name: SelectAllCompositionsByComposerID :many
SELECT id, composer_id, name, version FROM compositions
WHERE composer_id = $1
OFFSET $2 LIMIT $3
`

type SelectAllCompositionsByComposerIDParams struct {
	ComposerID int64
	Offset     int32
	Limit      int32
}

func (q *Queries) SelectAllCompositionsByComposerID(ctx context.Context, arg SelectAllCompositionsByComposerIDParams) ([]Composition, error) {
	rows, err := q.db.Query(ctx, selectAllCompositionsByComposerID, arg.ComposerID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Composition
	for rows.Next() {
		var i Composition
		if err := rows.Scan(
			&i.ID,
			&i.ComposerID,
			&i.Name,
			&i.Version,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectComposition = `-- name: SelectComposition :one
SELECT id, composer_id, name, version FROM compositions
WHERE id = $1 LIMIT 1
`

func (q *Queries) SelectComposition(ctx context.Context, id int64) (Composition, error) {
	row := q.db.QueryRow(ctx, selectComposition, id)
	var i Composition
	err := row.Scan(
		&i.ID,
		&i.ComposerID,
		&i.Name,
		&i.Version,
	)
	return i, err
}

const updateComposition = `-- name: UpdateComposition :one
UPDATE compositions
SET composer_id = $1, name = $2, version = version + 1
WHERE id = $3
RETURNING id, composer_id, name, version
`

type UpdateCompositionParams struct {
	ComposerID int64
	Name       string
	ID         int64
}

func (q *Queries) UpdateComposition(ctx context.Context, arg UpdateCompositionParams) (Composition, error) {
	row := q.db.QueryRow(ctx, updateComposition, arg.ComposerID, arg.Name, arg.ID)
	var i Composition
	err := row.Scan(
		&i.ID,
		&i.ComposerID,
		&i.Name,
		&i.Version,
	)
	return i, err
}
