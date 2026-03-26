-- name: GetNote :one
SELECT * FROM notes WHERE id = ? LIMIT 1;

-- name: GetNotesAbbr :many
SELECT id, title, createdAt FROM notes ORDER BY createdAt DESC;

-- name: DeleteNote :exec
DELETE FROM notes WHERE id = ?;

-- name: CreateNote :one
INSERT INTO notes (title, content) VALUES(?, ?) RETURNING *;
