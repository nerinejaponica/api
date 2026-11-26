-- name: AddSongArtist :exec
INSERT INTO song_artists (song_id, artist_id)
VALUES (sqlc.arg(song_id), sqlc.arg(artist_id))
ON CONFLICT DO NOTHING;

-- name: ListSongArtists :many
SELECT a.*
FROM artists a
JOIN song_artists sa ON sa.artist_id = a.id
WHERE sa.song_id = sqlc.arg(song_id)
ORDER BY a.name_ja;

-- name: AddSongSource :exec
INSERT INTO song_sources (song_id, source_id)
VALUES (sqlc.arg(song_id), sqlc.arg(source_id))
ON CONFLICT DO NOTHING;

-- name: ListSongSources :many
SELECT s.*
FROM sources s
JOIN song_sources ss ON ss.source_id = s.id
WHERE ss.song_id = sqlc.arg(song_id)
ORDER BY s.name_ja;

-- name: AddSongTag :exec
INSERT INTO content_tags (tag_id, song_id)
VALUES (sqlc.arg(tag_id), sqlc.arg(song_id))
ON CONFLICT DO NOTHING;

-- name: AddLiteratureTag :exec
INSERT INTO content_tags (tag_id, literature_id)
VALUES (sqlc.arg(tag_id), sqlc.arg(literature_id))
ON CONFLICT DO NOTHING;

-- name: ListSongTags :many
SELECT t.*
FROM tags t
JOIN content_tags ct ON ct.tag_id = t.id
WHERE ct.song_id = sqlc.arg(song_id)
ORDER BY t.name;

-- name: ListLiteratureTags :many
SELECT t.*
FROM tags t
JOIN content_tags ct ON ct.tag_id = t.id
WHERE ct.literature_id = sqlc.arg(literature_id)
ORDER BY t.name;

-- name: AddPostTag :exec
INSERT INTO content_tags (tag_id, post_id)
VALUES (sqlc.arg(tag_id), sqlc.arg(post_id))
ON CONFLICT DO NOTHING;

-- name: ListPostTags :many
SELECT t.*
FROM tags t
JOIN content_tags ct ON ct.tag_id = t.id
WHERE ct.post_id = sqlc.arg(post_id)
ORDER BY t.name;
