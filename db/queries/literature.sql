-- name: ListLiterature :many
SELECT id, slug, title_ja, title_id, original_author, year, cover_url, status, published_at
FROM literature
ORDER BY published_at DESC NULLS LAST, id DESC
LIMIT 50;

-- name: GetLiteratureBySlug :one
SELECT * FROM literature WHERE slug = sqlc.arg(slug);

-- name: CreateLiterature :one
INSERT INTO literature (slug, title_ja, title_id, original_author, year, aozora_url, cover_url, body, footnotes, status, published_at)
VALUES (
    sqlc.arg(slug), sqlc.arg(title_ja), sqlc.narg(title_id), sqlc.narg(original_author), sqlc.narg(year),
    sqlc.narg(aozora_url), sqlc.narg(cover_url), sqlc.arg(body), sqlc.arg(footnotes), sqlc.arg(status), sqlc.narg(published_at)
)
RETURNING *;

-- name: DeleteLiterature :execrows
DELETE FROM literature WHERE id = sqlc.arg(id);
