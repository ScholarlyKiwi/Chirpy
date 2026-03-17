-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *;

-- name: GetRefreshTokenByUserID :many
SELECT *
  FROM refresh_tokens
 WHERE user_id = $1
 ORDER BY expires_at desc;

-- name: GetRefreshTokenByToken :one
SELECT *
  FROM refresh_tokens
 WHERE token = $1;

-- name: GetUserByRefreshToken :one
SELECT r.expires_at,
       r.token,
       r.revoked_at,
       u.*
  FROM refresh_tokens   r
  JOIN users            u ON r.user_id = u.id
 WHERE token = $1;

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
   SET revoked_at = NOW(),
       updated_at = NOW()
 WHERE token = $1
RETURNING *;