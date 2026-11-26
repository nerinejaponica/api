-- name: ListArtists :many
SELECT * FROM artists ORDER BY name_romaji, name_ja;

-- name: GetArtistBySlug :one
SELECT * FROM artists WHERE slug = sqlc.arg(slug);

-- name: CreateArtist :one
INSERT INTO artists (slug, name_ja, name_romaji)
VALUES (sqlc.arg(slug), sqlc.arg(name_ja), sqlc.narg(name_romaji))
RETURNING *;

-- name: ListSources :many
SELECT * FROM sources ORDER BY kind, name_romaji, name_ja;

-- name: GetSourceBySlug :one
SELECT * FROM sources WHERE slug = sqlc.arg(slug);

-- name: CreateSource :one
INSERT INTO sources (slug, name_ja, name_romaji, kind)
VALUES (sqlc.arg(slug), sqlc.arg(name_ja), sqlc.narg(name_romaji), sqlc.arg(kind))
RETURNING *;

-- name: ListTags :many
SELECT * FROM tags ORDER BY name;

-- name: GetTagBySlug :one
SELECT * FROM tags WHERE slug = sqlc.arg(slug);

-- name: CreateTag :one
INSERT INTO tags (slug, name)
VALUES (sqlc.arg(slug), sqlc.arg(name))
RETURNING *;
