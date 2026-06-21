-- name: CreateGroupUpload :one
INSERT INTO group_uploads (group_id, folder_id, code, uploader_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListGroupUploadsByFolder :many
SELECT gu.id, gu.code, gu.uploader_id, u.username AS uploader_username, gu.created_at
FROM group_uploads gu
JOIN users u ON u.id = gu.uploader_id
WHERE gu.group_id = $1 AND gu.folder_id IS NOT DISTINCT FROM $2
ORDER BY gu.created_at DESC;
