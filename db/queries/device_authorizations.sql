-- name: CreateDeviceAuthorization :one
INSERT INTO device_authorizations (device_code, user_code, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetDeviceAuthorizationByDeviceCode :one
SELECT * FROM device_authorizations
WHERE device_code = $1 AND expires_at > now();

-- name: GetDeviceAuthorizationByUserCode :one
SELECT * FROM device_authorizations
WHERE user_code = $1 AND expires_at > now();

-- name: ApproveDeviceAuthorization :one
UPDATE device_authorizations
SET status = 'approved', user_id = $2, session_token = $3
WHERE user_code = $1 AND status = 'pending'
RETURNING *;

-- name: DenyDeviceAuthorization :one
UPDATE device_authorizations
SET status = 'denied'
WHERE user_code = $1 AND status = 'pending'
RETURNING *;

-- name: DeleteExpiredDeviceAuthorizations :exec
DELETE FROM device_authorizations
WHERE expires_at <= now();
