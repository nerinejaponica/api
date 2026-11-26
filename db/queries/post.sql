-- name: ListPosts :many
SELECT id, slug, title, excerpt, cover_url, status, published_at
FROM posts
ORDER BY published_at DESC NULLS LAST, id DESC
LIMIT 50;

-- name: GetPostBySlug :one
SELECT * FROM posts WHERE slug = sqlc.arg(slug);

-- name: CreatePost :one
INSERT INTO posts (slug, title, excerpt, cover_url, body, status, published_at)
VALUES (
    sqlc.arg(slug), sqlc.arg(title), sqlc.narg(excerpt), sqlc.narg(cover_url),
    sqlc.arg(body), sqlc.arg(status), sqlc.narg(published_at)
)
RETURNING *;

-- name: DeletePost :execrows
DELETE FROM posts WHERE id = sqlc.arg(id);
