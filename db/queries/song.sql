-- name: ListSongs :many
-- Sementara semua status tampil; Fase 4 memisahkan publik (published) dan admin.
SELECT id, slug, title_ja, status, published_at
FROM songs
ORDER BY published_at DESC NULLS LAST, id DESC
LIMIT 50;

-- name: CreateSong :one
INSERT INTO songs (slug, title_ja, status, published_at)
VALUES (sqlc.arg(slug), sqlc.arg(title_ja), sqlc.arg(status), sqlc.narg(published_at))
RETURNING id, slug, title_ja, status, published_at;
