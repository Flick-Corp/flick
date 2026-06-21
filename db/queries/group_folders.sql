-- name: CreateGroupFolder :one
INSERT INTO group_folders (group_id, parent_id, name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetGroupFolderByID :one
SELECT id, group_id, parent_id, name, created_at FROM group_folders
WHERE id = $1;

-- name: ListGroupFoldersByParent :many
SELECT id, group_id, parent_id, name, created_at FROM group_folders
WHERE group_id = $1 AND parent_id IS NOT DISTINCT FROM $2
ORDER BY name;

-- name: DeleteGroupFolder :exec
DELETE FROM group_folders
WHERE id = $1;
