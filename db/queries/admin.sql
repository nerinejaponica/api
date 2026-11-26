-- name: GetAdminByEmail :one
SELECT * FROM admin_users WHERE email = sqlc.arg(email);

-- name: CreateAdmin :one
INSERT INTO admin_users (email, password_hash)
VALUES (sqlc.arg(email), sqlc.arg(password_hash))
RETURNING *;
