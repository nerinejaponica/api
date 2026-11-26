-- name: ListSiteSettings :many
SELECT * FROM site_settings ORDER BY key;

-- name: UpsertSiteSetting :one
INSERT INTO site_settings (key, value)
VALUES (sqlc.arg(key), sqlc.arg(value))
ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = now()
RETURNING *;
